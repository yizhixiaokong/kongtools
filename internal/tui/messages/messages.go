package messages

import "time"

// SwitchPageMsg 页面切换消息
type SwitchPageMsg struct {
	Page string
}

// NotificationLevel 通知级别
type NotificationLevel int

const (
	NotificationInfo NotificationLevel = iota
	NotificationSuccess
	NotificationWarning
	NotificationError
)

// NotificationMsg 通知消息
type NotificationMsg struct {
	Message string
	Level   NotificationLevel
}

// ClearNotificationMsg 清除通知消息
type ClearNotificationMsg struct{}

// SaveSuccessMsg 保存成功消息
type SaveSuccessMsg struct {
	Path string
}

// SaveFailedMsg 保存失败消息
type SaveFailedMsg struct {
	Err error
}

// WelcomeTimeoutMsg 欢迎页超时消息
type WelcomeTimeoutMsg struct{}

// ClearHintMsg 清除提示消息
type ClearHintMsg struct{}

// TickMsg 定时器消息
type TickMsg struct {
	Time time.Time
}
