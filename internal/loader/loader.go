package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

func LoadFile(path string) ([]User, []User,error) {

	lowerPath := strings.ToLower(path)

	if strings.HasSuffix(lowerPath, ".csv") {
		users, invaliduser,err := LoadCSV(path)
		if err != nil {
			return nil, nil,err
		}
		return users, invaliduser,nil

	} else if strings.HasSuffix(lowerPath, ".json") {
		users, err := LoadJSON(path)
		if err != nil {
			return nil, nil,err
		}
		return users, nil,nil

	} else {
		got := filepath.Ext(path)
		err := fmt.Errorf("unsupported file extension for %q: got %q, expected .csv or .json", path, got)
		return nil, nil,err
	}

}
