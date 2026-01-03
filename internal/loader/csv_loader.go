package loader

import (
	"encoding/csv"
	// "errors"
	"fmt"
	"os"
	// "strconv"
	// "strings"
)

func LoadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", path, err)
	}

	defer f.Close()

	filereader := csv.NewReader(f)
	records, err := filereader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("cannot read CSV file %q: %w", path, err)

	}

	return records, nil
}
