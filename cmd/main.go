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

	validuser, invaliduser, err := loader.LoadFile(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	 fmt.Println(invaliduser)


	log.Printf("Loaded %d valid users", len(validuser))
	log.Printf("Skipped %d invalid users", len(invaliduser))

	userhandler := api.UserHandler(file)
	cacheuserhandler := api.CacheUsersHandler(validuser)
	invaliduserhandler := api.InvalidUserHandler(invaliduser)

	http.HandleFunc("/users-nocaching", userhandler)
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/users", cacheuserhandler)
	http.HandleFunc("/users/invalid", invaliduserhandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
