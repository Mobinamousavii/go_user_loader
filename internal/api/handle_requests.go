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


func MakeGetHttpUers(path string) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){

	countloader , itemsloader , err := loader.LoadFile(path)
	if err != nil{
		http.Error(w,"sth went wrong", http.StatusBadRequest)
		return
	} 

	Users := Getusers{Count: countloader, Items: itemsloader}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Users)

	}
}


