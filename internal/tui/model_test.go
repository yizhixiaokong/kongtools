package tui

import (
	"log/slog"
	"os"
	"testing"

	"kongtools/internal/tui/messages"
)

func TestNewModel(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := Config{TasksSavePath: "/tmp/test-tasks.json"}

	model := NewModel(logger, cfg)

	if model.currentPage != PageWelcome {
		t.Errorf("Expected PageWelcome, got %v", model.currentPage)
	}

	if len(model.pages) != 6 {
		t.Errorf("Expected 6 pages, got %d", len(model.pages))
	}
}

func TestSwitchPage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := Config{TasksSavePath: "/tmp/test-tasks.json"}
	model := NewModel(logger, cfg)

	msg := messages.SwitchPageMsg{Page: "main"}
	newModel, _ := model.Update(msg)

	m := newModel.(Model)
	if m.currentPage != PageMain {
		t.Errorf("Expected PageMain, got %v", m.currentPage)
	}
}

func TestContentSize(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := Config{TasksSavePath: "/tmp/test-tasks.json"}
	model := NewModel(logger, cfg)

	model.height = 30
	expectedContentHeight := 30 - TotalNonContent

	if h := model.getContentHeight(); h != expectedContentHeight {
		t.Errorf("Expected content height %d, got %d", expectedContentHeight, h)
	}
}

func TestNotification(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := Config{TasksSavePath: "/tmp/test-tasks.json"}
	model := NewModel(logger, cfg)

	msg := messages.SaveSuccessMsg{Path: "/tmp/test.json"}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.notification == "" {
		t.Error("Notification should be set")
	}
}

func TestPageTypeString(t *testing.T) {
	tests := []struct {
		page     PageType
		expected string
	}{
		{PageWelcome, "欢迎"},
		{PageMain, "主菜单"},
		{PageTodo, "待办事项"},
		{PageImage, "图片预览"},
		{PageSettings, "设置"},
		{PageAbout, "关于"},
	}

	for _, tt := range tests {
		result := tt.page.String()
		if result != tt.expected {
			t.Errorf("PageType(%d).String() = %s, expected %s", tt.page, result, tt.expected)
		}
	}
}
