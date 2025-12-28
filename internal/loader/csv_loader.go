package loader

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"net/mail"
)

func createListCsv(records [][]string) ( validuser []User, invaliduser []User, err error) {
	var emptyFile = errors.New("file is empty and missing header")
	if len(records) == 0 {
		return nil, nil, emptyFile
	}

	
	var invalidHeaderCount = errors.New("invalid header: expected 4 columns")
	if len(records[0]) != 4 {
		return nil, nil, invalidHeaderCount
	}
	
	var invalidHeaderName = errors.New("invalid header column name")
	expected := [4]string{"id", "first_name", "last_name", "email"}
	for i := range expected {
		if records[0][i] != expected[i] {
			return nil, nil, invalidHeaderName
		}
	}
	
	var invalidColumnsCount = errors.New("invalid record: expected 4 columns")
	for i := range records[1:] {
		if len(records[i]) != 4 {
			return nil, nil, invalidColumnsCount
		}
	}

	invaliduser = make([]User,0)

	validuser = make([]User, 0, len(records)-1)

	type ValidationErrors struct{
		IDError        error
		EmailError     error
		FirstNameError error
	}
	
	for i, line := range records {
		var vErrs ValidationErrors
		var rec User
		if i > 0 {
			for j, field := range line {
				if j == 0 {
					number, err := strconv.Atoi(strings.TrimSpace(field))
					if err != nil {
						return nil, nil, err
					}
					if number <= 0 {
						vErrs.IDError = errors.New("id must be a positive integer")
					}
					rec.ID = number
				} else if j == 1 {
					if field == ""{
						vErrs.FirstNameError = errors.New("firstname is empty")

					}
					rec.FirstName = field
				} else if j == 2 {
					rec.LastName = field
				} else if j == 3 {
					email := strings.TrimSpace(field)
					rec.Email = email
					if email == ""{
						vErrs.EmailError = errors.New("email is empty")
						rec.Email = "no-email"
					}

					if vErrs.EmailError == nil{
						_, err := mail.ParseAddress(email)
						if err!= nil{
							vErrs.EmailError = errors.New("invalid email")
							rec.Email = "invalid-email"
						}
					}

				}
			}

			if vErrs.IDError != nil ||
				vErrs.FirstNameError != nil ||
				vErrs.EmailError != nil{
					invaliduser = append(invaliduser, rec)
					
				}else {
					validuser = append(validuser, rec)
				}
			
		}
	}
	return validuser, invaliduser, nil
}

func LoadCSV(path string) ([]User, []User,error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, nil, err
	}

	defer f.Close()

	filereader := csv.NewReader(f)
	records, err := filereader.ReadAll()

	if err != nil {
		return nil, nil, err

	}

	validuser, invaliduser, err := createListCsv(records)

	if err != nil {
		return nil, nil, err
	}

	return validuser, invaliduser, nil
}
