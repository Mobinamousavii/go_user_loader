package main

import (
	"flag"
	"fmt"
	"goproject/internal/loader"
	"net/http"
	"goproject/internal/api"
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

	fmt.Println(loader.LoadFile(file))

	http.HandleFunc("/users", api.GetHttpUsers)
	http.ListenAndServe(":8080", nil)


	
}
