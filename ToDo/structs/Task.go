package structs

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Name        string
	Description string
	SubTasks    TaskSlice
}

type TaskSlice []Task

func (tasks TaskSlice) ToJson() []byte {
	result, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("err while converting task array to byte array")
		panic(err)
	}
	return result
}

func TasksFromJson(bytes []byte) TaskSlice {
	var tasks TaskSlice
	err := json.Unmarshal(bytes, &tasks)
	if err != nil {
		fmt.Println("errr while converting string to task array")
		panic(err)
	}
	return tasks
}
