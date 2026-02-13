# 架构指南

本文档记录了学习 Bubbletea 过程中的架构设计和设计决策。

## 核心概念：Elm 架构

Bubbletea 采用 **Elm 架构**，这是一个函数式的、可预测的 UI 架构模式：

```
┌──────────┐
│  Model   │ ← 应用状态（不可变）
└────┬─────┘
     │
     ▼
┌──────────┐
│  Update  │ ← Message → (Model, Cmd)
└────┬─────┘
     │
     ▼
┌──────────┐
│   View   │ → String（终端输出）
└──────────┘
```

**示例流程**：
1. 用户按键 → 创建 `tea.KeyMsg`
2. `Update` 函数接收消息
3. 根据消息更新 Model
4. 返回可选的命令（`tea.Cmd`）
5. `View` 函数渲染新状态

## 核心接口：Page

所有页面实现统一的 `Page` 接口：

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

**设计优点**：
- **多态性**: 可以在不更改核心逻辑的情况下交换页面
- **可测试性**: 每个页面可以独立测试
- **可扩展性**: 易于添加新页面

## 消息系统

应用使用类型化消息进行通信：

### 核心消息

```go
// 页面导航
type SwitchPageMsg struct {
    Page string  // "main", "todo", "image", "settings", "about"
}

// 通知
type NotificationMsg struct {
    Message string
    Level   NotificationLevel
}

// 文件操作
type SaveSuccessMsg struct {
    Path string
}
```

**优点**：
- **类型安全**: 编译器捕获消息不匹配
- **清晰性**: 明确的消息流
- **可调试性**: 易于跟踪消息路径

## 目录结构

```
kongtools/
├── cmd/                    # CLI 命令（Cobra）
├── internal/
│   ├── config/            # 配置管理（基于 Viper）
│   ├── pkg/               # 可复用公共库
│   │   ├── log/          # slog + lumberjack
│   │   ├── paths/        # XDG 路径管理
│   │   └── version/      # 构建时版本注入
│   └── tui/               # TUI 应用核心
│       ├── pages/         # 各个页面实现
│       ├── styles/        # Lipgloss 样式
│       ├── messages/      # 消息类型定义
│       └── keys/          # 键绑定
├── main.go
└── Makefile
```

## 关键设计决策

### 1. 为什么选择 Bubbletea？

**理由**：
1. **Elm 架构**: 更可预测的状态管理
2. **函数式方法**: 更易于测试和推理
3. **可组合性**: 易于组合组件
4. **性能**: 通过 diffing 实现高效渲染

**学习价值**：
- 理解函数式编程在 Go 中的应用
- 掌握不可变状态管理
- 学习响应式 UI 设计

### 2. 为什么使用页面映射而不是继承？

**决策**: 使用 `map[PageType]Page`

**理由**：
1. **灵活性**: 页面可以动态添加/删除
2. **解耦**: 页面不需要相互了解
3. **可测试性**: 每个页面可以独立测试

**实现**：
```go
type Model struct {
    currentPage PageType
    pages       map[PageType]pages.Page
}
```

### 3. 为什么使用 XDG 基础目录规范？

**理由**：
1. **可预测性**: 用户知道文件在哪里
2. **整洁的主目录**: 避免用点文件弄乱 `~/`
3. **可移植性**: 跨 Linux 发行版工作

### 4. 为什么使用结构化日志？

**决策**: 使用 `log/slog` + JSON 输出

**理由**：
1. **机器可解析**: 易于分析日志
2. **结构化数据**: 用于上下文的键值对
3. **性能**: 即使有详细日志也很高效

## 学习要点

### 1. 状态管理

```go
// ❌ 错误：在 Update 中进行繁重工作
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    heavyComputation()  // 这会阻塞 UI！
    return m, nil
}

// ✅ 正确：使用 Cmd 进行异步工作
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, heavyComputationCmd()
}
```

### 2. 错误处理

```go
// ❌ 错误：panic
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

### 3. 性能优化

```go
// 增量更新而非全量重建
func (m *Model) updateListItems(changedIndices ...int) {
    if len(changedIndices) == 0 {
        // 全量更新
        items := make([]list.Item, len(m.tasks))
        for i, t := range m.tasks {
            items[i] = item{task: &t}
        }
        m.list.SetItems(items)
    } else {
        // 增量更新：只更新改变的项
        for _, idx := range changedIndices {
            m.list.SetItem(idx, item{task: &m.tasks[idx]})
        }
    }
}
```

## 扩展点

### 添加新页面

1. 在 `internal/tui/pages/newpage/` 创建新包
2. 实现 `Page` 接口
3. 在 `internal/tui/model.go` 中注册：
   ```go
   m.pages[PageNew] = newpage.NewPage()
   ```
4. 添加导航消息处理器

### 添加新消息类型

1. 在 `internal/tui/messages/messages.go` 中定义：
   ```go
   type NewMsg struct {
       Data string
   }
   ```
2. 在适当的 `Update()` 方法中处理

## 性能考虑

### 消息处理

- 消息在 Update 循环中同步处理
- 长时间运行的操作应返回 `tea.Cmd` 进行异步执行

### 渲染

- Bubbletea 执行高效的 diffing
- 避免在 View() 中进行不必要的分配
- 尽可能缓存计算样式

## 参考资源

- [Bubbletea 文档](https://github.com/charmbracelet/bubbletea)
- [Elm 架构教程](https://guide.elm-lang.org/architecture/)
- [Lipgloss 样式](https://github.com/charmbracelet/lipgloss)

---

实现细节见 [DEVELOPMENT.md](DEVELOPMENT.md)。
API 参考见 [API.md](API.md)。
