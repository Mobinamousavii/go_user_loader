package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func LoadJSON(path string , out chan <- Record, errCh chan <- error) {
	file, err := os.Open(path)

	if err != nil {
		errCh <- fmt.Errorf("cannot open file %q: %w", path, err)
		return
	}

	defer file.Close()

	dec := json.NewDecoder(file)


	t, err := dec.Token()
	if err!=nil{
		errCh <- fmt.Errorf("cannot decode json: %w", err)
		return
	}

	if t!= json.Delim('['){
		errCh <- errors.New("invalid json: expected top-level array")
		return

	}



	idx := 0
	
	for dec.More(){
		var raw json.RawMessage
		if err := dec.Decode(&raw); err!= nil{
			errCh <- fmt.Errorf("cannot decode json: %w", err)
			return
		}
		out <- Record{Index: idx, Fields: []string{string(raw)}}
		idx++
		
	}

	if _, err := dec.Token(); err != nil {
		errCh <- fmt.Errorf("cannot decode json: %w", err)
		return
	}

	var extra json.RawMessage
	if err := dec.Decode(&extra); err == nil {
		errCh <- errors.New("invalid json: multiple top-level values")
		return
	} else if !errors.Is(err, io.EOF) {
		errCh <- fmt.Errorf("invalid json: %w", err)
		return
	}

}
