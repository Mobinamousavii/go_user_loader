package loader

import (
	"goproject/internal/service"
	"reflect"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	//Will update to use relative paths or create temporary files
	path := "/home/mobina-mousavi/Go_project/users.csv"
	expected_userlist := []service.User{
		{
			ID:        1,
			FirstName: "Ali",
			LastName:  "Ahmadi",
			Email:     "ali@example.com",
		},
	}

	userlist, err := LoadCSV(path)

	if err != nil {
		t.Errorf("Failed to load CSV: %v", err)
		return
	}

	if !reflect.DeepEqual(expected_userlist, userlist) {
		t.Errorf("Expected userlist: %+v,  got: %+v", expected_userlist, userlist)
	}

}
