package keys

import "github.com/charmbracelet/bubbles/key"

// CommonKeyMap 提供通用的键绑定（Back 和 Quit）
// 可以嵌入到其他页面的 KeyMap 结构中
type CommonKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

// NewCommonKeyMap 创建通用键绑定
func NewCommonKeyMap() CommonKeyMap {
	return CommonKeyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "返回"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "退出"),
		),
	}
}

// QuitKeyMap 只包含退出键的简单 KeyMap
// 用于欢迎页等只需要退出的页面
type QuitKeyMap struct {
	Quit key.Binding
}

// NewQuitKeyMap 创建只包含退出键的 KeyMap
func NewQuitKeyMap() QuitKeyMap {
	return QuitKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "退出"),
		),
	}
}
