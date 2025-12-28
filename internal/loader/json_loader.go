package loader

import (
	"encoding/json"
	"net/mail"
	"os"
	"strings"
)

func LoadJSON(path string) (validuser []User,invaliduser []User, err error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, nil,err
	}


	var userlist []User

	err = json.Unmarshal(data, &userlist)

	if err != nil {
		return nil, nil, err
	}

	
	for i := range userlist{
		var validationerrors []ValidationError
		email := strings.TrimSpace(userlist[i].Email)
		if userlist[i].ID <= 0{
			validationerrors = append(validationerrors, *ErrIDInvalid)

		}
		
		if userlist[i].FirstName == ""{
			validationerrors = append(validationerrors, *ErrFirstNameEmpty)

		}
		
		_, err := mail.ParseAddress(email)
		if err != nil{
			validationerrors = append(validationerrors, *ErrEmailInvalid)
			userlist[i].Email = "invalid-email"
			}
		


		if len(validationerrors) != 0{
			invaliduser = append(invaliduser, userlist[i])
		}else{
			validuser = append(validuser, userlist[i])
		}
	}

	return validuser, invaliduser,nil

}
