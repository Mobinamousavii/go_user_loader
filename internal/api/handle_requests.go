package api

import (
	"encoding/json"
	"goproject/internal/loader"
	"net/http"
)

type Getusers struct {
	Count int           `json:"count"`
	Items []loader.User `json:"items"`
}

type health struct {
	status string `josn:"status"`
}

func MakeGetHttpUers(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		countloader, itemsloader, err := loader.LoadFile(path)
		if err != nil {
			http.Error(w, "sth went wrong", http.StatusBadRequest)
			return
		}

		Users := Getusers{Count: countloader, Items: itemsloader}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Users)

	}
}

func MakeHealthUsers(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// _, _, err := loader.LoadFile(path)
		// if err!= nil{
		// 	http.Error(w,"sth went wrong", http.StatusBadRequest)
		// 	return
		// }

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(health{status: "ok"})

	}
}
