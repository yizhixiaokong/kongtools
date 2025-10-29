package pages

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/styles"
)

// WelcomeMsg 欢迎页 ASCII 艺术字（大版本）
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

// WelcomeMsgSmall 欢迎页 ASCII 艺术字（小版本）
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

// WelcomePage 欢迎页面
type WelcomePage struct {
	width       int
	height      int
	keys        welcomeKeyMap
	autoSkipped bool // 标记是否已经自动跳过
}

// welcomeKeyMap 欢迎页面快捷键映射
type welcomeKeyMap struct {
	Enter key.Binding
	Quit  key.Binding
}

// ShortHelp 返回简短帮助信息
func (k welcomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Quit}
}

// FullHelp 返回完整帮助信息
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

// NewWelcomePage 创建欢迎页面
func NewWelcomePage() *WelcomePage {
	return &WelcomePage{
		keys:        welcomeKeys,
		autoSkipped: false,
	}
}

// WelcomeTimeoutMsg 欢迎页超时消息
type WelcomeTimeoutMsg struct{}

// Init 实现 Page 接口
func (m *WelcomePage) Init() tea.Cmd {
	// 启动3秒定时器
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return WelcomeTimeoutMsg{}
	})
}

// Update 实现 Page 接口
func (m *WelcomePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			// Enter 或空格键，切换到主菜单
			return m, func() tea.Msg {
				return SwitchPageMsg{Page: "main"}
			}
		case key.Matches(msg, m.keys.Quit):
			// q 或 ctrl+c，退出
			return m, tea.Quit
		}
	case WelcomeTimeoutMsg:
		// 超时后发送切换页面消息
		if !m.autoSkipped {
			m.autoSkipped = true
			return m, func() tea.Msg {
				return SwitchPageMsg{Page: "main"}
			}
		}
	}

	return m, nil
}

// View 实现 Page 接口
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

	// 标题
	titleText := "Welcome to KongTools"
	if useSmallVersion && m.width < 50 {
		titleText = "KongTools"
	}
	title := styles.TitleStyle.Width(m.width).Render(titleText)

	// ASCII 艺术字 - 直接拼接不要逐行渲染，避免错位
	asciiArt := styles.WelcomeAsciiStyle.Render(strings.Join(asciiMsg, "\n"))

	// 描述文本
	descText := "🛠️  A collection of useful tools\n\nPress Enter to continue or wait 3 seconds..."
	if useSmallVersion && m.width < 60 {
		descText = "🛠️  Useful tools\n\nPress Enter or wait 3s..."
	}
	description := styles.WelcomeDescStyle.Width(m.width).Render(descText)

	// 组合所有内容
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		asciiArt,
		description,
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
func (m *WelcomePage) Title() string {
	return "欢迎"
}

// Help 实现 Page 接口
func (m *WelcomePage) Help() help.KeyMap {
	return m.keys
}

// SetSize 实现 Page 接口
func (m *WelcomePage) SetSize(width, height int) {
	m.width = width
	m.height = height
}
