package structs

import (
	style "ToDo/defs"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/nerdmaster/terminal"
)

type EditingField int64

const (
	Tasks EditingField = iota
	Descriptions
)

type Menu struct {
	Tasks           []Task
	Focused_Task_Id int
	Editing         EditingField
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
				if menu.Editing == Tasks {
					fmt.Println(style.Default_Style.Render("> " + task.Name))
					fmt.Println(style.Remark_Style.Render(task.Description))
				} else if menu.Editing == Descriptions {
					fmt.Println(style.Default_Style.Render(task.Name))
					fmt.Println(style.Remark_Style.Render("> " + task.Description))
				}
			} else {
				fmt.Println(style.Default_Style.Render(task.Name))
				fmt.Println(style.Remark_Style.Render(task.Description))
			}
		}
	}

	// move cursor up to focused Id
	fmt.Printf("\033[%dA", min((len(menu.Tasks))*2, height)-(menu.Focused_Task_Id-printStartIndex)*2)
	return (menu.Focused_Task_Id - printStartIndex) * 2
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

func (menu *Menu) UpKeyPressed() {
	switch menu.Editing {
	case Tasks:
		menu.MoveCursorUp()
	}
}

func (menu *Menu) DownKeyPressed() {
	switch menu.Editing {
	case Tasks:
		menu.MoveCursorDown()
	}
}

func (menu *Menu) EnterKeyPressed() {
	switch menu.Editing {
	case Tasks:
		menu.Editing = Descriptions
	case Descriptions:
		menu.Editing = Tasks
	}
}

func (menu *Menu) EscKeyPressed() (kill bool) {
	switch menu.Editing {
	case Tasks:
		fmt.Printf("\033[%dB", 2*len(menu.Tasks))
		return true
	case Descriptions:
		menu.Editing = Tasks
	}
	return false
}

func (menu *Menu) DefaultKeyPress(event keyboard.KeyEvent) {
	switch menu.Editing {
	case Descriptions:
		if event.Key == keyboard.KeyBackspace {
			desc := menu.Tasks[menu.Focused_Task_Id].Description
			menu.Tasks[menu.Focused_Task_Id].Description = desc[0:max(len(desc)-1, 0)]
			return
		}
		menu.Tasks[menu.Focused_Task_Id].Description += string(event.Rune)
	}
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

		if menu.HandleKeyInput(keyEvent) {
			return nil
		}
	}
}

func (menu *Menu) HandleKeyInput(event keyboard.KeyEvent) (kill bool) {
	switch event.Key {
	case keyboard.KeyArrowUp:
		menu.UpKeyPressed()
	case keyboard.KeyArrowDown:
		menu.DownKeyPressed()
	case keyboard.KeyEnter:
		menu.EnterKeyPressed()
	case keyboard.KeyEsc:
		return menu.EscKeyPressed()
	default:
		menu.DefaultKeyPress(event)
	}
	return false
}
