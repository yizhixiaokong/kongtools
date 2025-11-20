package todolist

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"kongtools/internal/tui/messages"
	"kongtools/internal/tui/styles"
)

// Task 任务结构
type Task struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ListPage Todo List 页面
type ListPage struct {
	// 数据
	tasks    []Task
	savePath string
	mutex    sync.Mutex
	logger   *slog.Logger

	// UI 状态
	width     int
	height    int
	input     string
	selected  int
	inputMode bool // 是否处于输入模式（焦点在输入框）
	editMode  bool
	editIndex int
	hint      string

	// 快捷键
	keys todoKeyMap
}

// type todoKeyMap Todo 页面快捷键映射
type todoKeyMap struct {
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Toggle  key.Binding
	Back    key.Binding
	Quit    key.Binding
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
	Input   key.Binding
}

// ShortHelp 返回简短帮助信息
func (k todoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Edit, k.Toggle, k.Delete, k.Back}
}

// FullHelp 返回完整帮助信息
func (k todoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Toggle},
		{k.Edit, k.Delete, k.Add},
		{k.Back, k.Quit},
	}
}

var todoKeys = todoKeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "进入输入模式"),
	),
	Edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "编辑"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete", "d"),
		key.WithHelp("d/del", "删除"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "完成/未完成"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "返回"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "退出"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "上移"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "下移"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "确认"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "取消"),
	),
}

// NewListPage 创建 Todo List 页面
func NewListPage(logger *slog.Logger, savePath string) *ListPage {
	return &ListPage{
		tasks:     []Task{},
		savePath:  savePath,
		editMode:  false,
		editIndex: -1,
		selected:  0,
		logger:    logger.With("module", "tui-todo"),
		keys:      todoKeys,
	}
}

// Init 实现 Page 接口
func (m *ListPage) Init() tea.Cmd {
	if err := m.LoadTasks(); err != nil {
		m.logger.Error("failed to load tasks", slog.String("error", err.Error()))
	}
	return nil
}

// Update 实现 Page 接口
func (m *ListPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.MouseMsg:
		// 处理鼠标滚动
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				if m.selected > 0 {
					m.selected--
				}
			case tea.MouseButtonWheelDown:
				if m.selected < len(m.tasks)-1 {
					m.selected++
				}
			}
		}

	case tea.KeyMsg:
		// 输入模式的按键处理
		if m.inputMode {
			return m.handleInputKeys(msg)
		}

		// 列表模式的按键处理
		return m.handleListKeys(msg)

	case messages.SaveSuccessMsg:
		m.hint = "✓ 已保存: " + msg.Path
		cmds = append(cmds, m.clearHintAfter(3*time.Second))

	case messages.SaveFailedMsg:
		m.hint = "✗ 保存失败: " + msg.Err.Error()
		cmds = append(cmds, m.clearHintAfter(3*time.Second))

	case messages.ClearHintMsg:
		m.hint = ""
	}

	return m, tea.Batch(cmds...)
}

// handleInputKeys 处理输入模式的按键
func (m *ListPage) handleInputKeys(msg tea.KeyMsg) (*ListPage, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.Type {
	case tea.KeyEnter:
		if m.editMode {
			cmds = append(cmds, m.saveEdit())
		} else if m.input != "" {
			cmds = append(cmds, m.addTask())
		}
		// 添加/编辑完成后退出输入模式
		m.inputMode = false

	case tea.KeyEsc:
		// ESC 退出输入模式
		m.cancelEdit()
		m.inputMode = false

	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	case tea.KeySpace:
		// 在输入模式下，空格键就是输入空格，不触发 Toggle
		if len(m.input) < 80 {
			m.input += " "
		}

	case tea.KeyRunes:
		if len(m.input) < 80 {
			m.input += string(msg.Runes)
		} else {
			m.hint = "⚠ 任务长度不能超过 80 个字符"
			cmds = append(cmds, m.clearHintAfter(3*time.Second))
		}
	}

	return m, tea.Batch(cmds...)
}

