package structs

import "strconv"

type Task struct {
	ID          int
	Name        string
	Description string
}

func (x *Task) ToCsvLine() []string {
	return []string{strconv.Itoa(x.ID), x.Name, x.Description}
}

func TaskFromCsv(line []string) Task {
	id, _ := strconv.Atoi(line[0])
	return Task{
		ID:          id,
		Name:        line[1],
		Description: line[2],
	}
}
