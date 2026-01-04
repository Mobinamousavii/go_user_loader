package service

import (
	"goproject/internal/validation"
	"sync"
)

type Job struct{
	Index int
	Record []string
}


type Result struct{
	Index int
	User User
	ParseErr error
	ValidationErrs []*validation.ValidationError
}


func Worker(jobs <-chan Job,results chan<- Result, wg *sync.WaitGroup,strategy SetUsersStrategy){

	defer wg.Done()

	for job := range jobs{

		user , err := strategy.Parse(job.Record,job.Index)

		if err!=nil{
			results <- Result{Index: job.Index,ParseErr: err} 
			return
		}


		errs := ValidateUser(user)
		results <- Result{Index: job.Index, User: user, ValidationErrs: errs}
	}
}

func RunPool(records [][]string, workers int , strategy SetUsersStrategy)([]Result, error){

	jobs := make(chan Job , len(records))
	results := make(chan Result, len(records))

	var wg sync.WaitGroup

	for i , record := range records{
		jobs <-  Job{Index: i, Record: record}
	}

	close(jobs)

	for i:=0; i< workers; i++{
		wg.Add(1)
		go Worker(jobs, results, &wg , strategy)
	}

	wg.Wait()

	close(results)

	var finalResults []Result
	for result := range results {
		finalResults = append(finalResults, result)
	}


	return finalResults, nil

}
