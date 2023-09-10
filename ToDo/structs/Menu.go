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
	Tasks         TaskSlice
	Focused_Task  FocusedTask
	Editing       EditingField
	OutputChannel chan []Task
}

func NewMenu(tasks TaskSlice) (menu *Menu) {
	taskPointer := &tasks[0]
	taskNode := FocusedTaskPathNode{
		Node: &tasks[0],
		Next: nil,
	}

	return &Menu{
		Tasks: tasks,
		Focused_Task: FocusedTask{
			Task: taskPointer,
			Path: FocusedTaskPath{
				First: &taskNode,
				Last:  &taskNode,
			},
		},
		Editing:       Tasks,
		OutputChannel: make(chan []Task),
	}
}

func (menu *Menu) TaskStrings(task *Task, indentation int, collapsed bool) (taskStrings TaskStrings, focusedIndentation *int) {
	var taskString string
	var descString string
	var focused bool
	if menu.Focused_Task.Task == task {
		if menu.Editing == Tasks {
			taskString = " " + task.Name
			descString = task.Description
		}
		if menu.Editing == Descriptions {
			taskString = " " + task.Name
			descString = "> " + task.Description
		}
		focused = true

	} else {
		taskString = " " + task.Name
		descString = task.Description
		focused = false
	}

	var statusString string
	if task.Complete {
		statusString = "{C}"
	} else {
		statusString = "{ }"
	}

	taskStyle, taskStatuStyle, descStyle, descStatusStyle := style.GetTaskStyles(indentation, task.Complete)

	taskString = taskStyle.Render(taskString)
	descString = descStyle.Render(descString)

	taskStatuStyle.Height(lipgloss.Height(taskString))
	descStatusStyle.Height(lipgloss.Height(descString))

	taskString = lipgloss.JoinHorizontal(lipgloss.Center, taskString, taskStatuStyle.Render(statusString))
	descString = lipgloss.JoinHorizontal(lipgloss.Center, descString, descStatusStyle.Render(""))

	descLines := make([]string, 0)
	if !collapsed {
		descLines = strings.Split(descString, "\n")
	}

	focusedIndentation = nil
	if focused {
		focusedIndentation = &indentation
	}

	return TaskStrings{
		taskLines: strings.Split(taskString, "\n"),
		descLines: descLines,
		focused:   focused,
	}, focusedIndentation
}

func (menu *Menu) AppendToLines(lines *[]string, printData TaskStrings) (focusedLine int) {
	if printData.focused {
		focusedLine = len(*lines)
	}

	*lines = append(*lines, printData.taskLines...)
	*lines = append(*lines, printData.descLines...)

	return focusedLine
}

func (menu *Menu) FormattedPrintLines(tasks *TaskSlice, lines *[]string, focusedRow *int, focusedIndentation *int, indentation int) {
	var collapsed bool
	parentOfFocusedTask := menu.Focused_Task.Path.Last.Prev
	if parentOfFocusedTask != nil {
		collapsed = !(tasks == &parentOfFocusedTask.Node.SubTasks)
	} else {
		collapsed = !(tasks == &menu.Tasks)
	}

	for i := 0; i < len(*tasks); i++ {
		printData, indent := menu.TaskStrings(&(*tasks)[i], indentation, collapsed)

		if indent != nil {
			*focusedIndentation = *indent
		}

		*focusedRow = max(menu.AppendToLines(lines, printData), *focusedRow)

		menu.FormattedPrintLines(&(*tasks)[i].SubTasks, lines, focusedRow, focusedIndentation, indentation+1)
	}
}

func (menu Menu) Print() (printHeight int, cursorReset int, cursorIndent int) {
	_, terminalHeight, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error while printing.")
		panic(err)
	}

	// print from startindex, height of at most terminal height
	var lines []string
	var focusedLine int
	var focusedIndentation int
	menu.FormattedPrintLines(&menu.Tasks, &lines, &focusedLine, &focusedIndentation, 0)

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

	return printHeight, cursorReset, focusedIndentation
}

func (menu *Menu) Display() {
	keys, err := keyboard.GetKeys(1)
	if err != nil {
		fmt.Println("Error getting key")
		panic(err)
	}
	defer keyboard.Close()

	for {
		printHeight, cursor, cursorIndent := menu.Print()
		menu.MoveCursorToFocus(printHeight-cursor, cursorIndent)

		keyEvent := <-keys
		if keyEvent.Err != nil {
			fmt.Println("Error with key event")
			panic(err)
		}

		if menu.HandleKeyInput(keyEvent) {
			menu.MoveCursorToEnd(printHeight - cursor)
			break
		}

		menu.ResetDisplayCursor(cursor)
	}
}

