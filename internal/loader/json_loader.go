package loader

import (
	"encoding/json"
	"errors"
	"net/mail"
	"os"
	"strings"
)

func LoadJSON(path string) (validuser []User,invaliduser []User, err error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, nil,err
	}


	type ValidationErrors struct{
		IDError        error
		EmailError     error
		FirstNameError error
	}
	
	var vErrs ValidationErrors
	var userlist []User

	err = json.Unmarshal(data, &userlist)

	if err != nil {
		return nil, nil, err
	}

	
	for i := range userlist{
		email := strings.TrimSpace(userlist[i].Email)
		if userlist[i].ID <= 0{
			vErrs.IDError = errors.New("id must be a positive integer")

		}else if userlist[i].FirstName == ""{
			vErrs.FirstNameError = errors.New("firstname is empty")

		}else if email == ""{
			vErrs.EmailError = errors.New("email is empty")
			userlist[i].Email = "no_email"
			

		}else if vErrs.EmailError == nil{
			_,err := mail.ParseAddress(email)
			if err!=nil{
				vErrs.EmailError = errors.New("invalid email")
				userlist[i].Email = "invalid-email"
				
			}

		}


		if vErrs.IDError != nil ||
			vErrs.FirstNameError != nil ||
			vErrs.EmailError != nil {
				invaliduser = append(invaliduser,userlist[i])
			}else{
				validuser = append(validuser, userlist[i])
			}
	}

	return validuser, invaliduser,nil

}
