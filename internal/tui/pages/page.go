package pages

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// Page 页面接口 - 所有页面都需要实现这个接口
type Page interface {
	// Init 初始化页面
	Init() tea.Cmd

	// Update 更新页面状态
	Update(tea.Msg) (tea.Model, tea.Cmd)

	// View 渲染页面视图
	View() string

	// Title 返回页面标题
	Title() string

	// Help 返回帮助信息的 KeyMap
	Help() help.KeyMap

	// SetSize 设置页面尺寸
	SetSize(width, height int)
}
