package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)


type Record struct{
	Index int
	Fields []string
}


func LoadFile(path string) (<- chan Record, <- chan error) {
	records := make(chan Record)
	errCh:= make (chan error, 1)

	go func(){
		defer close(records)
		defer close(errCh)

		ext := strings.ToLower(filepath.Ext(path))

		switch ext{
		case ".csv":
			LoadCSV(path, records, errCh)
		case ".json":
			LoadJSON(path, records, errCh)
		default:
			errCh <- fmt.Errorf("unsupported file type: %q (only .csv or .json)", ext)
		}
	}()

	return records, errCh

}
