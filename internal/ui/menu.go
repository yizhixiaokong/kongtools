package ui

import "github.com/rivo/tview"

// Menu 菜单
type Menu struct {
	*tview.List
}

// NewMenu 新建
func NewMenu() *Menu {
	m := &Menu{
		List: tview.NewList(),
	}
	m.SetBorder(true).SetTitle("Menu")

	return m
}

// SetTitle 设置标题
func (m *Menu) SetTitle(title string) *Menu {
	m.List.SetTitle(title)
	return m
}

// AddItem 添加菜单项
func (m *Menu) AddItem(text, secondaryText string, shortcut rune, selected func()) *Menu {
	m.List.AddItem(text, secondaryText, shortcut, selected)
	return m
}
