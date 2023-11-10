package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/sagikazarmark/slog-shim"
)

// Menu 菜单
type Menu struct {
	*tview.List

	logger *slog.Logger
}

// NewMenu 新建
func NewMenu(logger *slog.Logger) *Menu {
	m := &Menu{
		List:   tview.NewList(),
		logger: logger.With("module", "ui-menu"),
	}
	m.SetBorder(true).SetTitle("Menu")

	return m
}

// SetTitle 设置标题
func (m *Menu) SetTitle(title string) *Menu {
	m.logger.Debug(fmt.Sprintf("set menu title, from: %s, to: %s.", m.GetTitle(), title))
	m.List.SetTitle(title)
	return m
}

// AddItem 添加菜单项
func (m *Menu) AddItem(text, secondaryText string, shortcut rune, selected func()) *Menu {
	m.logger.Debug(fmt.Sprintf("add menu item, text: %s, secondaryText: %s, shortcut: %s.", text, secondaryText, string(shortcut)))
	m.List.AddItem(text, secondaryText, shortcut, selected)
	return m
}
