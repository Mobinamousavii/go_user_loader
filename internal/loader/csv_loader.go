package loader

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
)

func createListCsv(records [][]string) ([]User, error) {
	var emptyFile = errors.New("file is empty and missing header")
	if len(records) == 0 {
		return nil, emptyFile
	}
	userlist := make([]User, 0, len(records) - 1)

	var invalidHeaderCount = errors.New("invalid header: expected 4 columns")
	if len(records[0])!= 4{
		return nil, invalidHeaderCount
	}
	
	var invalidHeaderName = errors.New("invalid header column name")
	expected := [4]string{"id", "first_name", "last_name", "email"}
	for i := range expected{
		if records[0][i] != expected[i]{
			return nil , invalidHeaderName
		}
	}
	
	var invalidColumnsCount = errors.New("invalid record: expected 4 columns")
	for i := range records[1:]{
		if len(records[i])!=4{
			return nil , invalidColumnsCount
		}
	}

	var invalidID = errors.New("id must be a positive integer")
	for i, line := range records {
		var rec User
		if i > 0 {
			for j, field := range line{
				if j == 0 {
					number, err := strconv.Atoi(strings.TrimSpace(field))
					if err != nil {
						return nil, err
					}
					if number <= 0{
						return nil, invalidID
					}
					rec.ID = number
				} else if j == 1 {
					rec.FirstName = field
				} else if j == 2 {
					rec.LastName = field
				} else if j == 3 {
					rec.Email = field
				}
			}
			userlist = append(userlist, rec)
		}
	}
	return userlist, nil
}

func LoadCSV(path string) ([]User, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	filereader := csv.NewReader(f)
	records, err := filereader.ReadAll()

	if err != nil {
		return nil, err

	}

	userlist, err := createListCsv(records)

	if err != nil {
		return nil, err
	}

	return userlist, nil
}
