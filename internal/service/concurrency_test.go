package service

import (
	"goproject/internal/loader"
	"testing"
)


func TestProcessStream(t *testing.T){
	records := make(chan loader.Record)

	// go func ()  {
	// 	defer close(records)
	// 	records <- loader.Record{Index: 0, Fields: []string{"id", "first_name", "last_name", "email"}}
	// 	records <- loader.Record{Index: 1, Fields: []string{"5", "person1", "last1", ""}}
	// 	records <- loader.Record{Index: 2, Fields: []string{"2", "person2", "last2", "person2@gmail.com"}}
	// 	records <- loader.Record{Index: 3, Fields: []string{"3", "person3", "last3", "person3@gmail.com"}}
	// }()

	go func() {
    defer close(records)

    records <- loader.Record{
        Index: 0,
        Fields: []string{`{"id":5,"first_name":"person1","last_name":"last1","email":""}`},
    }
    records <- loader.Record{
        Index: 1,
        Fields: []string{`{"id":2,"first_name":"person2","last_name":"last2","email":"person2@gmail.com"}`},
    }
    records <- loader.Record{
        Index: 2,
        Fields: []string{`{"id":3,"first_name":"person3","last_name":"last3","email":"person3@gmail.com"}`},
    }
	}()


	svc := NewService()

	if err := svc.ProcessStream(records, 2) ; err != nil{
		t.Fatalf("ProcessStream returned error: %v", err)
	}

	valid, _ := svc.GetValidUsers(1, 5, "")
	invalid := svc.GetInvalidUsers()

	if len(valid) != 2{
		// t.Fatalf("expected 2 valid users, got %d", total)
		t.Fatalf("validusers, got %+v", valid)
	}


	if valid[0].ID != 2 || valid[1].ID != 3{
		t.Fatalf("unexpected valid order: %+v", valid)
	}

	if len(invalid) != 1{
		t.Fatalf("expected 1 invalid user, got %d", len(invalid))
	}

	if invalid[0].Email != "invalid-email" {
		t.Fatalf("unexpected invalid user: %+v", invalid[0])
	}

}

func TestProcessStreamInvalidHeaderReturnsError(t *testing.T) {
	records := make(chan loader.Record)
	go func() {
		defer close(records)
		records <- loader.Record{Index: 0, Fields: []string{"id", "first", "last_name", "email"}}
		records <- loader.Record{Index: 1, Fields: []string{"1", "Ali", "Ahmadi", "ali@example.com"}}
	}()

	svc := NewService()
	if err := svc.ProcessStream(records, 1); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

