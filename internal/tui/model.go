package tui

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/pages"
	"kongtools/internal/tui/pages/about"
	"kongtools/internal/tui/pages/image"
	"kongtools/internal/tui/pages/settings"
	"kongtools/internal/tui/pages/todolist"
	"kongtools/internal/tui/styles"
)

// PageType 页面类型
type PageType int

const (
	PageWelcome PageType = iota
	PageMain
	PageTodo
	PageImage
	PageSettings
	PageAbout
)

func (p PageType) String() string {
	switch p {
	case PageWelcome:
		return "欢迎"
	case PageMain:
		return "主菜单"
	case PageTodo:
		return "待办事项"
	case PageImage:
		return "图片预览"
	case PageSettings:
		return "设置"
	case PageAbout:
		return "关于"
	default:
		return "未知"
	}
}

// Model 主模型
type Model struct {
	// 配置
	cfg Config

	// 当前页面
	currentPage PageType

	// 页面实例
	pages map[PageType]pages.Page

	// 通知
	notification      string
	notificationLevel messages.NotificationLevel

	// 窗口大小
	width  int
	height int

	// 帮助
	help help.Model

	// 日志
	logger *slog.Logger
}

// NewModel 创建新模型
func NewModel(logger *slog.Logger, cfg Config) *Model {
	m := &Model{
		cfg:         cfg,
		currentPage: PageWelcome,
		logger:      logger.With("module", "tui-model"),
		help:        help.New(),
		pages:       make(map[PageType]pages.Page),
	}

	// 初始化页面
	m.pages[PageWelcome] = pages.NewWelcomePage()
	m.pages[PageMain] = pages.NewMainPage()
	m.pages[PageTodo] = todolist.NewTodoListPage(logger, cfg.TasksSavePath)
	m.pages[PageImage] = image.NewImagePage()
	m.pages[PageSettings] = settings.NewSettingsPage()
	m.pages[PageAbout] = about.NewAboutPage()

	return m
}

// Init 实现 tea.Model 接口
func (m Model) Init() tea.Cmd {
	// 启动欢迎页的定时器
	// 同时初始化 Todo 页面（加载数据）
	var cmds []tea.Cmd
	if p, ok := m.pages[PageWelcome]; ok {
		cmds = append(cmds, p.Init())
	}
	if p, ok := m.pages[PageTodo]; ok {
		cmds = append(cmds, p.Init())
	}
	return tea.Batch(cmds...)
}

// Update 实现 tea.Model 接口
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		contentHeight := m.getContentHeight()
		contentWidth := m.width

		for _, p := range m.pages {
			p.SetSize(contentWidth, contentHeight)
		}

		m.logger.Debug("window size changed",
			slog.Int("width", m.width),
			slog.Int("height", m.height))

		return m, nil

	case tea.KeyMsg:
		// 全局快捷键
		switch msg.String() {
		case "ctrl+c":
			m.logger.Debug("quit app by ctrl+c")
			return m, tea.Quit
		}

	case messages.SwitchPageMsg:
		switch msg.Page {
		case "main":
			m.currentPage = PageMain
			m.logger.Debug("switch to main page")
		case "todo":
			m.currentPage = PageTodo
			m.logger.Debug("switch to todo page")
		case "image":
			m.currentPage = PageImage
			m.logger.Debug("switch to image page")
		case "settings":
			m.currentPage = PageSettings
			m.logger.Debug("switch to settings page")
		case "about":
			m.currentPage = PageAbout
			m.logger.Debug("switch to about page")
		}
		return m, nil

	case messages.WelcomeTimeoutMsg:
		// 欢迎页超时，自动跳转到主菜单
		if m.currentPage == PageWelcome {
			m.currentPage = PageMain
			m.logger.Debug("auto switch to main from welcome page")
		}
		return m, nil

	case messages.SaveSuccessMsg:
		m.notification = "✓ 已保存: " + msg.Path
		m.notificationLevel = messages.NotificationSuccess
		return m, m.clearNotificationAfter()

	case messages.SaveFailedMsg:
		m.notification = "✗ 保存失败: " + msg.Err.Error()
		m.notificationLevel = messages.NotificationError
		return m, m.clearNotificationAfter()

	case messages.NotificationMsg:
		m.notification = msg.Message
		m.notificationLevel = msg.Level
		return m, m.clearNotificationAfter()

	case messages.ClearNotificationMsg:
		m.notification = ""
		return m, nil

	case todolist.TasksLoadedMsg:
		if p, ok := m.pages[PageTodo]; ok {
			newModel, cmd := p.Update(msg)
			m.pages[PageTodo] = newModel.(pages.Page)
			return m, cmd
		}
	}

	// 将消息传递给当前页面
	if p, ok := m.pages[m.currentPage]; ok {
		newModel, cmd := p.Update(msg)
		m.pages[m.currentPage] = newModel.(pages.Page)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View 实现 tea.Model 接口
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.currentPage == PageWelcome {
		if p, ok := m.pages[PageWelcome]; ok {
			return p.View()
		}
	}

	// Header
	header := m.renderHeader()

	// Content
	content := m.renderContent()

	// Footer
	footer := m.renderFooter()

	// 分隔线
	separator := styles.SeparatorStyle.Width(m.width).Render(strings.Repeat("─", m.width/2))

	// 组合所有部分
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		separator,
		content,
		separator,
		footer,
	)
}

// renderHeader 渲染头部
func (m Model) renderHeader() string {
	var headerText string

	headerText = fmt.Sprintf("📋 当前页面: %s", m.currentPage)

	if m.notification != "" {
		var notifStyle lipgloss.Style
		switch m.notificationLevel {
		case messages.NotificationSuccess:
			notifStyle = styles.SuccessStyle
		case messages.NotificationWarning:
			notifStyle = styles.WarningStyle
		case messages.NotificationError:
			notifStyle = styles.ErrorStyle
		default:
			notifStyle = styles.InfoStyle
		}
		headerText += " | " + notifStyle.Render(m.notification)
	}

	return styles.HeaderStyle.Width(m.width).Render(headerText)
}

// renderContent 渲染内容
func (m Model) renderContent() string {
	if p, ok := m.pages[m.currentPage]; ok {
		return p.View()
	}
	return m.renderPlaceholder("❓ 未知页面", "页面不存在")
}

// renderFooter 渲染页脚
func (m Model) renderFooter() string {
	sizeInfo := fmt.Sprintf("📐 终端尺寸: %dx%d", m.width, m.height)

	var helpInfo string
	if p, ok := m.pages[m.currentPage]; ok {
		helpInfo = m.help.View(p.Help())
	} else {
		helpInfo = "q/ctrl+c: 退出"
	}

	footer := sizeInfo + "\n" + helpInfo

	return styles.FooterStyle.Width(m.width).Render(footer)
}

// renderPlaceholder 渲染占位页面
func (m Model) renderPlaceholder(title, description string) string {
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
		titleStyle.Render(title),
		descStyle.Render(description),
		backHint,
	)

	// 垂直居中
	contentHeight := lipgloss.Height(content)
	availableHeight := m.getContentHeight()
	topPadding := (availableHeight - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	paddingStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(availableHeight).
		Align(lipgloss.Center, lipgloss.Center)

	return paddingStyle.Render(content)
}

// getContentHeight 获取内容区域高度
func (m Model) getContentHeight() int {
	return m.height - 5
}

// clearNotificationAfter 延时清除通知
func (m *Model) clearNotificationAfter() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return messages.ClearNotificationMsg{}
	})
}
