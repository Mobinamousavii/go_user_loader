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
			handler := UserHandler(path)
			
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
			
			handler := CacheUsersHandler(userlist)
			
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