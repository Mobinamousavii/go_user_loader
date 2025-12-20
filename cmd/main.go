package main

import (
	"flag"
	"fmt"
	"goproject/internal/api"
	"goproject/internal/loader"
	"net/http"
	"os"
)


var (
	file string
	port int
)

func init(){
	flag.IntVar(&port, "port", 8080, "the port usage for api")
	flag.StringVar(&file, "file", "", "the file path")
}
func main(){

	flag.Parse()

	fmt.Println("file:", file)
	fmt.Println("port:", port)

	_,_,err := loader.LoadFile(file)
	if err!= nil{
		fmt.Println(err)
		os.Exit(1)
	}

	userhandler := api.MakeGetHttpUers(file)
	healthhandler := api.MakeHealthUsers(file)
	http.HandleFunc("/users", userhandler)
	http.HandleFunc("/health",healthhandler)
	http.ListenAndServe(":8080", nil)

	
}
