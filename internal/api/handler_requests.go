package api

import (
	"encoding/json"
	"goproject/internal/loader"
	"log"
	"net/http"

)

type GetUsersResponse struct {
	Count int           `json:"count"`
	Items []loader.User `json:"items"`
}

type HealthResponse struct {
	Status string `json:"status"`
}


func UserHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		itemsloader, _,err := loader.LoadFile(path)
		if err != nil {
			http.Error(w, "failed to load users", http.StatusInternalServerError)
			return
		}

		users := GetUsersResponse{Count: len(itemsloader), Items: itemsloader}
		writeJSON(w,http.StatusOK,users)


	}
}




func CacheUsersHandler(validuser []loader.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method!=http.MethodGet{
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w,"method not allowed", http.StatusMethodNotAllowed)
			return 
		}

		users := GetUsersResponse{Count: len(validuser), Items: validuser}
		writeJSON(w,http.StatusOK,users)
	}

}


func InvalidUserHandler(invaliduser []loader.User) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method!=http.MethodGet{
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w,"method not allowed", http.StatusMethodNotAllowed)
			return 
		}

		users := GetUsersResponse{Count: len(invaliduser), Items: invaliduser}
		writeJSON(w,http.StatusOK,users)
	}

}


func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{Status: "ok"}
	writeJSON(w,http.StatusOK,health)
	
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
	
}