func (menu *Menu) ScrollFocusUp() (success bool) {
	var containerSlice *TaskSlice
	if menu.Focused_Task.Path.First == menu.Focused_Task.Path.Last {
		containerSlice = &menu.Tasks
	} else {
		containerSlice = &menu.Focused_Task.Path.Last.Prev.Node.SubTasks
	}

	if menu.Focused_Task.Task != &(*containerSlice)[0] {
		for index := range *containerSlice {
			if &((*containerSlice)[index]) == menu.Focused_Task.Task {
				menu.Focused_Task.Refocus(&((*containerSlice)[index-1]))
				break
			}
		}
		return true
	}
	return false
}

func (menu *Menu) ScrollFocusDown() (success bool) {
	var containerSlice *TaskSlice
	if menu.Focused_Task.Path.First == menu.Focused_Task.Path.Last {
		containerSlice = &menu.Tasks
	} else {
		containerSlice = &menu.Focused_Task.Path.Last.Prev.Node.SubTasks
	}

	if menu.Focused_Task.Task != &(*containerSlice)[len(*containerSlice)-1] {
		for index := range *containerSlice {
			if &((*containerSlice)[index]) == menu.Focused_Task.Task {
				menu.Focused_Task.Refocus(&((*containerSlice)[index+1]))
				break
			}
		}
		return true
	}
	return false
}

func (menu *Menu) MoveCursorToFocus(cursorOffset int, cursorIndent int) {
	if cursorOffset != 0 {
		fmt.Printf("\033[%dF", cursorOffset)
	}
	if cursorIndent != 0 {
		fmt.Printf("\033[%dC", cursorIndent*2)
	}
}

func (menu *Menu) MoveCursorToEnd(cursorOffset int) {
	fmt.Printf("\033[%dE", cursorOffset)
}

func (menu *Menu) ResetDisplayCursor(cursor int) (success bool) {
	if cursor != 0 {
		fmt.Printf("\033[%dF", cursor)
		return true
	}
	return false
}

func (menu *Menu) UpKeyPressed() {
	switch menu.Editing {
	case Tasks:
		menu.ScrollFocusUp()
	}
}

func (menu *Menu) DownKeyPressed() {
	switch menu.Editing {
	case Tasks:
		menu.ScrollFocusDown()
	}
}

func (menu *Menu) RightKeyPressed() {
	switch menu.Editing {
	case Tasks:
		if len(menu.Focused_Task.Task.SubTasks) != 0 {
			menu.Focused_Task.FocusDown(&menu.Focused_Task.Task.SubTasks[0])
		}
	}
}

func (menu *Menu) LeftKeyPressed() {
	switch menu.Editing {
	case Tasks:
		if menu.Focused_Task.Path.First != menu.Focused_Task.Path.Last {
			menu.Focused_Task.FocusUp()
		}
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
			desc := menu.Focused_Task.Task.Description
			menu.Focused_Task.Task.Description = desc[0:max(len(desc)-1, 0)]
			return
		} else if event.Key == keyboard.KeySpace {
			desc := menu.Focused_Task.Task.Description
			menu.Focused_Task.Task.Description = desc + " "
			return
		}
		menu.Focused_Task.Task.Description += string(event.Rune)
	case Tasks:
		if event.Rune == 'c' {
			menu.Focused_Task.Task.Complete = !menu.Focused_Task.Task.Complete
		}
		if event.Rune == 's' {
			menu.Focused_Task.Task.SubTasks = append([]Task{{}}, menu.Focused_Task.Task.SubTasks...)
			menu.OutputChannel <- menu.Tasks
		}
	}
}

func (menu *Menu) HandleKeyInput(event keyboard.KeyEvent) (kill bool) {
	switch event.Key {
	case keyboard.KeyArrowUp:
		menu.UpKeyPressed()
	case keyboard.KeyArrowDown:
		menu.DownKeyPressed()
	case keyboard.KeyArrowRight:
		menu.RightKeyPressed()
	case keyboard.KeyArrowLeft:
		menu.LeftKeyPressed()
	case keyboard.KeyEnter:
		menu.EnterKeyPressed()
	case keyboard.KeyEsc:
		return menu.EscKeyPressed()
	default:
		menu.DefaultKeyPress(event)
	}
	return false
}
