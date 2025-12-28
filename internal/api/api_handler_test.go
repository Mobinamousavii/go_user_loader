package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"goproject/internal/loader"
	"encoding/json"

)



func TestInvalidUserAPI(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "users/invalid", nil)

	if err!=nil{
		t.Error(err)
	}

	userlist := []loader.User{

	{ID: 1,
	FirstName: "Ali",
	LastName: "Ahmadi",
	Email: "ali@example.com",
	},
	{
		ID : -5,
		FirstName:"farnaz",
		LastName:"radaie",
		Email: "hjvv",
	},

	}




	expected_response := GetUsersResponse{
		Count: 1,
		Items: []loader.User{
			{ID: -5, FirstName: "farnaz", LastName: "radaie", Email: "invalid-email"},
		},
	}

	

	rr := httptest.NewRecorder()

	invaliduserhandler := InvalidUserHandler(userlist)
	handler := http.HandlerFunc(invaliduserhandler)
	
	handler.ServeHTTP(rr , req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rr.Code)
	}

	var response GetUsersResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.Count != expected_response.Count {
		t.Errorf("Expected count %d, got %d", expected_response.Count, response.Count)
	}


	if len(response.Items) != len(expected_response.Items) {
		t.Errorf("Expected %d items, got %d", len(expected_response.Items), len(response.Items))
	}

	


}