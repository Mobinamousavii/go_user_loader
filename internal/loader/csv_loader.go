package loader

import (
	"encoding/csv"
	"os"
	"strconv"
)

func createListCsv(records [][]string) ([]User, error) {
	var userlist []User
	for i, line := range records {
		var rec User
		if i > 0 {
			for j, field := range line {
				if j == 0 {
					number, err := strconv.Atoi(field)
					if err != nil {
						return nil, err
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
