package structs

import (
	style "ToDo/defs"
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/nerdmaster/terminal"
)

type EditingField int64

const (
	Tasks EditingField = iota
	Descriptions
)

type Menu struct {
	Tasks         []Task
	Focused_Index int
	Editing       EditingField
	OutputChannel chan []Task
}

type TaskStrings struct {
	taskString string
	descString string
}

func (menu *Menu) TaskStrings(index int) (taskStrings TaskStrings) {
	task := menu.Tasks[index]

	var taskString string
	var descString string

	if index == menu.Focused_Index {
		if menu.Editing == Tasks {
			taskString = style.Default_Style.Render("> " + task.Name)
			descString = style.Remark_Style.Render(task.Description)
		}
		if menu.Editing == Descriptions {
			taskString = style.Default_Style.Render(task.Name)
			descString = style.Remark_Style.Render("> " + task.Description)
		}
	} else {
		taskString = style.Default_Style.Render(task.Name)
		descString = style.Remark_Style.Render(task.Description)
	}

	return TaskStrings{
		taskString: taskString,
		descString: descString,
	}
}

func (menu Menu) Print() (cursorRow int) {
	_, terminalHeight, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error while printing.")
		panic(err)
	}

	// print from startindex, height of at most terminal height
	stringBuffer := ""
	focusedLine := 0
	currentLine := 0

	for i := 0; i < len(menu.Tasks); i++ {
		printData := menu.TaskStrings(i)
		if i == menu.Focused_Index && menu.Editing == Tasks {
			focusedLine = currentLine
		}

		for _, line := range strings.Split(printData.taskString, "\n") {
			stringBuffer += line + "\n"
			currentLine++
		}

		if i == menu.Focused_Index && menu.Editing == Descriptions {
			focusedLine = currentLine
		}
		for _, line := range strings.Split(printData.descString, "\n") {
			stringBuffer += line + "\n"
			currentLine++
		}
	}

	printHeight := 0
	cursorReset := 0
	if terminalHeight >= currentLine {
		// print all lines
		fmt.Print(strings.Trim(stringBuffer, "\n"))
		printHeight = currentLine - 1
		cursorReset = focusedLine

	} else if focusedLine <= currentLine-terminalHeight {
		// print from focused up to terminal height
		lines := strings.Split(stringBuffer, "\n")

		i := 0
		for i < terminalHeight-1 {
			fmt.Println(lines[focusedLine+i])
			i++
		}
		fmt.Print(lines[i])
		printHeight = terminalHeight
		cursorReset = 0

	} else {
		// print from end-terminalHeight up to end
		lines := strings.Split(stringBuffer, "\n")

		i := len(lines) - terminalHeight
		for i < len(lines)-1 {
			fmt.Println(lines[i])
			i++
		}
		fmt.Print(lines[i])
		printHeight = terminalHeight
		cursorReset = focusedLine - (len(lines) - terminalHeight) + 1
	}

	// move cursor up to focused Id
	fmt.Printf("\033[%dF", printHeight-cursorReset)
	return cursorReset
}

func (x *Menu) MoveCursorUp() (success bool) {
	if x.Focused_Index > 0 {
		x.Focused_Index--
		return true
	}
	return false
}

func (x *Menu) MoveCursorDown() (success bool) {
	if x.Focused_Index < len(x.Tasks)-1 {
		x.Focused_Index++
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
		menu.OutputChannel <- menu.Tasks
		menu.Editing = Tasks
	}
}

func (menu *Menu) EscKeyPressed() (kill bool) {
	switch menu.Editing {
	case Tasks:
		fmt.Printf("\033[%dB", 2*len(menu.Tasks))
		menu.OutputChannel <- nil
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
			desc := menu.Tasks[menu.Focused_Index].Description
			menu.Tasks[menu.Focused_Index].Description = desc[0:max(len(desc)-1, 0)]
			return
		} else if event.Key == keyboard.KeySpace {
			desc := menu.Tasks[menu.Focused_Index].Description
			menu.Tasks[menu.Focused_Index].Description = desc + " "
			return
		}
		menu.Tasks[menu.Focused_Index].Description += string(event.Rune)
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
