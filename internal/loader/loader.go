package loader


import (
	"errors"
	"log"
	"strings"
)


func LoadFile(path string)  (int ,error){

	var errnotsupported = errors.New("not Supported File")

	if strings.HasSuffix(path, ".csv"){
		users , err := LoadCSV(path)
		if err != nil{
			return 0 , err
		}
		count := len(users)
		log.Printf("Loaded %b users", count)
		return count, nil

	}else if strings.HasSuffix(path, ".json"){
		users , err := LoadJSON(path)
		if err != nil{
			return 0, err
		}
		count := len(users)
		log.Printf("Loaded %b users", count)
		return  count, nil

	}else{
		return 0,errnotsupported
	}

}