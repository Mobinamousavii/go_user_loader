package loader

import (
	"reflect"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	//Will update to use relative paths or create temporary files in the future
	path := "/home/mobina-mousavi/Go_project/users.json"
	expected_userlist := []User{
		{
			ID:        1,
			FirstName: "Ali",
			LastName:  "Ahmadi",
			Email:     "ali@example.com",
		},
	}



	userlist,err := LoadJSON(path)

	if err != nil {
		t.Errorf("Failed to load CSV: %v", err)
		return
	}

	if !reflect.DeepEqual(expected_userlist, userlist) {
		t.Errorf("Expected userlist: %+v,  got: %+v", expected_userlist, userlist)
	}


	}




