package view

import (
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Task struct {
	Title     string
	Completed bool
}

type TodoList struct {
	// ui
	*tview.Flex
	input *tview.InputField
	hint  *tview.TextView
	tasks *tview.List

	// data
	taskItems []Task

	// control
	editMode  bool
	editIndex int
	hintTimer *time.Timer

	// global
	logger *slog.Logger
}

func NewTodoList(logger *slog.Logger, tasks []Task) *TodoList {
	todoList := &TodoList{
		Flex:      tview.NewFlex(),
		input:     tview.NewInputField(),
		hint:      tview.NewTextView(),
		tasks:     tview.NewList(),
		taskItems: tasks,
		editMode:  false,
		editIndex: -1,
		hintTimer: nil,
		logger:    logger.With("module", "view-todo-list"),
	}

	todoList.initTasks()
	todoList.updateInputLabel()
	todoList.configureHandlers()
	todoList.setupLayout()

	return todoList
}

func (t *TodoList) initTasks() {
	t.tasks.Clear()
	t.logger.Debug("init tasks", slog.Int("count", len(t.taskItems)))
	t.tasks.ShowSecondaryText(false)

	if len(t.taskItems) == 0 {
		t.logger.Debug("No tasks found, adding help messages")
		t.addHelpMessages()
		return
	}
	t.updateTasksDisplay(0)
}

func (t *TodoList) displayTask(task Task) {
	title := task.Title

	if task.Completed {
		title = "[gray]" + tview.Escape("[x]") + title + "[-]"
	} else {
		title = "[white]" + tview.Escape("[ ]") + title + "[-]"
	}

	t.tasks.AddItem(title, "", 0, nil)
}

func (t *TodoList) displayTaskUpdate(index int, task Task) {
	if index >= len(t.taskItems) {
		return
	}

	title := task.Title

	if task.Completed {
		title = "[gray]" + tview.Escape("[x]") + title + "[-]"
	} else {
		title = "[white]" + tview.Escape("[ ]") + title + "[-]"
	}

	if index < t.tasks.GetItemCount() {
		t.tasks.SetItemText(index, title, "")
	} else {
		t.tasks.AddItem(title, "", 0, nil)
	}
}

func (t *TodoList) updateTasksDisplay(index int) {
	if index >= len(t.taskItems) {
		index = len(t.taskItems) - 1
	}
	// 删除index之后
	for i := t.tasks.GetItemCount() - 1; i > index; i-- {
		t.tasks.RemoveItem(i)
	}

	if index < 0 {
		return
	}
	t.displayTaskUpdate(index, t.taskItems[index])

	for i, task := range t.taskItems {
		if i <= index {
			continue
		}

		t.displayTask(task)
	}
}

func (t *TodoList) updateInputLabel() {
	label := "New To-Do: "
	if t.editMode {
		label = "Edit To-Do: "
	}
	t.input.SetLabel(label).
		SetLabelColor(tcell.ColorYellow).
		SetLabelWidth(len(label))
}

func (t *TodoList) updateHint(message string) {
	if t.hintTimer != nil {
		t.hintTimer.Stop()
	}

	t.hint.SetText(message).
		SetTextColor(tcell.ColorRed)

	if message != "" {
		t.hintTimer = time.AfterFunc(3*time.Second, func() {
			t.hint.SetText("")
		})
	}
}

var helpMessage = []string{
	"💡Write your first to-do task in the input field above.",
	"👏Press Enter to add the task to the list.",
	"📝Select a task and press Enter to edit it.",
	"🤷Press Esc to cancel editing a task.",
	"🥷Press Delete to remove a selected task.",
	"✅Press Space to mark a task as completed.",
}

func (t *TodoList) addHelpMessages() {
	for _, msg := range helpMessage {
		newTask := Task{
			Title:     msg,
			Completed: false,
		}
		t.taskItems = append(t.taskItems, newTask)
	}
	t.updateTasksDisplay(0)
}

func (t *TodoList) AddTask() {
	task := t.input.GetText()
	if task != "" {
		newTask := Task{
			Title:     task,
			Completed: false,
		}
		t.taskItems = append(t.taskItems, newTask)
		t.updateTasksDisplay(len(t.taskItems) - 1)
		// t.displayTask(newTask)
		t.input.SetText("")
		t.logger.Debug("Task added", slog.String("task", task))
	}
}

func (t *TodoList) DeleteTask() {
	if t.editMode {
		t.logger.Debug("Task edit mode, can't delete")
		t.updateHint("Cannot delete while editing a task.")
		return
	}
	if t.tasks.GetItemCount() == 0 {
		t.logger.Debug("Task list is empty, can't delete")
		t.updateHint("Task list is empty. Add a task first.")
		return
	}

	index := t.tasks.GetCurrentItem()
	if index >= len(t.taskItems) {
		return
	}

	task := t.taskItems[index].Title
	t.taskItems = append(t.taskItems[:index], t.taskItems[index+1:]...)
	t.updateTasksDisplay(index)
	t.logger.Debug("Task deleted", slog.String("task", task))

}

func (t *TodoList) EditTask() {
	if t.tasks.GetItemCount() > 0 {
		index := t.tasks.GetCurrentItem()
		if index >= len(t.taskItems) {
			return
		}

		task := t.taskItems[index].Title
		t.input.SetText(task)
		t.editMode = true
		t.editIndex = index
		t.updateInputLabel()
		t.logger.Debug("Task edit", slog.String("task", task))
	}
}

func (t *TodoList) SaveEdit() {
	if t.editMode {
		task := t.input.GetText()
		index := t.editIndex
		if task != "" && index >= 0 {
			t.taskItems[index].Title = task
			t.input.SetText("")
			t.editMode = false
			t.editIndex = -1
			t.updateInputLabel()
			t.updateTasksDisplay(index)
			t.logger.Debug("Task edited", slog.String("task", task))
		}
	}
}

func (t *TodoList) CancelEdit() {
	if t.editMode {
		t.input.SetText("")
		t.editMode = false
		t.editIndex = -1
		t.updateInputLabel()
		t.logger.Debug("Task edit canceled")
	}
}

func (t *TodoList) CompleteTask() {
	index := t.tasks.GetCurrentItem()
	if index >= len(t.taskItems) {
		return
	}

	t.taskItems[index].Completed = !t.taskItems[index].Completed
	t.updateTasksDisplay(index)
	t.logger.Debug("Task completion toggled", slog.String("task", t.taskItems[index].Title), slog.Bool("completed", t.taskItems[index].Completed))
}

func (t *TodoList) configureHandlers() {
	t.input.SetDoneFunc(t.handleInputDone)
	t.input.SetChangedFunc(t.handleInputText)
	t.tasks.SetInputCapture(t.handleListInput)
}

func (t *TodoList) handleInputText(text string) {
	if len(text) > 80 {
		t.input.SetText(text[:80])
		t.updateHint("Task length should not exceed 80 characters.")
	}
}

func (t *TodoList) handleInputDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		if t.editMode {
			t.SaveEdit()
		} else {
			t.AddTask()
		}
	} else if key == tcell.KeyEsc {
		t.CancelEdit()
	}
}

func (t *TodoList) handleListInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case ' ':
			t.CompleteTask()
			return nil
		default:
			return event
		}
	case tcell.KeyDelete:
		t.DeleteTask()
		return nil
	case tcell.KeyEnter:
		if t.editMode {
			t.SaveEdit()
		} else {
			t.EditTask()
		}
		return nil
	case tcell.KeyEsc:
		t.CancelEdit()
		return nil
	default:
		// t.logger.Debug("Unhandled key", slog.String("key", event.Name()), slog.Int("key", int(event.Key())))
		return event
	}
}

func (t *TodoList) setupLayout() {
	t.SetDirection(tview.FlexRow).
		AddItem(t.input, 1, 1, true).
		AddItem(t.hint, 1, 1, false).
		AddItem(t.tasks, 0, 1, false)

	t.SetBorder(true).
		SetTitle("To-Do List").
		SetTitleAlign(tview.AlignCenter)
}
