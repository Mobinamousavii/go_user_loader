package loader

import (
	"encoding/json"
	"fmt"
	"os"

)


func LoadJSON(path string)(userlist []User, err error){
	data , err := os.Open(path)

	if err!= nil{
		return nil, fmt.Errorf("cannot open file %q: %w", path, err)
	}

	dec := json.NewDecoder(data)

	
	if err := dec.Decode(&userlist) ; err!=nil{
		return nil, fmt.Errorf("cannot decode json: %w", err)
	}

	return userlist, nil


}

