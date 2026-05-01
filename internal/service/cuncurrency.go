package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"goproject/internal/loader"
	"goproject/internal/validation"
	"strconv"
	"strings"
	"sync"

)

type Service struct {
	mu sync.RWMutex
	validUsers   []User
	invalidUsers []User
}

func NewService() *Service {
	return &Service{validUsers: []User{}, invalidUsers: []User{}}
}

type SetUsersStrategy interface {
	Parse(record []string, row int) (User, error)
}

func DetectStrategy(rec []string) SetUsersStrategy {

	if err := isCSVHeader(rec); err != nil {
		return JSONStrategy{}
	}

	return CSVStrategy{}
}


type job struct{
	Seq    int 
	Record []string	
}

type recordResult struct{
	Seq  		   int
	User           User
	ParseErr       error
	ValidationErrs []*validation.ValidationError
	IsHeader bool
}

func (s *Service) ProcessStream(records <-chan loader.Record, workers int)error{
	jobs := make(chan job)
	results := make(chan recordResult)
	producerErr := make(chan error, 1)

	
	first,_ := <- records
	strategy := DetectStrategy(first.Fields)

	//producer
	go func (){
		defer close(jobs)
		defer close(producerErr)
		if errs := isCSVHeader(first.Fields); errs == nil{

		}else if lookslikeCSVHeader(first.Fields){
			for range records {
			}
			producerErr <- fmt.Errorf("invalid csv header: got %v, want %v", first.Fields, []string{"id", "first_name", "last_name", "email"})
			return

		}

		jobs <- job{Seq: 0, Record: first.Fields}
		
		seq := 1
		for rec := range records{
			jobs <- job{Seq: seq, Record: rec.Fields}
			seq++
		}
		
		}()
		
	var wg sync.WaitGroup

	//workers
	for i:= 0 ;i < workers ; i++{
		wg.Add(1)
		
		go func (strategy SetUsersStrategy) {
			defer wg.Done()
			for job := range jobs{
				if _, ok := strategy.(CSVStrategy); ok && job.Seq ==0 {
					results <- recordResult{Seq: job.Seq, IsHeader: true}
					continue
				}

				u, err := strategy.Parse(job.Record, job.Seq)
			
				if err != nil {
					results <- recordResult{Seq: job.Seq, ParseErr: err}
					continue
				}

				errs := ValidateUser(u)
				results <- recordResult{Seq: job.Seq, User: u, ValidationErrs: errs}

			}			
		}(strategy)

	}
	//closer
	go func ()  {
		wg.Wait()
		close(results)
	}()

	//collector
	pending := make(map[int]recordResult)
	next := 0
	valids := make([]User, 0)
	invalids := make([]User, 0)

	for res := range results{
		pending[res.Seq] = res

		for {
			r, ok := pending[next]
			if !ok{
				break
			}
			if r.IsHeader{
				next++
				continue
			}

			delete(pending, next)
			if len(r.ValidationErrs) > 0{
				if hasEmailError(r.ValidationErrs) {
				r.User.Email = "invalid-email"
			}
				invalids = append(invalids, r.User)

			} else {
				valids = append(valids, r.User)
			}
			next++
		}
	}

	if err := <-producerErr; err != nil {
		return err
	}

	s.mu.Lock()
	s.invalidUsers = invalids
	s.validUsers = valids
	s.mu.Unlock()

	return nil
}

func (s *Service)AddUser(user User)[]*validation.ValidationError{
	s.mu.Lock()
	defer s.mu.Unlock()
	
	errs := ValidateUser(user)

	if s.isDuplicateEmail(user.Email){
		errs = append(errs, &validation.ValidationError{Field: "email",Message: "email already exists", Code: 1004 })
	}
	
	if errs != nil{
		return errs
	}

	s.validUsers = append(s.validUsers, user)
	return nil
}

func (s *Service) ValidCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.validUsers)
}

func (s *Service) InvalidCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.invalidUsers)
}

