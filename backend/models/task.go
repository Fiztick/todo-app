package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskNode struct {
	Task Task
	Next *TaskNode
}

// Now has both Head and Tail
type TaskList struct {
	Head *TaskNode
	Tail *TaskNode
	Size int
}

// Enqueue adds a new task to the BACK of the list
func (tl *TaskList) Enqueue(task Task) {
	newNode := &TaskNode{Task: task}
	if tl.Tail == nil {
		// list is empty, head and tail are the same node
		tl.Head = newNode
		tl.Tail = newNode
	} else {
		// attach to the back, move tail forward
		tl.Tail.Next = newNode
		tl.Tail = newNode
	}
	tl.Size++
}

// Dequeue removes and returns the task at the FRONT of the list
func (tl *TaskList) Dequeue() (Task, bool) {
	if tl.Head == nil {
		return Task{}, false
	}
	task := tl.Head.Task
	tl.Head = tl.Head.Next
	if tl.Head == nil {
		// list is now empty, reset tail too
		tl.Tail = nil
	}
	tl.Size--
	return task, true
}

// GetAll walks the list front to back
func (tl *TaskList) GetAll() []Task {
	var tasks []Task
	current := tl.Head
	for current != nil {
		tasks = append(tasks, current.Task)
		current = current.Next
	}
	return tasks
}

// Remove finds and removes a task by ID
func (tl *TaskList) Remove(id int) bool {
	if tl.Head == nil {
		return false
	}

	if tl.Head.Task.ID == id {
		tl.Head = tl.Head.Next
		if tl.Head == nil {
			tl.Tail = nil
		}
		tl.Size--
		return true
	}

	current := tl.Head
	for current.Next != nil {
		if current.Next.Task.ID == id {
			if current.Next == tl.Tail {
				// removing the tail, update tail pointer
				tl.Tail = current
			}
			current.Next = current.Next.Next
			tl.Size--
			return true
		}
		current = current.Next
	}
	return false
}