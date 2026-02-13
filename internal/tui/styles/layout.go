package styles

import "github.com/charmbracelet/lipgloss"

// CenterVertically 将内容垂直居中
// width: 可用宽度
// height: 可用高度
// content: 要居中的内容字符串
func CenterVertically(content string, width, height int) string {
	contentHeight := lipgloss.Height(content)
	topPadding := (height - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	paddingStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return paddingStyle.Render(content)
}

// CenterHorizontally 将内容水平居中
// width: 可用宽度
// content: 要居中的内容字符串
func CenterHorizontally(content string, width int) string {
	paddingStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center)

	return paddingStyle.Render(content)
}

// Center 将内容同时水平和垂直居中
// 这是 CenterVertically 的别名，因为它已经实现了同时居中
func Center(content string, width, height int) string {
	return CenterVertically(content, width, height)
}