// handleListKeys 处理列表模式的按键
func (m *ListPage) handleListKeys(msg tea.KeyMsg) (*ListPage, tea.Cmd) {
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Add):
		// 按 'a' 进入输入模式
		m.inputMode = true
		m.input = ""
		m.editMode = false
		m.hint = "💡 输入新任务内容，按 Enter 确认，ESC 取消"
		cmds = append(cmds, m.clearHintAfter(5*time.Second))

	case key.Matches(msg, m.keys.Up):
		if m.selected > 0 {
			m.selected--
		}

	case key.Matches(msg, m.keys.Down):
		if m.selected < len(m.tasks)-1 {
			m.selected++
		}

	case key.Matches(msg, m.keys.Edit):
		m.editTask()
		m.inputMode = true

	case key.Matches(msg, m.keys.Toggle):
		cmds = append(cmds, m.toggleComplete())

	case key.Matches(msg, m.keys.Delete):
		cmds = append(cmds, m.deleteTask())

	case key.Matches(msg, m.keys.Back):
		// 返回主菜单
		return m, func() tea.Msg {
			return messages.SwitchPageMsg{Page: "main"}
		}

	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

// View 实现 Page 接口
func (m *ListPage) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// 定义列表的最大宽度（用于居中）
	const maxListWidth = 80
	listWidth := maxListWidth
	if m.width < maxListWidth {
		listWidth = m.width - 4 // 留一些边距
	}

	// 标题部分（居中）
	title := lipgloss.NewStyle().
		Width(listWidth).
		Align(lipgloss.Center).
		Render(styles.TitleStyle.Render("📋 Todo List"))

	// 输入框部分（左对齐）
	var inputLine string
	if m.inputMode {
		// 在输入模式下显示输入框和光标
		inputLabel := "＋ 新任务: "
		if m.editMode {
			inputLabel = "✎ 编辑: "
		}

		labelStyle := styles.TodoInputLabelStyle
		inputStyle := styles.TodoInputStyle
		cursorStyle := styles.TodoCursorStyle

		inputText := m.input
		if len(inputText) < 80 {
			inputText += cursorStyle.Render("▊") // 使用更明显的光标
		}

		inputLine = labelStyle.Render(inputLabel) + inputStyle.Render(inputText)
	} else {
		// 非输入模式，显示提示
		hintStyle := lipgloss.NewStyle().Foreground(styles.TextSubtle)
		inputLine = hintStyle.Render("按 a 进入输入模式添加任务")
	}

	// 提示信息（左对齐）
	var hintLine string
	if m.hint != "" {
		hintStyle := styles.TodoHintStyle
		hintLine = hintStyle.Render(m.hint)
	}

	// 任务列表
	listHeight := m.height - 5
	if listHeight < 0 {
		listHeight = 0
	}

	// 计算显示范围
	start := 0
	end := len(m.tasks)

	if end > listHeight {
		if m.selected >= listHeight {
			start = m.selected - listHeight + 1
		}
		end = start + listHeight
		if end > len(m.tasks) {
			end = len(m.tasks)
			start = end - listHeight
			if start < 0 {
				start = 0
			}
		}
	}

	// 构建任务列表项（左对齐）
	var taskLines []string
	for i := start; i < end; i++ {
		task := m.tasks[i]

		checkbox := "[ ]"
		if task.Completed {
			checkbox = "[✓]"
		}

		line := fmt.Sprintf("%s %s", checkbox, task.Title)

		var renderedLine string
		if i == m.selected {
			if task.Completed {
				renderedLine = styles.TodoSelectedCompletedStyle.Render(line)
			} else {
				renderedLine = styles.TodoSelectedStyle.Render(line)
			}
		} else {
			if task.Completed {
				renderedLine = styles.TodoCompletedStyle.Render(line)
			} else {
				renderedLine = styles.TodoItemStyle.Render(line)
			}
		}

		taskLines = append(taskLines, renderedLine)
	}

	// 组合所有部分，确保任务列表左对齐
	var contentParts []string
	contentParts = append(contentParts, title)
	contentParts = append(contentParts, inputLine)
	if hintLine != "" {
		contentParts = append(contentParts, hintLine)
	} else {
		contentParts = append(contentParts, "")
	}
	contentParts = append(contentParts, taskLines...)

	// 使用固定宽度的容器，内部左对齐
	contentBox := lipgloss.NewStyle().
		Width(listWidth).
		Align(lipgloss.Left).
		Render(strings.Join(contentParts, "\n"))

	// 使用 lipgloss.Place 将内容块水平居中
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		contentBox,
	)
}

// Title 实现 Page 接口
func (m *ListPage) Title() string {
	return "待办列表"
}

// Help 实现 Page 接口
func (m *ListPage) Help() help.KeyMap {
	return m.keys
}

