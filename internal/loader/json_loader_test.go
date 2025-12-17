package loader

import (
	"testing"
	"reflect"

)

func TestLoadJSON(t *testing.T) {
	path := "/home/mobina-mousavi/go-practice/GoUserLoader/users.json"
	expected_user := []User{
		{
		ID : 1,
		FirstName : "Ali",
		LastName : "Ahmadi",
		Email : "ali@example.com",
		},
	}

	actual_user , err := LoadJSON(path)
	if reflect.DeepEqual(expected_user, actual_user) == false || err != nil{
		t.Error("sth went wrong" , err)
	}

}
