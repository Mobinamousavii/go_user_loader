package loader

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"fmt"
)



func createListCsv(records [][]string)( []User,  error){
	if len(records) == 0 {
		return nil, errors.New("file is empty and missing header")
	}

	if len(records[0]) != 4 {
		return nil, errors.New("invalid header: expected 4 columns")
	}
	
	expected := [4]string{"id", "first_name", "last_name", "email"}
	for i := range expected {
		if records[0][i] != expected[i] {
			return nil, errors.New("invalid header column name")
		}
	}
	

	userlist := make([]User, 0)

	for row, record:= range records[1:]{

		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid id at csv row %d: %w", row, err)
		}

		user := User{
			ID: id,
			FirstName: record[1],
			LastName: record[2],
			Email: record[3],
		} 

		userlist = append(userlist, user)
	}

	return userlist, nil


}

func LoadCSV(path string) ([]User, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil,fmt.Errorf("cannot open file %q: %w", path, err)
	}

	defer f.Close()

	filereader := csv.NewReader(f)
	records, err := filereader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("cannot read CSV file %q: %w", path, err)

	}

	userlist, err := createListCsv(records)

	if err != nil {
		return nil, err
	}

	return userlist, nil
}
