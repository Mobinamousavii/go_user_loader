package loader

import (
	"reflect"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	path := "/home/mobina-mousavi/go-practice/GoUserLoader/users.csv"
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

	if reflect.DeepEqual(expected_validuser, validuser) == false ||
	reflect.DeepEqual(expected_invaliduser,invaliduser) == false||
	err != nil {

		t.Error("sth went wrong", err)

	}
}
