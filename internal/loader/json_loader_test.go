package loader

import (
	"reflect"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	//Will update to use relative paths or create temporary files 
	path := "/home/mobina-mousavi/Go_project/users.json"
	expected_validuser := []User{
		{
			ID:        1,
			FirstName: "Ali",
			LastName:  "Ahmadi",
			Email:     "ali@example.com",
		},
	}

	expected_invaliduser := []User{
		{
			ID:        -5,
			FirstName: "farnaz",
			LastName:  "radaie",
			Email:     "invalid-email",
		},
	}


	validuser, invaliduser,err := LoadJSON(path)

	if err != nil {
		t.Errorf("Failed to load CSV: %v", err)
		return
	}

	if !reflect.DeepEqual(expected_validuser, validuser) {
		t.Errorf("Expected valid users: %+v,  got: %+v", expected_validuser, validuser)
	}

	if !reflect.DeepEqual(expected_invaliduser, invaliduser) {
		t.Errorf("Expected invalid users: %+v,  got: %+v", expected_invaliduser, invaliduser)
	}

	}




