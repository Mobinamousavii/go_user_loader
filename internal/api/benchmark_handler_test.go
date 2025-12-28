package api

import (
	"bufio"
	"encoding/json"
	"goproject/internal/loader"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
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

		itemsloader, _,err := loader.LoadFile(path)
		if err != nil {
			http.Error(w, "failed to load users", http.StatusInternalServerError)
			return
		}

		users := GetUsersExample{Count: len(itemsloader), Items: itemsloader}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		erro := json.NewEncoder(w).Encode(users)
		if erro != nil {
			return  
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




func makeuser(n int)[]loader.User{

	userlist := make([]loader.User, n)

	for i:= 0; i < n ; i++{
		userlist[i] = loader.User{ID: i, FirstName: "User" + strconv.Itoa(i),
		LastName: "Test" + strconv.Itoa(i),
		Email: "user" + strconv.Itoa(i) + "@example.com"}
	}

	return  userlist
}

func writeUserJson( dir string, n int)(string ,error){

	path := filepath.Join(dir , "users.json")


	f, err := os.Create(path)
	if err!= nil{
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	w.WriteByte('[')
	enc := json.NewEncoder(w)


	for i := 1; i <= n; i++{

		if i > 1{
			w.WriteByte(',')

		}

		user := loader.User{
		ID: i,
		FirstName: "User" + strconv.Itoa(i),
		LastName: "Test" + strconv.Itoa(i),
		Email: "user" + strconv.Itoa(i) + "@example.com"}

		err := enc.Encode(user)
		if err!= nil{
			return "", err
		}
	}
	w.WriteByte(']')

	return path, nil

}


func BenchmarkUserGet(b *testing.B) {
	
	for n := 10000; n<=1000000; n += 10000{
		size := n
		
		b.Run("N="+strconv.Itoa(size), func(b *testing.B) {
			dir := b.TempDir()
			path, err := writeUserJson(dir,n)
			if err!=nil{
				b.Fatal(err)
			}
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
			
		})
	}
}
	 




func BenchmarkUserCacheGet(b *testing.B) {
	
	for n := 10000; n <= 1000000; n += 10000{
		size := n

		b.Run("N=" + strconv.Itoa(size), func(b *testing.B) {
			userlist := makeuser(n)
			
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
			
		})
	}
}