// SetSize 实现 Page 接口
func (m *ListPage) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// LoadTasks 加载任务
func (m *ListPage) LoadTasks() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data, err := os.ReadFile(m.savePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Debug("no tasks file found, creating with help tasks")
			m.tasks = m.getHelpTasks()
			if saveErr := m.saveTasksLocked(); saveErr != nil {
				m.logger.Error("failed to save initial help tasks", slog.String("error", saveErr.Error()))
			}
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &m.tasks); err != nil {
		return err
	}

	m.logger.Info("tasks loaded", slog.String("path", m.savePath), slog.Int("count", len(m.tasks)))
	return nil
}

// SaveTasks 保存任务
func (m *ListPage) SaveTasks() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.saveTasksLocked()
}

// saveTasksLocked 保存任务（已加锁版本）
func (m *ListPage) saveTasksLocked() error {
	data, err := json.MarshalIndent(m.tasks, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(m.savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(m.savePath, data, 0644); err != nil {
		return err
	}

	m.logger.Info("tasks saved", slog.String("path", m.savePath), slog.Int("count", len(m.tasks)))
	return nil
}

// scheduleSave 计划保存
func (m *ListPage) scheduleSave() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		if err := m.SaveTasks(); err != nil {
			m.logger.Error("failed to save tasks", slog.String("error", err.Error()))
			return messages.SaveFailedMsg{Err: err}
		}
		return messages.SaveSuccessMsg{Path: m.savePath}
	}
}

// getHelpTasks 获取帮助任务
func (m *ListPage) getHelpTasks() []Task {
	return []Task{
		{Title: "💡 按 a 进入输入模式添加任务", Completed: false},
		{Title: "👏 输入完成后按 Enter 确认，ESC 取消", Completed: false},
		{Title: "📝 选中任务并按 Enter 编辑任务", Completed: false},
		{Title: "❌ 按 Delete 或 x 删除选中的任务", Completed: false},
		{Title: "✅ 按空格键标记任务为已完成/未完成", Completed: false},
	}
}

// addTask 添加任务
func (m *ListPage) addTask() tea.Cmd {
	if m.input == "" {
		return nil
	}

	task := Task{
		Title:     m.input,
		Completed: false,
	}

	m.tasks = append(m.tasks, task)
	m.input = ""
	m.selected = len(m.tasks) - 1
	m.logger.Debug("task added", slog.String("title", task.Title))

	return m.scheduleSave()
}

// editTask 编辑任务
func (m *ListPage) editTask() {
	if len(m.tasks) == 0 || m.selected >= len(m.tasks) {
		return
	}

	m.input = m.tasks[m.selected].Title
	m.editMode = true
	m.editIndex = m.selected
	m.logger.Debug("editing task", slog.Int("index", m.editIndex))
}

// saveEdit 保存编辑
func (m *ListPage) saveEdit() tea.Cmd {
	if m.input == "" {
		m.cancelEdit()
		return nil
	}

	if m.editIndex >= 0 && m.editIndex < len(m.tasks) {
		m.tasks[m.editIndex].Title = m.input
		m.logger.Debug("task edited", slog.Int("index", m.editIndex), slog.String("title", m.input))
	}

	m.input = ""
	m.editMode = false
	m.editIndex = -1

	return m.scheduleSave()
}

// cancelEdit 取消编辑
func (m *ListPage) cancelEdit() {
	m.input = ""
	m.editMode = false
	m.editIndex = -1
	m.logger.Debug("edit cancelled")
}

// deleteTask 删除任务
func (m *ListPage) deleteTask() tea.Cmd {
	if len(m.tasks) == 0 || m.selected >= len(m.tasks) {
		return nil
	}

	if m.editMode {
		m.hint = "⚠ 编辑模式下不能删除任务"
		return m.clearHintAfter(3 * time.Second)
	}

	task := m.tasks[m.selected]
	m.tasks = append(m.tasks[:m.selected], m.tasks[m.selected+1:]...)

	if m.selected >= len(m.tasks) && m.selected > 0 {
		m.selected--
	}

	m.logger.Debug("task deleted", slog.String("title", task.Title))

	return m.scheduleSave()
}

// toggleComplete 切换完成状态
func (m *ListPage) toggleComplete() tea.Cmd {
	if len(m.tasks) == 0 || m.selected >= len(m.tasks) {
		return nil
	}

	m.tasks[m.selected].Completed = !m.tasks[m.selected].Completed
	m.logger.Debug("task completion toggled",
		slog.Int("index", m.selected),
		slog.Bool("completed", m.tasks[m.selected].Completed))

	return m.scheduleSave()
}

// clearHintAfter 延时清除提示
func (m *ListPage) clearHintAfter(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return messages.ClearHintMsg{}
	})
}
