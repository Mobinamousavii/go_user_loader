package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

func LoadFile(path string) ([][]string, error) {

	lowerPath := strings.ToLower(path)

	if strings.HasSuffix(lowerPath, ".csv") {
		records, err := LoadCSV(path)
		if err != nil {
			return nil, err
		}
		return records, nil

	} else if strings.HasSuffix(lowerPath, ".json") {
		records, err := LoadJSON(path)
		if err != nil {
			return nil, err
		}
		return records, nil

	} else {
		got := filepath.Ext(path)
		err := fmt.Errorf("unsupported file extension for %q: got %q, expected .csv or .json", path, got)
		return nil, err
	}

}
