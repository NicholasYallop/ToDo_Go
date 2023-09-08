/*
Copyright © 2023 nick yallop
*/
package main

import (
	"ToDo/cmd"
	datastore "ToDo/data"
)

func main() {
	cmd.Execute()
	defer datastore.SaveCache()
}
