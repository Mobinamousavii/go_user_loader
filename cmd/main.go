package main

import (
	"flag"
	"fmt"
	"goproject/internal/api"
	"goproject/internal/loader"
	"log"
	"net/http"
	"os"
)

var (
	file string
	port int
)

func init() {
	flag.IntVar(&port, "port", 8080, "port the program listens on for incoming requests")
	flag.StringVar(&file, "file", "", "path to users file (.csv/.json)")
}
func main() {

	flag.Parse()

	if file == "" {
		fmt.Fprintln(os.Stderr, "-file is required for running the program")
		flag.Usage()
		os.Exit(1)
	}

	_, _, err := loader.LoadFile(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	userhandler := api.MakeGetHttpUsers(file)
	healthhandler := api.MakeHealthUsers(file)
	http.HandleFunc("/users", userhandler)
	http.HandleFunc("/health", healthhandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
