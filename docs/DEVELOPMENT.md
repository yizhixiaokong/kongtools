# 开发指南

本文档帮助你理解项目结构，以及如何扩展这个学习项目。

## 项目结构

```
kongtools/
├── cmd/                    # CLI 命令（Cobra）
│   ├── root.go            # 根命令和初始化
│   └── version.go         # 版本子命令
│
├── internal/              # 私有应用代码
│   ├── config/           # 配置管理（基于 Viper）
│   │
│   ├── pkg/              # 可复用的公共库
│   │   ├── log/         # 日志系统（slog + lumberjack）
│   │   ├── paths/       # XDG 路径管理
│   │   ├── sysinfo/     # 系统信息收集
│   │   └── version/     # 构建时版本注入
│   │
│   └── tui/             # TUI 应用核心
│       ├── pages/        # 各个页面实现
│       │   ├── page.go  # Page 接口定义
│       │   ├── welcome.go
│       │   ├── main.go
│       │   ├── todolist/
│       │   ├── image/
│       │   ├── settings/
│       │   └── about/
│       ├── styles/       # Lipgloss 样式
│       ├── messages/     # 消息类型定义
│       └── keys/         # 键绑定
│
├── main.go               # 应用入口
├── Makefile              # 构建自动化
└── go.mod                # Go 模块定义
```

## 添加新页面

添加新页面是最常见的任务之一。

### 分步指南

#### 1. 创建页面包

在 `internal/tui/pages/` 中创建新目录：

```bash
mkdir internal/tui/pages/newpage
```

#### 2. 实现 Page 接口

创建 `internal/tui/pages/newpage/newpage.go`：

```go
package newpage

import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    tea "github.com/charmbracelet/bubbletea"
    "kongtools/internal/tui/messages"
    "kongtools/internal/tui/styles"
)

// KeyMap 定义键盘快捷键
type newPageKeyMap struct {
    Action key.Binding
    Back   key.Binding
    Quit   key.Binding
}

func (k newPageKeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Action, k.Back, k.Quit}
}

func (k newPageKeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Action},
        {k.Back, k.Quit},
    }
}

var keys = newPageKeyMap{
    Action: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "执行操作"),
    ),
    Back: key.NewBinding(
        key.WithKeys("esc"),
        key.WithHelp("esc", "返回主菜单"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("ctrl+c"),
        key.WithHelp("ctrl+c", "退出"),
    ),
}

// NewPage 表示新页面
type NewPage struct {
    width  int
    height int
    keys   newPageKeyMap
    help   help.Model
}

// NewNewPage 创建新的 NewPage 实例
func NewNewPage() *NewPage {
    return &NewPage{
        help: help.New(),
        keys: keys,
    }
}

// Init 初始化页面
func (p *NewPage) Init() tea.Cmd {
    return nil
}

// Update 处理消息
func (p *NewPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
            // 处理操作
        }

    case tea.WindowSizeMsg:
        p.width = msg.Width
        p.height = msg.Height
        p.help.Width = msg.Width
    }

    return p, nil
}

// View 渲染页面
func (p *NewPage) View() string {
    // 使用 Lipgloss 构建 UI
    title := styles.TitleStyle.Render("新页面")
    content := styles.SubtitleStyle.Render("你的内容在这里")

    return title + "\n\n" + content
}

// Title 返回页面标题
func (p *NewPage) Title() string {
    return "新页面"
}

// Help 返回 keymap 用于帮助显示
func (p *NewPage) Help() help.KeyMap {
    return p.keys
}

// SetSize 设置页面尺寸
func (p *NewPage) SetSize(width, height int) {
    p.width = width
    p.height = height
    p.help.Width = width
}
```

#### 3. 注册页面

编辑 `internal/tui/model.go`:

```go
// 添加到 PageType 常量
const (
    // ... 现有页面
    PageNewPage  // 添加这个
)

// 添加到 String() 方法
func (p PageType) String() string {
    switch p {
    // ... 现有情况
    case PageNewPage:
        return "新页面"  // 添加这个
    }
}

// 在 NewModel 函数中注册
func NewModel(logger *slog.Logger, cfg Config) (*Model, error) {
    // ... 现有代码

    m.pages[PageNewPage] = newpage.NewNewPage()  // 添加这个

    return m, nil
}
```

#### 4. 添加导航

