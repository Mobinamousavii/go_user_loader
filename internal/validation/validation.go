package validation

import (
	"fmt"
	"net/mail"
)

type ValidationError struct {
	Field   string
	Message string
	Code    int
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Message: %s, Code: %d", v.Field, v.Message, v.Code)
}

var (
	ErrIDInvalid      = &ValidationError{Field: "id", Message: "id must be a positive integer", Code: 1001}
	ErrEmailInvalid   = &ValidationError{Field: "email", Message: "invalid email format", Code: 1002}
	ErrFirstNameEmpty = &ValidationError{Field: "first_name", Message: "first name cannot be empty", Code: 1003}
)

func ValidID(id int) error {
	if id < 0 {
		return ErrIDInvalid
	}
	return nil
}

func ValidFirstname(firstname string) error {
	if firstname == "" {
		return ErrFirstNameEmpty
	}
	return nil
}

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}
	return nil

}
