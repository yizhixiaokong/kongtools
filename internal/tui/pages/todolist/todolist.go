package todolist

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

// TasksLoadedMsg 任务加载成功消息
type TasksLoadedMsg struct {
	Tasks []Task
}

// item 实现 list.Item 接口
type item struct {
	task *Task
}

func (i item) Title() string       { return i.task.Title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.task.Title }

// itemDelegate 列表项渲染委托
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Title())

	// 选中状态
	fn := styles.TodoItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.TodoSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	} else {
		fn = func(s ...string) string {
			return styles.TodoItemStyle.Render("  " + strings.Join(s, " "))
		}
	}

	// 完成状态
	checkbox := "[ ]"
	if i.task.Completed {
		checkbox = "[✓]"
		str = styles.TodoCompletedStyle.Render(str)
	}

	fmt.Fprint(w, fn(checkbox, str))
}

// ListPage Todo List 页面
type ListPage struct {
	// 数据
	tasks    []Task
	savePath string
	mutex    sync.Mutex
	logger   *slog.Logger

	// UI 组件
	list  list.Model
	input textinput.Model

	// 状态
	width     int
	height    int
	adding    bool // 是否正在添加/编辑
	editIndex int  // -1 表示添加，>=0 表示编辑

	// 快捷键
	keys todoKeyMap
}

// type todoKeyMap Todo 页面快捷键映射
type todoKeyMap struct {
	Add    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Toggle key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// ShortHelp 返回简短帮助信息
func (k todoKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Edit, k.Toggle, k.Delete, k.Back}
}

// FullHelp 返回完整帮助信息
func (k todoKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Edit, k.Delete},
		{k.Toggle, k.Back, k.Quit},
	}
}

var todoKeys = todoKeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "添加任务"),
	),
	Edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "编辑任务"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete", "d"),
		key.WithHelp("d/del", "删除任务"),
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
}

// NewListPage 创建 Todo List 页面
func NewListPage(logger *slog.Logger, savePath string) *ListPage {
	// 初始化列表
	delegate := itemDelegate{}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "📋 Todo List"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // 我们使用自己的帮助系统或页面底部的帮助
	l.Styles.Title = styles.TitleStyle
	l.DisableQuitKeybindings()

	// 初始化输入框
	ti := textinput.New()
	ti.Placeholder = "输入任务内容..."
	ti.CharLimit = 80
	ti.Width = 40

	return &ListPage{
		tasks:     []Task{},
		savePath:  savePath,
		list:      l,
		input:     ti,
		adding:    false,
		editIndex: -1,
		logger:    logger.With("module", "tui-todo"),
		keys:      todoKeys,
	}
}

// Init 实现 Page 接口
func (m *ListPage) Init() tea.Cmd {
	return m.loadTasksCmd
}

// loadTasksCmd 加载任务
func (m *ListPage) loadTasksCmd() tea.Msg {
	data, err := os.ReadFile(m.savePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Debug("no tasks file found, returning help tasks")
			return TasksLoadedMsg{Tasks: m.getHelpTasks()}
		}
		m.logger.Error("failed to load tasks", slog.String("error", err.Error()))
		return nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		m.logger.Error("failed to unmarshal tasks", slog.String("error", err.Error()))
		return nil
	}

	m.logger.Info("tasks loaded", slog.String("path", m.savePath), slog.Int("count", len(tasks)))
	return TasksLoadedMsg{Tasks: tasks}
}

// Update 实现 Page 接口
func (m *ListPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4) // 留出 header/footer 空间

	case TasksLoadedMsg:
		m.tasks = msg.Tasks
		m.updateListItems()

	case messages.SaveSuccessMsg:
		// 可以显示通知，或者什么都不做
	case messages.SaveFailedMsg:
		// 可以显示错误通知
	}

	// 如果正在添加/编辑，处理输入框逻辑
	if m.adding {
		return m.updateInput(msg)
	}

	// 否则处理列表逻辑
	return m.updateList(msg)
}

// updateInput 处理输入模式的更新
func (m *ListPage) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.input.Value() != "" {
				if m.editIndex >= 0 {
					// 编辑现有任务
					m.tasks[m.editIndex].Title = m.input.Value()
				} else {
					// 添加新任务
					m.tasks = append(m.tasks, Task{Title: m.input.Value()})
				}
				m.updateListItems()
				cmd = m.scheduleSave()
			}
			m.adding = false
			m.input.Blur()
			return m, cmd

		case tea.KeyEsc:
			m.adding = false
			m.input.Blur()
			return m, nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// updateList 处理列表模式的更新
func (m *ListPage) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 检查是否匹配自定义快捷键
		if key.Matches(msg, m.keys.Add) {
			m.adding = true
			m.editIndex = -1
			m.input.SetValue("")
			m.input.Focus()
			return m, textinput.Blink
		}

		if key.Matches(msg, m.keys.Edit) {
			if len(m.tasks) > 0 && m.list.Index() >= 0 {
				m.adding = true
				m.editIndex = m.list.Index()
				m.input.SetValue(m.tasks[m.editIndex].Title)
				m.input.Focus()
				return m, textinput.Blink
			}
		}

		if key.Matches(msg, m.keys.Delete) {
			if len(m.tasks) > 0 && m.list.Index() >= 0 {
				index := m.list.Index()
				m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
				m.updateListItems()
				// 调整选中项
				if index >= len(m.tasks) && index > 0 {
					m.list.Select(index - 1)
				}
				return m, m.scheduleSave()
			}
		}

		if key.Matches(msg, m.keys.Toggle) {
			if len(m.tasks) > 0 && m.list.Index() >= 0 {
				index := m.list.Index()
				m.tasks[index].Completed = !m.tasks[index].Completed
				// 重新生成列表项以更新显示
				m.updateListItems()
				// 保持选中位置
				m.list.Select(index)
				return m, m.scheduleSave()
			}
		}

		if key.Matches(msg, m.keys.Back) {
			return m, func() tea.Msg {
				return messages.SwitchPageMsg{Page: "main"}
			}
		}

		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// updateListItems 更新列表组件的数据
func (m *ListPage) updateListItems() {
	items := make([]list.Item, len(m.tasks))
	for i, t := range m.tasks {
		// 注意：这里需要传递指针，否则修改不会反映到原始切片
		// 但由于我们每次都重新生成 items，所以直接传值也可以，
		// 只要保证 m.tasks 是最新的。
		// 为了在 item 方法中访问 Task，我们创建一个新的 Task 副本或指针
		taskCopy := t // 复制一份
		items[i] = item{task: &taskCopy}
	}
	m.list.SetItems(items)
}

// View 实现 Page 接口
func (m *ListPage) View() string {
	if m.adding {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				styles.TitleStyle.Render(m.inputTitle()),
				m.input.View(),
				styles.SubtitleStyle.Render("(Enter 确认, Esc 取消)"),
			),
		)
	}

	// 使用 lipgloss.Place 居中显示列表
	// 注意：list 组件自带了分页和样式，我们只需要给它足够的空间
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		m.list.View(),
	)
}

func (m *ListPage) inputTitle() string {
	if m.editIndex >= 0 {
		return "编辑任务"
	}
	return "添加新任务"
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
	m.list.SetWidth(width)
	m.list.SetHeight(height - 4)
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
		{Title: "💡 按 a 添加任务", Completed: false},
		{Title: "📝 按 Enter 编辑任务", Completed: false},
		{Title: "❌ 按 d 删除任务", Completed: false},
		{Title: "✅ 按空格键切换完成状态", Completed: false},
	}
}
