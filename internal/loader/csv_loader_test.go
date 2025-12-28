package loader

import (
	"reflect"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	path := "/home/mobina-mousavi/Go_project/users.csv"
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
			ID: 2,
			FirstName: "Amir",
			LastName: "Mousavi",
			Email: "no-email",
		},

	}

	validuser, invaliduser,err := LoadCSV(path)

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
