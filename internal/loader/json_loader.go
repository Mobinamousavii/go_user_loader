package loader

import (
	"encoding/json"
	"os"
)

func LoadJSON(path string) ([]User, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var userlist []User
	err = json.Unmarshal(data, &userlist)

	if err != nil {
		return nil, err
	}

	return userlist, nil

}
