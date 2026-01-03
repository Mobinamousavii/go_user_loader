package api

import (
	"encoding/json"
	"goproject/internal/service"
	"log"
	"net/http"
)

type GetUsersResponse struct {
	Count int            `json:"count"`
	Items []service.User `json:"items"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type Server struct {
	svc *service.Service
}

func New(svc *service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /users", s.handleUsers)
	mux.HandleFunc("GET /users/invalid", s.handleInvalidUsers)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{Status: "ok"}
	writeJSON(w, http.StatusOK, health)

}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := s.svc.GetValidUsers()
	validcount := s.svc.ValidCount()

	users := GetUsersResponse{Count: validcount, Items: items}
	writeJSON(w, http.StatusOK, users)

}

func (s *Server) handleInvalidUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := s.svc.GetInvalidUsers()
	invalidcount := s.svc.InvalidCount()

	users := GetUsersResponse{Count: invalidcount, Items: items}
	writeJSON(w, http.StatusOK, users)

}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}

}
