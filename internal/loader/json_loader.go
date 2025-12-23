package loader

import (
	"encoding/json"
	"os"
)

func LoadJSON(path string) ([]User, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var userlist []User
	err = json.Unmarshal(data, &userlist)

	if err != nil {
		return nil, err
	}

	return userlist, nil

}
