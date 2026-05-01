package config

import (
	"fmt"
	"flag"
	"errors"
)

type Config  struct{
	Filepath string
	Port int
	Workers int
}

func Parse(args []string)(Config , error){
	fs := flag.NewFlagSet("go-project", flag.ContinueOnError)

	var cfg Config

	fs.IntVar(&cfg.Port, "port", 8080, "port the program listens on for incoming requests")
	fs.StringVar(&cfg.Filepath, "file", "", "path to users file (.csv/.json)")
	fs.IntVar(&cfg.Workers, "workers", 5, "Number of workers to process records concurrently")
	

	if err := fs.Parse(args); err != nil{
		return Config{}, err
	}
	
	if cfg.Filepath == ""{
		return Config{}, errors.New("--file is required for running the program")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("invalid --port: %d", cfg.Port)
	}

	if cfg.Workers <= 0 {
		return  Config{}, fmt.Errorf("invalid value for 'workers': %d. The number of workers must be greater than or equal to 1.", cfg.Workers)
	}


	return cfg , nil
}