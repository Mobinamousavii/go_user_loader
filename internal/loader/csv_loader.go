package loader

import (
	"encoding/csv"
	"io"
	"fmt"
	"os"

)

func LoadCSV(path string , out chan <- Record, errCh chan <- error){
	f, err := os.Open(path)

	if err != nil {
		errCh <- fmt.Errorf("cannot open file %q: %w", path, err)
		return
	}

	defer f.Close()

	filereader := csv.NewReader(f)

	idx := 0
	for{
		rec, err := filereader.Read()
		if err == io.EOF{
			return
		}
		if err!= nil{
			errCh <- fmt.Errorf("cannot read CSV file %q: %w", path, err)
			return
		}

		out <- Record{Index: idx, Fields: rec}
		idx++

	}
}
