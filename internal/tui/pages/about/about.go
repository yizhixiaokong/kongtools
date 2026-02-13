package about

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/pkg/version"
	"kongtools/internal/tui/keys"
	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/styles"
)

type AboutPage struct {
	width  int
	height int
	keys   aboutKeyMap
}

type aboutKeyMap struct {
	keys.CommonKeyMap
}

func (k aboutKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k aboutKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var aboutKeys = aboutKeyMap{
	CommonKeyMap: keys.NewCommonKeyMap(),
}

func NewAboutPage() *AboutPage {
	return &AboutPage{
		keys: aboutKeys,
	}
}

func (m *AboutPage) Init() tea.Cmd {
	return nil
}

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

	return styles.CenterVertically(content, m.width, m.height)
}

func (m *AboutPage) Title() string {
	return "关于"
}

func (m *AboutPage) Help() help.KeyMap {
	return m.keys
}

func (m *AboutPage) SetSize(width, height int) {
	m.width = width
	m.height = height
}
