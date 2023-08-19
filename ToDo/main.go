/*
Copyright © 2023 nick yallop
*/
package main

import (
	"ToDo/cmd"
	datastore "ToDo/data"
	"os"
)

func main() {
	path, err := os.Getwd()
	datastore.InitVariables(path + "\\data\\data.csv")
	if err != nil {
		panic(err)
	}
	cmd.Execute()
}
