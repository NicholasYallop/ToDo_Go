package structs

import (
	style "ToDo/defs"
	"fmt"
	"strconv"
)

// region Task
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

//endregion Task

// region Menu
type Menu struct {
	Tasks           []Task
	Focused_Task_Id int
}

func (x Menu) Print() {
	for _, task := range x.Tasks {
		if task.ID == x.Focused_Task_Id {
			blinker := style.Default_Style.Copy()
			fmt.Println(blinker.Render("> " + task.Name))
		} else {
			fmt.Println(style.Default_Style.Render(task.Name))
		}
		fmt.Println(style.Remark_Style.Render(task.Description))
	}

}

func (x *Menu) MoveCursorUp() bool {
	if x.Focused_Task_Id > 0 {
		x.Focused_Task_Id--
		return true
	}
	return false
}

func (x *Menu) MoveCursorDown() bool {
	if x.Focused_Task_Id < len(x.Tasks)-1 {
		x.Focused_Task_Id++
		return true
	}
	return false
}

//endregion Menu
