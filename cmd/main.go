package main

import (
	"fmt"
	"goproject/internal/app"
	"os"
)





func main() {

	if err := app.Run(os.Args[1:]); err != nil{
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}
