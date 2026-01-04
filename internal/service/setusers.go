package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"goproject/internal/validation"
	"strconv"
	"strings"
)

type Service struct {
	validUsers   []User
	invalidUsers []User
}

func NewService() *Service {
	return &Service{validUsers: []User{}, invalidUsers: []User{}}
}

type SetUsersStrategy interface {
	Parse(record []string, row int) (User, error)
}

func DetectStrategy(records [][]string) SetUsersStrategy {

	if err := isCSVHeader(records); err != nil {
		return JSONStrategy{}
	}

	return CSVStrategy{}

}

func (s *Service) SetUsers(records [][]string, workers int) error {

	if len(records) == 0 {
		s.validUsers = []User{}
		s.invalidUsers = []User{}
		return nil
	}

	strategy := DetectStrategy(records)

	data := records
	if err := isCSVHeader(records); err == nil {
		data = records[1:]
	}

	results, err := RunPool(data, workers, strategy)
	if err != nil {
		return err
	}

	var validUsers []User
	var invalidUsers []User

	for _, result := range results {
		if result.ParseErr != nil {
			return fmt.Errorf("error parsing record %d: %w", result.Index, result.ParseErr)
		} else if len(result.ValidationErrs) > 0 {
			if hasEmailError(result.ValidationErrs) {
				result.User.Email = "invalid-email"
			}
			invalidUsers = append(invalidUsers, result.User)

		} else {
			validUsers = append(validUsers, result.User)
		}
	}

	s.validUsers = validUsers
	s.invalidUsers = invalidUsers
	return nil
}

func (s *Service) ValidCount() int {
	return len(s.validUsers)
}

func (s *Service) InvalidCount() int {
	return len(s.invalidUsers)
}

func (s *Service) GetValidUsers(page int, limit int, email string) ([]User, int) {

	var users []User

	if email != "" {
		for _, user := range s.validUsers {

			if user.Email == email {
				users = append(users, user)
			}
		}
	} else {
		users = s.validUsers
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(users) {
		users = []User{}
	} else if end > len(users) {
		users = users[start:]
	} else {
		users = users[start:end]
	}

	total := len(users)

	return users, total
}

func (s *Service) GetInvalidUsers() []User {
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

func isCSVHeader(records [][]string) error {
	if len(records) == 0 {
		return errors.New("file is empty and missing header")
	}

	if len(records[0]) != 4 {
		return errors.New("invalid header: expected 4 columns")
	}

	expected := [4]string{"id", "first_name", "last_name", "email"}
	for i := range expected {
		if records[0][i] != expected[i] {
			return errors.New("invalid header column name")
		}
	}
	return nil
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
