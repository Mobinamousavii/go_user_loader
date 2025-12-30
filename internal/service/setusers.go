package service

import (
	"goproject/internal/loader"
	"goproject/internal/validation"
)



func SetValidUsers(userlist []loader.User)(validuser []loader.User){
	for _, user := range userlist{
		validationerror:= validation.ValidateUser(user)

		if len(validationerror) == 0{
			validuser = append(validuser, user)
		} 	
		 
	}
	return validuser 
	
		
}

func SetInvalidUsers(userlist []loader.User)(invaliduser []loader.User){

	for _, user := range userlist{
		validationerror:= validation.ValidateUser(user)

		if len(validationerror) > 0 {
			for _, vErr := range validationerror{
				if vErr == *validation.ErrEmailInvalid{
					user.Email = "invalid-email"
				}
			}
			invaliduser = append(invaliduser, user)
		}

	}
	return invaliduser
}


