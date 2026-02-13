# API 参考

本文档提供 KongTools 内部 API 的详细文档。

## 目录

- [Page 接口](#page-接口)
- [消息类型](#消息类型)
- [配置](#配置)
- [日志](#日志)
- [路径管理](#路径管理)
- [版本信息](#版本信息)
- [扩展开发](#扩展开发)

## Page 接口

`Page` 接口是所有 UI 页面的核心抽象。

### 定义

```go
type Page interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
    Title() string
    Help() help.KeyMap
    SetSize(width, height int)
}
```

**位置**：`internal/tui/pages/page.go`

### 方法

#### Init()

```go
Init() tea.Cmd
```

**用途**：初始化页面并返回任何初始命令。

**返回值**：
- `tea.Cmd`：要执行的初始命令（可以是 `nil`）

**示例**：

```go
func (p *MyPage) Init() tea.Cmd {
    // 异步加载数据
    return p.loadDataCmd()
}
```

**调用时机**：页面首次注册或导航到该页面时。

---

#### Update()

```go
Update(tea.Msg) (tea.Model, tea.Cmd)
```

**用途**：处理传入的消息并更新页面状态。

**参数**：
- `msg`：要处理的消息（可以是任何类型）

**返回值**：
- `tea.Model`：更新后的模型（通常是 `p` 本身）
- `tea.Cmd`：要执行的命令（可以是 `nil`）

**示例**：

```go
func (p *MyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // 处理键盘输入
        switch msg.String() {
        case "enter":
            return p, p.handleEnter()
        }
    case tea.WindowSizeMsg:
        p.SetSize(msg.Width, msg.Height)
    }
    return p, nil
}
```

**调用时机**：应用程序事件循环中的每条消息。

---

#### View()

```go
View() string
```

**用途**：将页面渲染为字符串。

**返回值**：
- `string`：渲染后的 UI（通常使用 Lipgloss）

**示例**：

```go
func (p *MyPage) View() string {
    title := styles.TitleStyle.Render("我的页面")
    content := styles.InfoStyle.Render("内容在这里")
    return title + "\n\n" + content
}
```

**调用时机**：每一帧，在 `Update()` 之后。

---

#### Title()

```go
Title() string
```

**用途**：返回用于显示的页面标题。

**返回值**：
- `string`：人类可读的页面标题

**示例**：

```go
func (p *MyPage) Title() string {
    return "我的自定义页面"
}
```

---

#### Help()

```go
Help() help.KeyMap
```

**用途**：返回此页面的键盘快捷键。

**返回值**：
- `help.KeyMap`：用于帮助显示的键绑定

**示例**：

```go
func (p *MyPage) Help() help.KeyMap {
    return p.keys
}
```

---

#### SetSize()

```go
SetSize(width, height int)
```

**用途**：更新页面尺寸。

**参数**：
- `width`：可用宽度（字符数）
- `height`：可用高度（字符数）

**示例**：

```go
func (p *MyPage) SetSize(width, height int) {
    p.width = width
    p.height = height
    p.help.Width = width
    // 重新计算布局...
}
```

**调用时机**：终端调整大小事件时。

---

### 实现自定义页面

完整示例见[扩展开发](#扩展开发)。

## 消息类型

消息用于组件之间的通信。

### 核心消息

**位置**：`internal/tui/messages/messages.go`

#### SwitchPageMsg

导航到不同的页面。

```go
type SwitchPageMsg struct {
    Page string  // 页面标识符："main", "todo", "image", "settings", "about"
}
```

**使用方法**：

```go
// 导航到 Todo 页面
return m, func() tea.Msg {
    return messages.SwitchPageMsg{Page: "todo"}
}
```

---

#### NotificationMsg

向用户显示通知。

```go
type NotificationMsg struct {
    Message string
    Level   NotificationLevel
}

type NotificationLevel int

const (
    NotificationInfo NotificationLevel = iota
    NotificationSuccess
    NotificationWarning
    NotificationError
)
```

**使用方法**：

```go
// 显示成功通知
return m, func() tea.Msg {
    return messages.NotificationMsg{
        Message: "任务保存成功",
        Level:   messages.NotificationSuccess,
    }
}
```

---

#### SaveSuccessMsg / SaveFailedMsg

报告文件保存操作结果。

```go
type SaveSuccessMsg struct {
    Path string
}

type SaveFailedMsg struct {
    Err error
}
```

**使用方法**：

```go
// 在保存命令中
func (m *Model) saveCmd() tea.Msg {
    err := saveToFile(m.data, m.path)
    if err != nil {
        return messages.SaveFailedMsg{Err: err}
    }
    return messages.SaveSuccessMsg{Path: m.path}
}
```

---

#### ClearNotificationMsg

清除当前通知。

```go
type ClearNotificationMsg struct{}
```

**使用方法**：

```go
// 延迟后清除通知
return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
    return messages.ClearNotificationMsg{}
})
```

---

#### WelcomeTimeoutMsg

从欢迎页面自动切换的信号。

```go
type WelcomeTimeoutMsg struct{}
```

**使用方法**：

```go
// 2秒后自动导航
func (p *WelcomePage) Init() tea.Cmd {
    return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
        return messages.WelcomeTimeoutMsg{}
    })
}
```

### 页面特定消息

#### TodoList 消息

**位置**：`internal/tui/pages/todolist/todolist.go`

```go
type TasksLoadedMsg struct {
    Tasks []Task
}

type TasksLoadErrorMsg struct {
    Err error
}

type Task struct {
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

---

#### Image 消息

**位置**：`internal/tui/pages/image/image.go`

```go
type imageMsg struct {
    output string
}

type imageErrMsg struct {
    err error
}

type imageDownloadedMsg struct {
    path string
}
```

## 配置

### 配置结构

**位置**：`internal/config/config.go`

```go
type config struct {
    Log log.Config
    App tui.Config
}
```

### TUI 配置

**位置**：`internal/tui/config.go`

```go
type Config struct {
    TasksSavePath string  // tasks.json 的路径
}
```

### 访问配置

```go
import "kongtools/internal/config"

// 获取单例配置实例
cfg, err := config.Config()
if err != nil {
    log.Error("failed to load config", slog.String("error", err.Error()))
    return err
}

// 访问值
savePath := cfg.App.TasksSavePath
```

### 配置文件

默认配置模板：

```yaml
# KongTools 配置

# 日志配置
log:
  level: info        # debug, info, warn, error
  format: json       # json 或 text
  output: file       # file 或 stdout

# 应用配置
app:
  tasks_save_path: ~/.local/share/kongtools/tasks.json
```

## 日志

### 日志器结构

**位置**：`internal/pkg/log/log.go`

```go
type Config struct {
    Level      string // "debug", "info", "warn", "error"
    Format     string // "json" 或 "text"
    OutputPath string // 文件路径或 "stdout"
    MaxSize    int    // 兆字节
    MaxBackups int    // 备份数量
    MaxAge     int    // 天数
}
```

### 使用日志器

```go
import (
    "log/slog"
    "kongtools/internal/pkg/log"
)

// 创建日志器
logger := log.NewLogger(log.Config{
    Level:  "info",
    Format: "json",
    OutputPath: paths.LogFile("kongtools.log"),
})

// 记录消息
logger.Info("应用程序已启动",
    slog.String("version", "1.0.0"),
    slog.Int("pid", os.Getpid()),
)

logger.Error("加载配置失败",
    slog.String("error", err.Error()),
    slog.String("path", configPath),
)

logger.Debug("处理消息",
    slog.String("type", fmt.Sprintf("%T", msg)),
)
```

### 日志级别

- **Debug**：用于调试的详细信息
- **Info**：常规操作信息
- **Warn**：警告条件
- **Error**：错误条件

## 路径管理

### XDG 路径

**位置**：`internal/pkg/paths/paths.go`

```go
import "kongtools/internal/pkg/paths"

// 配置目录
configDir := paths.ConfigDir()  // ~/.config/kongtools

// 配置文件
configFile := paths.ConfigFile("config.yaml")  // ~/.config/kongtools/config.yaml

// 数据目录
dataDir := paths.DataDir()  // ~/.local/share/kongtools

// 数据文件
dataFile := paths.DataFile("tasks.json")  // ~/.local/share/kongtools/tasks.json

// 状态目录（用于日志）
stateDir := paths.StateDir()  // ~/.local/state/kongtools

// 日志文件
logFile := paths.LogFile("kongtools.log")  // ~/.local/state/kongtools/logs/kongtools.log

// 缓存目录
cacheDir := paths.CacheDir()  // ~/.cache/kongtools
```

### 创建自定义路径

```go
// 确保目录存在
dir := paths.DataDir()
if err := os.MkdirAll(dir, 0755); err != nil {
    log.Error("创建目录失败", slog.String("error", err.Error()))
}

// 在 XDG 路径中创建文件
filePath := paths.DataFile("custom.json")
file, err := os.Create(filePath)
```

## 版本信息

### 版本结构

**位置**：`internal/pkg/version/version.go`

```go
var (
    Version   string = "dev"      // 构建时注入
    Commit    string = "unknown"  // Git 提交哈希
    BuildTime string = "unknown"  // 构建时间戳
    GoVersion string = "unknown"  // 使用的 Go 版本
)
```

### 访问版本信息

```go
import "kongtools/internal/pkg/version"

fmt.Printf("版本: %s\n", version.Version)
fmt.Printf("提交: %s\n", version.Commit)
fmt.Printf("构建时间: %s\n", version.BuildTime)
fmt.Printf("Go 版本: %s\n", version.GoVersion)
```

### 构建时注入

版本信息在构建时通过 ldflags 注入：

```bash
go build -ldflags "\
    -X 'kongtools/internal/pkg/version.Version=v1.0.0' \
    -X 'kongtools/internal/pkg/version.Commit=abc123' \
    -X 'kongtools/internal/pkg/version.BuildTime=2024-01-01' \
    -X 'kongtools/internal/pkg/version.GoVersion=go1.24.0'"
```

## 扩展开发

### 创建自定义页面

添加新页面的完整示例：

#### 1. 定义页面结构

```go
package mypage

import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    tea "github.com/charmbracelet/bubbletea"
    "kongtools/internal/tui/messages"
    "kongtools/internal/tui/styles"
)

// 键绑定
type myKeyMap struct {
    Action key.Binding
    Back   key.Binding
    Quit   key.Binding
}

func (k myKeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Action, k.Back, k.Quit}
}

func (k myKeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Action},
        {k.Back, k.Quit},
    }
}

var keys = myKeyMap{
    Action: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "执行操作"),
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

// 页面模型
type MyPage struct {
    width  int
    height int
    keys   myKeyMap
    help   help.Model
    data   string
}

// 构造函数
func NewMyPage() *MyPage {
    return &MyPage{
        help: help.New(),
        keys: keys,
    }
}
```

#### 2. 实现 Page 接口

```go
// Init 初始化页面
func (p *MyPage) Init() tea.Cmd {
    // 加载数据、启动定时器等
    return nil
}

// Update 处理消息
func (p *MyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, p.keys.Back):
            return p, func() tea.Msg {
                return messages.SwitchPageMsg{Page: "main"}
            }
        case key.Matches(msg, p.keys.Quit):
            return p, tea.Quit
        case key.Matches(msg, p.keys.Action):
            p.data = "已执行操作！"
        }

    case tea.WindowSizeMsg:
        p.SetSize(msg.Width, msg.Height)
    }

    return p, nil
}

// View 渲染页面
func (p *MyPage) View() string {
    title := styles.TitleStyle.Render("我的自定义页面")
    content := styles.InfoStyle.Render(p.data)
    help := p.help.View(p.keys)

    return title + "\n\n" + content + "\n\n" + help
}

// Title 返回页面标题
func (p *MyPage) Title() string {
    return "我的页面"
}

// Help 返回键映射
func (p *MyPage) Help() help.KeyMap {
    return p.keys
}

// SetSize 更新尺寸
func (p *MyPage) SetSize(width, height int) {
    p.width = width
    p.height = height
    p.help.Width = width
}
```

#### 3. 注册页面

编辑 `internal/tui/model.go`：

```go
// 添加到 PageType 常量
const (
    // ... 现有的
    PageMyPage  // 添加这个
)

// 添加到 String() 方法
func (p PageType) String() string {
    switch p {
    // ... 现有的
    case PageMyPage:
        return "我的页面"
    }
    return "未知"
}

// 在 NewModel 中注册
func NewModel(logger *slog.Logger, cfg Config) (*Model, error) {
    // ... 现有代码
    m.pages[PageMyPage] = mypage.NewMyPage()
    return m, nil
}
```

### 创建自定义消息

```go
// 在适当的包中定义
type MyCustomMsg struct {
    Data string
    Code int
}

// 在 Update 中使用
func (p *MyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case MyCustomMsg:
        // 处理自定义消息
        p.data = msg.Data
        return p, nil
    }
    return p, nil
}

