package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func LoadJSON(path string) (records [][]string, err error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", path, err)
	}

	defer file.Close()

	dec := json.NewDecoder(file)

	var items []json.RawMessage
	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("cannot decode json: %w", err)
	}

	var extra json.RawMessage
	if err := dec.Decode(&extra); err == nil {
		return nil, errors.New("invalid json: multiple top-level values")
	} else if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	var out [][]string
	for _, raw := range items {
		out = append(out, []string{string(raw)})
	}

	return out, nil

}
