package loader

import (
	"reflect"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	path := "/home/mobina-mousavi/go-practice/GoUserLoader/users.json"
	expected_user := []User{
		{
			ID:        1,
			FirstName: "Ali",
			LastName:  "Ahmadi",
			Email:     "ali@example.com",
		},
	}

	validuser,_, err := LoadJSON(path)
	if reflect.DeepEqual(expected_user, validuser) == false || err != nil {
		t.Error("sth went wrong", err)
	}

}
