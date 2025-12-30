package validation

import (
	"fmt"
	"goproject/internal/loader"
	"net/mail"
)


type ValidationError struct{
	Field  string
	Message string
	Code int
}

func (v *ValidationError) Error() string{
	return fmt.Sprintf("Field: %s, Message: %s, Code: %d", v.Field, v.Message, v.Code)
}


var (
    ErrIDInvalid      = &ValidationError{Field: "id", Message: "id must be a positive integer", Code: 1001}
    ErrEmailInvalid   = &ValidationError{Field: "email", Message: "invalid email format", Code: 1002}
    ErrFirstNameEmpty = &ValidationError{Field: "first_name", Message: "first name cannot be empty", Code: 1003}
)



func ValidateUser(user loader.User)(validationerror []ValidationError){

	if user.ID < 0{
		validationerror = append(validationerror, *ErrIDInvalid)
	}

	if user.FirstName == ""{
		validationerror = append(validationerror, *ErrFirstNameEmpty)
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		validationerror = append(validationerror, *ErrEmailInvalid)
	}

	return validationerror

}