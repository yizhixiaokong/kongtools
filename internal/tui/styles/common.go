package styles

import "github.com/charmbracelet/lipgloss"

// 通用样式
var (
	// Header 样式
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center)

	// Footer 样式
	FooterStyle = lipgloss.NewStyle().
			Foreground(TextSubtle).
			Align(lipgloss.Center)

	// 分隔线样式
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(Border).
			Align(lipgloss.Center)

	// 标题样式
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	// 副标题样式
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(TextSubtle).
			Italic(true).
			Align(lipgloss.Center)

	// 内容样式
	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Align(lipgloss.Left)

	// 错误样式
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	// 成功样式
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	// 警告样式
	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	// 信息样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)
)

// 边框样式
var (
	// 普通边框
	NormalBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Border)

	// 激活边框
	ActiveBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(BorderActive)

	// 圆角边框
	RoundedBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Border)

	// 激活圆角边框
	ActiveRoundedBorderStyle = lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(BorderActive)
)

// 列表项样式
var (
	// 普通列表项
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(Text)

	// 选中列表项
	SelectedListItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Background(BackgroundSelected).
				Foreground(Background).
				Bold(true)

	// 悬停列表项
	HoverListItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Background(BackgroundHover).
				Foreground(Text)
)

// 帮助信息样式
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(TextSubtle).
			Padding(0, 1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(TextSubtle)
)
