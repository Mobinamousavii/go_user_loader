package api

import (
	"encoding/json"
	"fmt"
	"goproject/internal/service"
	"log"
	"net/http"
	"strconv"

)

type GetUsersResponse struct {
	Count int            `json:"count"`
	Items []service.User `json:"items"`
}

type GetFilteredResponse struct{
	Page int				`json:"page"`
	Total int				`json:"total"`
	Limit int				`json:"limit"`
	Items []service.User	`json:"items"`

}

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct{
	Error string`json:"error"`
}

type Server struct {
	svc *service.Service
}

func New(svc *service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/users", s.Users)
	mux.HandleFunc("/users/invalid", s.handleInvalidUsers)
	return mux
}


func getPaginationParams(r *http.Request)(page int , limit int , email string , err error){
	page = 1

	if p:= r.URL.Query().Get("page"); p!=""{
		page , err = strconv.Atoi(p)
		if err!=nil{
			return 0, 0, "", fmt.Errorf("invalid page parameter")
		}
	}

	limit = 10 
	if l := r.URL.Query().Get("limit"); l!=""{
		limit ,err = strconv.Atoi(l)
		if err!=nil{
			return 0, 0, "", fmt.Errorf("invalid limit paramter")
		}
	}

	email = r.URL.Query().Get("email")

	return page, limit, email, nil

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

func (s *Server)Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case http.MethodGet:
		s.handleUsers(w ,r)
	case http.MethodPost:
		s.handleAddUser(w, r)

	default:
		http.Error(w, "method not allowed" , http.StatusMethodNotAllowed)
	}
}


func (s *Server)handleAddUser(w http.ResponseWriter, r *http.Request){
	var user service.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil{
		http.Error(w , "invalid json body", http.StatusBadRequest)
	}

	errs := s.svc.AddUser(user)

	if len(errs) > 0{
		var resp ErrorResponse
		for _, e := range errs{
			resp = ErrorResponse{Error: e.Message}
		}
		writeJSON(w, http.StatusBadRequest, resp)
		return
	}

	created := HealthResponse{Status: "created"}
	writeJSON(w, http.StatusCreated, created)


}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {

	page, limit , email , err := getPaginationParams(r)

	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	users, total := s.svc.GetValidUsers(page, limit,email)

	response := GetFilteredResponse{
		Page: page,
		Total: total,
		Limit: limit,
		Items: users,
	}

	writeJSON(w, http.StatusOK, response)

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
