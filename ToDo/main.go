/*
Copyright © 2023 nick yallop
*/
package main

import (
	datastore "ToDo/data"
	"ToDo/structs"
)

func main() {
	tasks := datastore.FetchAll()
	defer datastore.SaveCache()

	menu := structs.NewMenu(tasks)

	go menu.Display()

	for {
		tasks := <-menu.OutputChannel
		if tasks != nil {
			datastore.SetCache(tasks)
		} else {
			break
		}
	}
}
