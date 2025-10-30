package about

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/pkg/version"
	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/styles"
)

// AboutPage 关于页面
type AboutPage struct {
	width  int
	height int
	keys   aboutKeyMap
}

// aboutKeyMap 关于页面快捷键
type aboutKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

// ShortHelp 返回简短帮助信息
func (k aboutKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp 返回完整帮助信息
func (k aboutKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var aboutKeys = aboutKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "返回"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "退出"),
	),
}

// NewAboutPage 创建关于页面
func NewAboutPage() *AboutPage {
	return &AboutPage{
		keys: aboutKeys,
	}
}

// Init 实现 Page 接口
func (m *AboutPage) Init() tea.Cmd {
	return nil
}

// Update 实现 Page 接口
func (m *AboutPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return messages.SwitchPageMsg{Page: "main"}
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

// View 实现 Page 接口
func (m *AboutPage) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	titleStyle := styles.TitleStyle.
		Width(m.width).
		Padding(2, 0)

	descStyle := styles.SubtitleStyle.
		Width(m.width).
		Padding(1, 0)

	backHint := styles.InfoStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Padding(2, 0).
		Render("按 ESC 返回主菜单")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("ℹ️  关于"),
		descStyle.Render(fmt.Sprintf("KongTools %s\n\n一个实用工具集合\n\nCommit: %s\nBuild Time: %s\nGo Version: %s",
			version.GetVersion(),
			version.GetCommit(),
			version.GetBuildTime(),
			version.GetGoVersion())),
		backHint,
	)

	// 垂直居中
	contentHeight := lipgloss.Height(content)
	topPadding := (m.height - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	paddingStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return paddingStyle.Render(content)
}

// Title 实现 Page 接口
func (m *AboutPage) Title() string {
	return "关于"
}

// Help 实现 Page 接口
func (m *AboutPage) Help() help.KeyMap {
	return m.keys
}

// SetSize 实现 Page 接口
func (m *AboutPage) SetSize(width, height int) {
	m.width = width
	m.height = height
}
