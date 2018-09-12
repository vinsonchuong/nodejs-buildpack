package main

import (
	"fmt"
	"os"
)

func main(){
	//detect, err := libbuildpackV3.DefaultDetect()
	//if err != nil {
	//	os.Exit(100)
	//}

	// TODO : pass detect somewhere to run logic
	fmt.Fprintf(os.Stdout, `nodejs = { version = "FIX ME" }`)
}