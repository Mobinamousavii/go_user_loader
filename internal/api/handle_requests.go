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

		countloader, itemsloader, err := loader.LoadFile(path)
		if err != nil {
			http.Error(w, "failed to load users", http.StatusInternalServerError)
			return
		}

		users := GetUsersResponse{Count: countloader, Items: itemsloader}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			log.Printf("failed to write response: %v", err)

		}

	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}

}
