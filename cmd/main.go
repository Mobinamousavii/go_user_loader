package main

import (
	"fmt"
	"goproject/internal/loader"

)

func main(){
	file := "/home/mobina-mousavi/Go_project/users.csv"

	fmt.Println(loader.LoadFile(file))
}
