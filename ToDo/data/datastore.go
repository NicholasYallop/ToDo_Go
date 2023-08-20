package datastore

import (
	structs "ToDo/structs"
	"encoding/csv"
	"fmt"
	"os"
)

var Reader *csv.Reader
var Writer *csv.Writer
var GreatestID int = 0
var tasks_cache []structs.Task = nil

func init() {
}

func InitVariables(datastore_path string) {
	file, err := os.OpenFile(datastore_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open datastore.")
		panic(err)
	}
	Reader = csv.NewReader(file)
	Writer = csv.NewWriter(file)

	FetchAll()
}

func FetchAll() []structs.Task {
	if tasks_cache != nil {
		return tasks_cache
	}

	records, err := Reader.ReadAll()
	if err != nil {
		fmt.Println("Error while reading.")
		panic(err)
	}

	tasks_cache = []structs.Task{}
	for _, record := range records {
		task := structs.TaskFromCsv(record)
		tasks_cache = append(tasks_cache, task)
		if task.ID > GreatestID {
			GreatestID = task.ID
		}
	}

	return tasks_cache
}

func Store(task structs.Task) {
	GreatestID += 1
	task.ID = GreatestID
	Writer.Write(task.ToCsvLine())
	Writer.Flush()
	err := Writer.Error()
	if err != nil {
		GreatestID -= 1
		fmt.Println("Error while writing.")
		fmt.Println(err)
	} else {
		tasks_cache = nil
	}
}
