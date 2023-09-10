package structs

type FocusedTask struct {
	Task *Task
	Path FocusedTaskPath
}

func (focusedTask *FocusedTask) FocusDown(target *Task) {
	focusedTask.Path.Last.Next = &FocusedTaskPathNode{
		Node: target,
		Prev: focusedTask.Path.Last,
	}

	if focusedTask.Path.First == focusedTask.Path.Last {
		focusedTask.Path.First.Next = focusedTask.Path.Last.Next
	}

	focusedTask.Path.Last = focusedTask.Path.Last.Next

	focusedTask.Task = focusedTask.Path.Last.Node
}

func (focusedTask *FocusedTask) FocusUp() {
	focusedTask.Path.Last = focusedTask.Path.Last.Prev
	focusedTask.Path.Last.Next = nil
	focusedTask.Task = focusedTask.Path.Last.Node
}

func (focusedTask *FocusedTask) Refocus(newFocus *Task) {
	focusedTask.Task = newFocus

	if focusedTask.Path.First == focusedTask.Path.Last {
		focusedTask.Path.First.Node = focusedTask.Task
	}
	focusedTask.Path.Last.Node = focusedTask.Task
}