// 发送自定义消息
return m, func() tea.Msg {
    return MyCustomMsg{
        Data: "你好",
        Code: 200,
    }
}
```

### 使用命令

```go
// 定义异步命令
func (p *MyPage) loadDataCmd() tea.Cmd {
    return func() tea.Msg {
        // 模拟异步工作
        time.Sleep(1 * time.Second)
        return MyCustomMsg{Data: "已加载！"}
    }
}

// 从 Init 或 Update 返回命令
func (p *MyPage) Init() tea.Cmd {
    return p.loadDataCmd()
}
```

## 最佳实践

### 错误处理

```go
// ❌ 错误：错误时 panic
func loadData() {
    data, err := os.ReadFile("data.json")
    if err != nil {
        panic(err)  // 不要这样做！
    }
}

// ✅ 正确：返回错误消息
func loadData() tea.Msg {
    data, err := os.ReadFile("data.json")
    if err != nil {
        return messages.NotificationMsg{
            Message: "加载数据失败",
            Level:   messages.NotificationError,
        }
    }
    return dataLoadedMsg{Data: data}
}
```

### 性能

```go
// ❌ 错误：在 Update 中进行繁重工作
func (p *MyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 这会阻塞 UI！
    heavyComputation()
    return p, nil
}

// ✅ 正确：使用命令进行异步工作
func (p *MyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return p, p.heavyComputationCmd()
}

func (p *MyPage) heavyComputationCmd() tea.Cmd {
    return func() tea.Msg {
        result := heavyComputation()
        return computationDoneMsg{Result: result}
    }
}
```

### 内存管理

```go
// ✅ 正确：尽可能重用缓冲区
type MyPage struct {
    buffer bytes.Buffer
}

func (p *MyPage) View() string {
    p.buffer.Reset()
    p.buffer.WriteString("内容")
    return p.buffer.String()
}
```

## 测试

测试策略和示例见 [TESTING.md](TESTING.md)。

## 其他资源

- [架构指南](ARCHITECTURE.md) - 系统设计
- [开发指南](DEVELOPMENT.md) - 设置和工作流程
- [Bubbletea 文档](https://github.com/charmbracelet/bubbletea) - 框架文档
- [Lipgloss 文档](https://github.com/charmbracelet/lipgloss) - 样式库

---

如有问题或疑问，请开启 GitHub issue。
