package structs

import (
	style "ToDo/defs"
	"fmt"
	"os"
	"strconv"

	"github.com/eiannone/keyboard"
	"github.com/nerdmaster/terminal"
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

func (menu Menu) Print() (cursorRow int) {
	_, height, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error while printing.")
		panic(err)
	}

	printStartIndex := min(menu.Focused_Task_Id, max(len(menu.Tasks)*2-height, 0))
	for index, task := range menu.Tasks {
		if index >= printStartIndex && index <= printStartIndex+(height/2-height%2) {
			if index == menu.Focused_Task_Id {
				fmt.Println(style.Default_Style.Render("> " + task.Name))
			} else {
				fmt.Println(style.Default_Style.Render(task.Name))
			}
			fmt.Println(style.Remark_Style.Render(task.Description))
		}
	}

	// move cursor up to focused Id
	fmt.Printf("\033[%dA", min((len(menu.Tasks))*2, height)-printStartIndex*2)
	return 2 * printStartIndex
}

func (x *Menu) MoveCursorUp() (success bool) {
	if x.Focused_Task_Id > 0 {
		x.Focused_Task_Id--
		return true
	}
	return false
}

func (x *Menu) MoveCursorDown() (success bool) {
	if x.Focused_Task_Id < len(x.Tasks)-1 {
		x.Focused_Task_Id++
		return true
	}
	return false
}

func (menu *Menu) ResetDisplayCursor(cursor int) (success bool) {
	if cursor != 0 {
		fmt.Printf("\033[%dA", cursor)
		return true
	}
	return false
}

func (menu *Menu) Display() error {
	keys, err := keyboard.GetKeys(1)
	if err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		cursor := menu.Print()

		keyEvent := <-keys
		if keyEvent.Err != nil {
			return err
		}

		menu.ResetDisplayCursor(cursor)

		if menu.HandleKeyInput(keyEvent.Key) {
			return nil
		}
	}
}

func (menu *Menu) HandleKeyInput(key keyboard.Key) (kill bool) {
	switch key {
	case keyboard.KeyArrowUp:
		menu.MoveCursorUp()
	case keyboard.KeyArrowDown:
		menu.MoveCursorDown()
	case keyboard.KeyEsc:
		fmt.Printf("\033[%dB", 2*len(menu.Tasks))
		return true
	}
	return false
}

//endregion Menu
