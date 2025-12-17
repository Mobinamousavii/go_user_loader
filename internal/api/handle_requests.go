package api

import (
	"encoding/json"
	"net/http"
	"goproject/internal/loader"
)




type Getusers struct{

	Count int           `json:"count"`
	Items []loader.User `json:"items"`

}



func GetHttpUsers(w http.ResponseWriter, r *http.Request){
	path := "/home/mobina-mousavi/Go_project"
	countloader , itemsloader , err := loader.LoadFile(path)
	if err != nil{
		http.Error(w,"sth went wrong", http.StatusBadRequest)
		return
	} 

	Users := Getusers{Count: countloader, Items: itemsloader}
	w.Header().Set("content-type", "applications/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Users)

}