更新主菜单或其他页面以允许导航到新页面：

```go
// 在 main.go 或适当的地方
case key.Matches(msg, m.keys.NewPage):
    return m, func() tea.Msg {
        return messages.SwitchPageMsg{Page: "newpage"}
    }
```

#### 5. 测试页面

```bash
# 运行应用
make dev

# 导航到新页面并测试所有功能
```

## 代码规范

### 格式化

始终用 `gofmt` 格式化代码：

```bash
gofmt -w .
```

### 错误处理

**不要**使用 `cobra.CheckErr()` - 它会 panic。而是返回错误：

```go
// ❌ 错误
func someFunction() {
    if err := doSomething(); err != nil {
        cobra.CheckErr(err)  // 这会 panic！
    }
}

// ✅ 正确
func someFunction() error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

### 日志

使用 `slog` 进行结构化日志：

```go
// ❌ 错误
log.Println("Something happened")

// ✅ 正确
logger.Info("task completed",
    slog.String("task_id", id),
    slog.Duration("duration", elapsed),
)
```

### 注释

- 优先使用英文注释
- 文档化公共函数和类型
- 解释"为什么"，而不是"什么"

```go
// ❌ 错误
// SetWidth sets the width
func SetWidth(w int) {
    width = w
}

// ✅ 正确
// SetWidth updates the viewport width and recalculates the layout.
// This ensures proper text wrapping in the UI.
func SetWidth(w int) {
    width = w
    recalculateLayout()
}
```

### 命名约定

- 使用描述性名称
- 遵循 Go 命名约定
- 避免缩写，除非显而易见

```go
// ❌ 错误
var n int
func calc() int { ... }

// ✅ 正确
var taskCount int
func calculateTotalDuration() int { ... }
```

## 测试

### 运行测试

```bash
# 运行所有测试
make test

# 带覆盖率运行
go test -cover ./...

# 运行特定包测试
go test ./internal/tui/...

# 带详细输出运行
go test -v ./...
```

### 编写测试

在相同包中放置测试，使用 `_test.go` 后缀：

```go
// internal/tui/model_test.go
package tui

import "testing"

func TestNewModel(t *testing.T) {
    model, err := NewModel(logger, cfg)
    if err != nil {
        t.Fatalf("NewModel failed: %v", err)
    }

    if model.currentPage != PageWelcome {
        t.Errorf("Expected PageWelcome, got %v", model.currentPage)
    }
}
```

### 测试覆盖率

我们的目标是 **80%+ 测试覆盖率**。检查覆盖率：

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 调试

### 查看日志

日志写入到：

```
~/.local/state/kongtools/logs/kongtools.log
```

实时查看日志：

```bash
tail -f ~/.local/state/kongtools/logs/kongtools.log
```

### 常见问题

**问题**: 应用无法启动

**解决**: 检查日志中的错误。验证配置文件存在。

```bash
cat ~/.local/state/kongtools/logs/kongtools.log
```

**问题**: 图片预览不工作

**解决**: 确保 `chafa` 已安装并在 PATH 中。

```bash
which chafa
```

**问题**: 配置文件错误

**解决**: 删除并重新生成配置：

```bash
rm ~/.config/kongtools/config.yaml
# 重启应用
```

## 构建 & 发布

### 构建

```bash
# 带版本信息构建
make build

# 输出: build/kongtools
```

### 版本信息

构建过程注入版本信息：

```bash
make version
# 显示: Version, Commit, Build Time, Go Version
```

### 创建发布

```bash
# 更新版本
# 标记发布
make release TAG=v1.0.0

# 这将：
# 1. 生成 CHANGELOG.md
# 2. 提交变更日志
# 3. 创建 git 标签
# 4. 推送说明
```

## 其他资源

- [架构指南](ARCHITECTURE.md) - 系统设计和组件
- [API 参考](API.md) - 内部 API 文档
- [用户指南](USER_GUIDE.md) - 功能文档
- [测试指南](TESTING.md) - 详细测试策略

## 获取帮助

- **Issues**: 为 bug 或功能请求开启 GitHub issue
- **Discussions**: 使用 GitHub Discussions 提问
- **Code Review**: 所有 PR 在合并前需要审查

## 许可证

通过为 KongTools 做贡献，你同意你的贡献将按照 MIT 许可证授权。

---

感谢为 KongTools 做贡献！🎉
