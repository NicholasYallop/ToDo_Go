/*
Copyright © 2023 nick yallop
*/
package main

import (
	"ToDo/cmd"
	datastore "ToDo/data"
	"fmt"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't get working directory.")
		panic(err)
	}
	datastore.InitVariables(path + "\\data\\data.csv")

	cmd.Execute()
}
