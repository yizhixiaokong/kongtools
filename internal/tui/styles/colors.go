package styles

import "github.com/charmbracelet/lipgloss"

// 主题颜色 - 使用 AdaptiveColor 支持明暗主题
var (
	// 主要颜色
	Primary   = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#7D56F4"}
	Secondary = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#874BFD"}
	Accent    = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"}

	// 状态颜色
	Success = lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}
	Warning = lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#F59E0B"}
	Error   = lipgloss.AdaptiveColor{Light: "#EF4444", Dark: "#EF4444"}
	Info    = lipgloss.AdaptiveColor{Light: "#3B82F6", Dark: "#3B82F6"}

	// 文本颜色
	Text       = lipgloss.AdaptiveColor{Light: "#212121", Dark: "#FAFAFA"}
	TextSubtle = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	TextMuted  = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"}

	// 边框颜色
	Border       = lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#404040"}
	BorderActive = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#874BFD"}

	// 背景颜色
	Background         = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#1A1A1A"}
	BackgroundSelected = lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#874BFD"}
	BackgroundHover    = lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#2A2A2A"}
)

// 特殊颜色值
var (
	ColorGreen  = lipgloss.Color("86")
	ColorRed    = lipgloss.Color("203")
	ColorYellow = lipgloss.Color("220")
	ColorBlue   = lipgloss.Color("39")
	ColorGray   = lipgloss.Color("240")
)
