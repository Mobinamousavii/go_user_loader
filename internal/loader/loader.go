package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

func LoadFile(path string) (int, []User, error) {

	lowerPath := strings.ToLower(path)

	if strings.HasSuffix(lowerPath, ".csv") {
		users, err := LoadCSV(path)
		if err != nil {
			return 0, nil, err
		}
		count := len(users)
		return count, users, nil

	} else if strings.HasSuffix(lowerPath, ".json") {
		users, err := LoadJSON(path)
		if err != nil {
			return 0, nil, err
		}
		count := len(users)
		return count, users, nil

	} else {
		got := filepath.Ext(path)
		err := fmt.Errorf("unsupported file extension for %q: got %q, expected .csv or .json", path, got)
		return 0, nil, err
	}

}
