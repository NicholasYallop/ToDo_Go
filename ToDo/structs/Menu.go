package structs

import (
	style "ToDo/defs"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	Focused_Task  *Task
	Editing       EditingField
	OutputChannel chan []Task
}

type TaskStrings struct {
	taskStrings []string
	descStrings []string
}

func (menu *Menu) TaskStrings(index int) (taskStrings TaskStrings, focused bool) {
	task := &menu.Tasks[index]

	var taskString string
	var descString string
	if menu.Focused_Task == task {
		if menu.Editing == Tasks {
			taskString = "> " + task.Name
			descString = task.Description
		}
		if menu.Editing == Descriptions {
			taskString = task.Name
			descString = "> " + task.Description
		}
		focused = true

	} else {
		taskString = task.Name
		descString = task.Description
		focused = false
	}

	if task.Complete {
		taskString = lipgloss.JoinHorizontal(lipgloss.Center, style.Defocused_Default_Style.Render(taskString), "{Θ}")
		descString = style.Defocused_Remark_Style.Render(descString)
	} else {
		taskString = lipgloss.JoinHorizontal(lipgloss.Center, style.Default_Style.Render(taskString), "{ }")
		descString = style.Remark_Style.Render(descString)
	}

	return TaskStrings{
		taskStrings: strings.Split(taskString, "\n"),
		descStrings: strings.Split(descString, "\n"),
	}, focused
}

func (menu Menu) Print() (cursorRow int) {
	_, terminalHeight, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error while printing.")
		panic(err)
	}

	// print from startindex, height of at most terminal height
	focusedLine := 0
	var lines []string

	for i := 0; i < len(menu.Tasks); i++ {
		printData, focused := menu.TaskStrings(i)

		if focused && menu.Editing == Tasks {
			focusedLine = len(lines)
		}

		lines = append(lines, printData.taskStrings...)

		if focused && menu.Editing == Descriptions {
			focusedLine = len(lines)
		}

		lines = append(lines, printData.descStrings...)
	}

	printHeight := 0
	cursorReset := 0
	if terminalHeight >= len(lines) {
		// print all lines
		for index, line := range lines {
			if index == len(lines)-1 {
				fmt.Print(line)
			} else {
				fmt.Println(line)
			}
		}
		printHeight = len(lines) - 1
		cursorReset = focusedLine

	} else if focusedLine <= len(lines)-terminalHeight {
		// print from focused up to terminal height
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

func (menu *Menu) MoveCursorUp() (success bool) {
	if menu.Focused_Task != &menu.Tasks[0] {
		for index := range menu.Tasks {
			if &menu.Tasks[index] == menu.Focused_Task {
				menu.Focused_Task = &menu.Tasks[index-1]
			}
		}
		return true
	}

	return false
}

func (menu *Menu) MoveCursorDown() (success bool) {
	if menu.Focused_Task != &menu.Tasks[len(menu.Tasks)-1] {
		for index := range menu.Tasks {
			if &menu.Tasks[index] == menu.Focused_Task {
				menu.Focused_Task = &menu.Tasks[index+1]
				break
			}
		}
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
			desc := menu.Focused_Task.Description
			menu.Focused_Task.Description = desc[0:max(len(desc)-1, 0)]
			return
		} else if event.Key == keyboard.KeySpace {
			desc := menu.Focused_Task.Description
			menu.Focused_Task.Description = desc + " "
			return
		}
		menu.Focused_Task.Description += string(event.Rune)
	case Tasks:
		if event.Rune == 'c' {
			menu.Focused_Task.Complete = !menu.Focused_Task.Complete
		}
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
