package api

import (
	"encoding/json"
	"goproject/internal/loader"
	"net/http"
	"net/http/httptest"
	"testing"
)


type GetUsersExample struct {
	Count int           `json:"count"`
	Items []loader.User `json:"items"`
}



func exampleUserHandler(path string)http.HandlerFunc{
	return func(w  http.ResponseWriter, r *http.Request){
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

		users := GetUsersExample{Count: countloader, Items: itemsloader}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		erro := json.NewEncoder(w).Encode(users)
		if erro != nil {
			return  
		}

	}
}


func BenchmarkUserGet(b *testing.B) {
	path := "/home/mobina-mousavi/Go_project/tests/testusers.csv"
	handler := exampleUserHandler(path)
	
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i <b.N ; i++{
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK{
			b.Fatalf("unexpected status code: %d", rr.Code)
		}
	}
}




func exampleUserCacheHandler(userlist []loader.User) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method!=http.MethodGet{
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w,"method not allowed", http.StatusMethodNotAllowed)
			return 
		}

		users := GetUsersExample{Count: len(userlist), Items: userlist}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(users)
		if err != nil {
			return 
		}

	}
}

func BenchmarkUserCacheGet(b *testing.B) {
	path := "/home/mobina-mousavi/Go_project/tests/testusers.csv"
	_,userlist,_ := loader.LoadFile(path)

	
	handler := exampleUserCacheHandler(userlist)
	
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i <b.N ; i++{
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK{
			b.Fatalf("unexpected status code: %d", rr.Code)
		}
	}

}