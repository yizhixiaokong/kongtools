package pages

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/styles"
)

var WelcomeMsg = []string{
	``,
	``,
	`██╗  ██╗ ██████╗ ███╗   ██╗ ██████╗     ████████╗ ██████╗  ██████╗ ██╗     ███████╗`,
	`██║ ██╔╝██╔═══██╗████╗  ██║██╔════╝     ╚══██╔══╝██╔═══██╗██╔═══██╗██║     ██╔════╝`,
	`█████╔╝ ██║   ██║██╔██╗ ██║██║  ███╗       ██║   ██║   ██║██║   ██║██║     ███████╗`,
	`██╔═██╗ ██║   ██║██║╚██╗██║██║   ██║       ██║   ██║   ██║██║   ██║██║     ╚════██║`,
	`██║  ██╗╚██████╔╝██║ ╚████║╚██████╔╝       ██║   ╚██████╔╝╚██████╔╝███████╗███████║`,
	`╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝        ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝`,
}

var WelcomeMsgSmall = []string{
	``,
	` _                       _              _     `,
	`| |                     | |            | |    `,
	`| | _____  _ __   __ _  | |_ ___   ___ | |___ `,
	`| |/ / _ \| '_ \ / _' | | __/ _ \ / _ \| / __|`,
	`|   < (_) | | | | (_| | | || (_) | (_) | \__ \`,
	`|_|\_\___/|_| |_|\__, |  \__\___/ \___/|_|___/`,
	`                  __/ |                       `,
	`                 |___/                        `,
}

type WelcomePage struct {
	width       int
	height      int
	keys        welcomeKeyMap
	autoSkipped bool
}

type welcomeKeyMap struct {
	Enter key.Binding
	Quit  key.Binding
}

func (k welcomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Quit}
}

func (k welcomeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Quit},
	}
}

var welcomeKeys = welcomeKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "继续"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "退出"),
	),
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{
		keys:        welcomeKeys,
		autoSkipped: false,
	}
}

func (m *WelcomePage) Init() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return messages.WelcomeTimeoutMsg{}
	})
}

func (m *WelcomePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			return m, func() tea.Msg {
				return messages.SwitchPageMsg{Page: "main"}
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case messages.WelcomeTimeoutMsg:
		if !m.autoSkipped {
			m.autoSkipped = true
			return m, func() tea.Msg {
				return messages.SwitchPageMsg{Page: "main"}
			}
		}
	}

	return m, nil
}

func (m *WelcomePage) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var asciiMsg []string
	var useSmallVersion bool

	if m.width < 95 {
		asciiMsg = WelcomeMsgSmall
		useSmallVersion = true
	} else {
		asciiMsg = WelcomeMsg
		useSmallVersion = false
	}

	titleText := "Welcome to KongTools"
	if useSmallVersion && m.width < 50 {
		titleText = "KongTools"
	}
	title := styles.TitleStyle.Width(m.width).Render(titleText)

	asciiArt := styles.WelcomeAsciiStyle.Render(strings.Join(asciiMsg, "\n"))

	descText := "🛠️  A collection of useful tools\n\nPress Enter to continue or wait 3 seconds..."
	if useSmallVersion && m.width < 60 {
		descText = "🛠️  Useful tools\n\nPress Enter or wait 3s..."
	}
	description := styles.WelcomeDescStyle.Width(m.width).Render(descText)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		asciiArt,
		description,
	)

	// 使用通用居中函数
	return styles.CenterVertically(content, m.width, m.height)
}

func (m *WelcomePage) Title() string {
	return "欢迎"
}

func (m *WelcomePage) Help() help.KeyMap {
	return m.keys
}

func (m *WelcomePage) SetSize(width, height int) {
	m.width = width
	m.height = height
}
