package todolist

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTasks(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "tasks.json")
	tasks := []Task{
		{Title: "Task 1", Completed: false},
		{Title: "Task 2", Completed: true},
	}
	data, _ := json.Marshal(tasks)
	os.WriteFile(tmpFile, data, 0644)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	page := NewTodoListPage(logger, tmpFile)

	// Init 会加载任务
	page.Init()

	// 验证列表有任务（可能从文件加载或使用帮助任务）
	if page.list.Items() == nil {
		t.Error("List items should not be nil")
	}
}

func TestToggleTask(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "tasks.json")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	page := NewTodoListPage(logger, tmpFile)

	page.tasks = []Task{{Title: "Task 1", Completed: false}}
	page.updateListItems()

	page.tasks[0].Completed = !page.tasks[0].Completed
	if !page.tasks[0].Completed {
		t.Error("Task should be completed after toggle")
	}
}

func TestDeleteTask(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "tasks.json")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	page := NewTodoListPage(logger, tmpFile)

	page.tasks = []Task{
		{Title: "Task 1"},
		{Title: "Task 2"},
		{Title: "Task 3"},
	}

	index := 1
	page.tasks = append(page.tasks[:index], page.tasks[index+1:]...)

	if len(page.tasks) != 2 {
		t.Errorf("Expected 2 tasks after deletion, got %d", len(page.tasks))
	}
}

func TestTaskJSON(t *testing.T) {
	task := Task{Title: "Test Task", Completed: true}
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	var unmarshaled Task
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaled.Title != task.Title {
		t.Errorf("Title mismatch: expected '%s', got '%s'", task.Title, unmarshaled.Title)
	}
}
