package loader

import (
	"reflect"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	path := "/home/mobina-mousavi/go-practice/GoUserLoader/users.csv"
	expected_user := []User{
		{
			ID:        1,
			FirstName: "Ali",
			LastName:  "Ahmadi",
			Email:     "ali@example.com",
		},
	}

	actual_user, _,err := LoadCSV(path)

	if reflect.DeepEqual(expected_user, actual_user) == false || err != nil {
		t.Error("sth went wrong", err)

	}
}
