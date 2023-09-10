package structs

type FocusedTaskPathNode struct {
	Node *Task
	Next *FocusedTaskPathNode
	Prev *FocusedTaskPathNode
}
