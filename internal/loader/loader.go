package loader

import (
	"errors"
	"log"
	"strings"

	)


func LoadFile(path string)  (int ,[]User ,error){

	
	var errnotsupported = errors.New("not Supported File")

	if strings.HasSuffix(path, ".csv"){
		users , err := LoadCSV(path)
		if err != nil{
			return 0 ,nil ,err
		}
		count := len(users)
		log.Printf("Loaded %b users", count)
		return count ,users ,nil

	}else if strings.HasSuffix(path, ".json"){
		users , err := LoadJSON(path)
		if err != nil{
			return 0 , nil,err
		}
		count := len(users)
		log.Printf("Loaded %b users", count)
		return  count ,users ,nil

	}else{
		return 0 ,nil ,errnotsupported
	}

}