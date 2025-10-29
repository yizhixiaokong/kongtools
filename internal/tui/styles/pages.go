package styles

import "github.com/charmbracelet/lipgloss"

// TabBorderWithBottom 创建带底部边框的标签页边框
func TabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	// 标签页边框
	InactiveTabBorder = TabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = TabBorderWithBottom("┘", " ", "└")

	// 标签页样式
	TabStyle = lipgloss.NewStyle().
			Border(InactiveTabBorder, true).
			BorderForeground(Primary).
			Padding(0, 1)

	// 直接复用并在需要时通过方法生成新样式（方法本身返回新的 style）
	InactiveTabStyle = TabStyle

	ActiveTabStyle = TabStyle.
			Border(ActiveTabBorder, true).
			Bold(true)

		// 标签页间隔
	TabGapStyle = TabStyle.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)

	// 标签页内容窗口
	TabWindowStyle = lipgloss.NewStyle().
			BorderForeground(Primary).
			Padding(1, 2).
			Align(lipgloss.Left)

	// 滚动指示器
	ScrollIndicatorStyle = lipgloss.NewStyle().
				Align(lipgloss.Right).
				Foreground(TextSubtle)
)

// Todo List 页面样式
var (
	// Todo 输入框标签
	TodoInputLabelStyle = lipgloss.NewStyle().
				Foreground(ColorYellow).
				Bold(true)

	// Todo 输入框
	TodoInputStyle = lipgloss.NewStyle().
			Foreground(Text)

	// Todo 输入光标
	TodoCursorStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(ColorGray)

	// Todo 列表项 - 普通
	TodoItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(Text)

	// Todo 列表项 - 选中
	TodoSelectedStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Background(ColorGreen).
				Foreground(lipgloss.Color("0")).
				Bold(true)

	// Todo 列表项 - 已完成
	TodoCompletedStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(ColorGray).
				Strikethrough(true)

	// Todo 列表项 - 选中且已完成
	TodoSelectedCompletedStyle = lipgloss.NewStyle().
					Padding(0, 2).
					Background(ColorGray).
					Foreground(lipgloss.Color("0")).
					Strikethrough(true)

	// Todo 提示信息
	TodoHintStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Padding(0, 0, 1, 0)
)

// Welcome 页面样式
var (
	WelcomeAsciiStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Align(lipgloss.Center)

	WelcomeDescStyle = lipgloss.NewStyle().
				Foreground(TextSubtle).
				Align(lipgloss.Center).
				Padding(2, 0)
)

// Main Menu 页面样式
var (
	MenuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGreen).
			Align(lipgloss.Center).
			Padding(0, 1)

	MenuItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(Text)

	MenuSelectedStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Background(ColorGreen).
				Foreground(lipgloss.Color("0")).
				Bold(true)

	MenuDescStyle = lipgloss.NewStyle().
			Foreground(TextSubtle).
			Padding(0, 4)
)