func (s *Service) GetValidUsers(page int, limit int, email string) ([]User, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]User, 0)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}


	if email != "" {
		for _, user := range s.validUsers {

			if user.Email == email {
				filtered = append(filtered, user)
			}
		}
	} else {
		filtered = s.validUsers
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []User{}, total
	}

	if end > total {
		end = total
	}

	users := filtered[start:end]

	return users, total
}

func (s *Service) GetInvalidUsers() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, len(s.invalidUsers))
	copy(out, s.invalidUsers)
	return out
}

type CSVStrategy struct{}

type JSONStrategy struct{}

func (CSVStrategy) Parse(record []string, row int) (User, error) {
	if len(record) != 4 {
		return User{}, fmt.Errorf("invalid csv row %d: expected 4 columns, got %d", row+2, len(record))
	}

	idStr := strings.TrimSpace(record[0])
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return User{}, fmt.Errorf("invalid id at csv row %d: %w", row+2, err)
	}

	return User{
		ID:        id,
		FirstName: strings.TrimSpace(record[1]),
		LastName:  strings.TrimSpace(record[2]),
		Email:     strings.TrimSpace(record[3]),
	}, nil
}

func (JSONStrategy) Parse(record []string, row int) (User, error) {

	if len(record) != 1 {
		return User{}, fmt.Errorf("json: row %d must have exactly 1 column (raw object), got %d", row+1, len(record))
	}

	raw := strings.TrimSpace(record[0])
	if raw == "" {
		return User{}, fmt.Errorf("json: empty object at row %d", row+1)
	}

	var obj map[string]any

	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return User{}, fmt.Errorf("invalid json at row %d: %w", row+1, err)
	}

	if _, ok := obj["id"]; !ok {
		return User{}, fmt.Errorf("json row %d: missing key 'id'", row+1)
	}
	if _, ok := obj["first_name"]; !ok {
		return User{}, fmt.Errorf("json row %d: missing key 'first_name'", row+1)
	}
	if _, ok := obj["last_name"]; !ok {
		return User{}, fmt.Errorf("json row %d: missing key 'last_name'", row+1)
	}
	if _, ok := obj["email"]; !ok {
		return User{}, fmt.Errorf("json row %d: missing key 'email'", row+1)
	}

	id, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(obj["id"])))
	if err != nil {
		return User{}, fmt.Errorf("json row %d: invalid id: %w", row+1, err)
	}

	first := fmt.Sprint(obj["first_name"])
	last := fmt.Sprint(obj["last_name"])
	email := fmt.Sprint(obj["email"])

	return User{
		ID:        id,
		FirstName: strings.TrimSpace(first),
		LastName:  strings.TrimSpace(last),
		Email:     strings.TrimSpace(email),
	}, nil
}

func isCSVHeader(rec []string) error {
	
	if len(rec) != 4 {
		return errors.New("invalid header: expected 4 columns")
	}

	expected := [4]string{"id", "first_name", "last_name", "email"}
	for i := range expected {
		if rec[i] != expected[i] {
			return errors.New("invalid header column name")
		}
	}
	return nil
}


func lookslikeCSVHeader(rec []string)bool{
	if len(rec) != 4 {
		return false
	}

	for _, v := range rec{
		v = strings.TrimSpace(strings.ToLower(v))
		if v == "id" || v == "first_name" || v == "last_name" || v == "email" {
			return true
		}
	}
	return false
}

func ValidateUser(u User) (errs []*validation.ValidationError) {

	if err := validation.ValidID(u.ID); err != nil {
		errs = append(errs, err.(*validation.ValidationError))
	}
	if err := validation.ValidFirstname(u.FirstName); err != nil {
		errs = append(errs, err.(*validation.ValidationError))
	}
	if err := validation.ValidEmail(u.Email); err != nil {
		errs = append(errs, err.(*validation.ValidationError))
	}

	return errs
}

func hasEmailError(errs []*validation.ValidationError) bool {
	for _, e := range errs {
		if e != nil && e.Code == 1002 {
			return true
		}
	}
	return false
}

func (s *Service)isDuplicateEmail(email string)bool{
	for _, user := range s.validUsers{
		if user.Email == email{
			return true
		}

	}
	return false
}
