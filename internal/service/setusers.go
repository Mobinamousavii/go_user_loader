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
	Parse(records [][]string) ([]User, error)
}

func DetectStrategy(records [][]string) SetUsersStrategy {

	if err := isCSVHeader(records); err != nil {
		return JSONStrategy{}
	}

	return CSVStrategy{}

}

func (s *Service) SetUsers(records [][]string) error {

	if len(records) == 0 {
		s.validUsers = []User{}
		s.invalidUsers = []User{}
		return nil
	}

	strategy := DetectStrategy(records)

	users, err := strategy.Parse(records)

	if err != nil {
		return err
	}

	var validUsers []User
	var invalidUsers []User

	for _, u := range users {
		errs := ValidateUser(u)

		if len(errs) == 0 {
			validUsers = append(validUsers, u)
			continue
		}

		if hasEmailError(errs) {
			u.Email = "invalid-email"
		}

		invalidUsers = append(invalidUsers, u)
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

func (s *Service) GetValidUsers() []User {
	out := make([]User, len(s.validUsers))
	copy(out, s.validUsers)
	return out
}

func (s *Service) GetInvalidUsers() []User {
	out := make([]User, len(s.invalidUsers))
	copy(out, s.invalidUsers)
	return out
}

type CSVStrategy struct{}

type JSONStrategy struct{}

func (CSVStrategy) Parse(records [][]string) ([]User, error) {
	var userlist []User

	for row, record := range records[1:] {

		if len(record) != 4 {
			return nil, fmt.Errorf("invalid csv row %d: expected 4 columns, got %d", row+2, len(record))
		}

		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid id at csv row %d: %w", row+2, err)
		}

		user := User{
			ID:        id,
			FirstName: strings.TrimSpace(record[1]),
			LastName:  strings.TrimSpace(record[2]),
			Email:     strings.TrimSpace(record[3]),
		}

		userlist = append(userlist, user)
	}
	return userlist, nil
}

func (JSONStrategy) Parse(records [][]string) ([]User, error) {
	var userlist []User

	for row, record := range records {
		if len(record) != 1 {
			return nil, fmt.Errorf("json: row %d must have exactly 1 column (raw object), got %d", row+1, len(record))
		}

		raw := strings.TrimSpace(record[0])
		if raw == "" {
			return nil, fmt.Errorf("json: empty object at row %d", row+1)
		}

		var obj map[string]any

		if err := json.Unmarshal([]byte(raw), &obj); err != nil {
			return nil, fmt.Errorf("invalid json at row %d: %w", row+1, err)
		}

		if _, ok := obj["id"]; !ok {
			return nil, fmt.Errorf("json row %d: missing key 'id'", row+1)
		}
		if _, ok := obj["first_name"]; !ok {
			return nil, fmt.Errorf("json row %d: missing key 'first_name'", row+1)
		}
		if _, ok := obj["last_name"]; !ok {
			return nil, fmt.Errorf("json row %d: missing key 'last_name'", row+1)
		}
		if _, ok := obj["email"]; !ok {
			return nil, fmt.Errorf("json row %d: missing key 'email'", row+1)
		}

		id, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(obj["id"])))
		if err != nil {
			return nil, fmt.Errorf("json row %d: invalid id: %w", row+1, err)
		}

		first := fmt.Sprint(obj["first_name"])
		last := fmt.Sprint(obj["last_name"])
		email := fmt.Sprint(obj["email"])

		user := User{
			ID:        id,
			FirstName: strings.TrimSpace(first),
			LastName:  strings.TrimSpace(last),
			Email:     strings.TrimSpace(email),
		}

		userlist = append(userlist, user)

	}
	return userlist, nil
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
