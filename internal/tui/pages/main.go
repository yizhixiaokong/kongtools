package pages

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/styles"
)

// MainPage 主菜单页面
type MainPage struct {
	width  int
	height int
	list   list.Model
	keys   mainKeyMap
}

// mainKeyMap 主菜单快捷键映射
type mainKeyMap struct {
	Enter key.Binding
	Quit  key.Binding
	Help  key.Binding
}

// ShortHelp 返回简短帮助信息
func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Quit}
}

// FullHelp 返回完整帮助信息
func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Help},
		{k.Quit},
	}
}

var mainKeys = mainKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "选择"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "退出"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "帮助"),
	),
}

// menuItem 菜单项
type menuItem struct {
	title       string
	description string
	page        string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// menuItemDelegate 菜单项委托
type menuItemDelegate struct{}

func (d menuItemDelegate) Height() int                             { return 2 }
func (d menuItemDelegate) Spacing() int                            { return 1 }
func (d menuItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d menuItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	var title, desc string
	if index == m.Index() {
		// 选中状态 - 使用更有科技感的箭头
		title = styles.MenuSelectedStyle.Render(fmt.Sprintf("› %s", i.title))
		desc = styles.MenuDescStyle.
			Foreground(styles.ColorGreen).
			Render(i.description)
	} else {
		// 普通状态
		title = styles.MenuItemStyle.Render(fmt.Sprintf("  %s", i.title))
		desc = styles.MenuDescStyle.Render(i.description)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

// NewMainPage 创建主菜单页面
func NewMainPage() *MainPage {
	items := []list.Item{
		menuItem{
			title:       "📝 Todo List",
			description: "管理你的待办事项",
			page:        "todo",
		},
		menuItem{
			title:       "⚙️  Settings",
			description: "配置应用设置",
			page:        "settings",
		},
		menuItem{
			title:       "ℹ️  About",
			description: "关于 KongTools",
			page:        "about",
		},
	}

	delegate := menuItemDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = "🏠 Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = styles.MenuTitleStyle

	return &MainPage{
		list: l,
		keys: mainKeys,
	}
}

// Init 实现 Page 接口
func (m *MainPage) Init() tea.Cmd {
	return nil
}

// Update 实现 Page 接口
func (m *MainPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				// 发送页面切换消息
				return m, func() tea.Msg {
					return SwitchPageMsg{Page: i.page}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 实现 Page 接口
func (m *MainPage) View() string {
	listView := m.list.View()

	// 使用 lipgloss.Place 实现水平和垂直居中
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		listView,
	)
}

// Title 实现 Page 接口
func (m *MainPage) Title() string {
	return "主菜单"
}

// Help 实现 Page 接口
func (m *MainPage) Help() help.KeyMap {
	return m.keys
}

// SetSize 实现 Page 接口
func (m *MainPage) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetWidth(width)
	m.list.SetHeight(height - 4)
}

// SwitchPageMsg 页面切换消息
type SwitchPageMsg struct {
	Page string
}
