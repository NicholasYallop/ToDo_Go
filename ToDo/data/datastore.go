package datastore

import (
	structs "ToDo/structs"
	"fmt"
	"os"
)

var tasks_cache structs.TaskSlice = nil
var datastore_path string

func init() {
	localappdata, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	datastore_path = localappdata + "/todo/data.json"
}

func AddToCache(task structs.Task) {
	if tasks_cache == nil {
		tasks_cache = structs.TaskSlice{}
	}
	tasks_cache = append(tasks_cache, task)
}

func SetCache(tasks structs.TaskSlice) {
	tasks_cache = tasks
}

func FetchAll() structs.TaskSlice {
	if tasks_cache != nil {
		return tasks_cache
	}

	content, err := os.ReadFile(datastore_path)
	if err != nil {
		tasks_cache = []structs.Task{{}}
		return tasks_cache
	}
	tasks := structs.TasksFromJson(content)
	if len(tasks) != 0 {
		tasks_cache = tasks
	} else {
		tasks_cache = []structs.Task{{}}
	}
	return tasks_cache
}

func SaveCache() {
	overwrite(tasks_cache)
}

func overwrite(tasks structs.TaskSlice) {
	err := os.WriteFile(datastore_path, tasks.ToJson(), 0666)
	if err != nil {
		fmt.Println("Error writing new datastore")
		panic(err)
	}
}
