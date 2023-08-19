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

func init() {
}

func InitVariables(datastore_path string) {
	file, _ := os.OpenFile(datastore_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Reader = csv.NewReader(file)
	Writer = csv.NewWriter(file)
	fmt.Println(datastore_path)
}

func FetchAll() []structs.Task {
	records, err := Reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var returnValue []structs.Task
	for _, record := range records {
		task := structs.TaskFromCsv(record)
		returnValue = append(returnValue, task)
		if task.ID > GreatestID {
			GreatestID = task.ID
		}
	}

	return returnValue
}

func Store(task structs.Task) {
	Writer.Write(task.ToCsvLine())
	Writer.Flush()
	err := Writer.Error()
	if err != nil {
		fmt.Println(err)
	}
}
