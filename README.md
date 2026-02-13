# KongTools

> 一个用于学习 Bubbletea TUI 框架和归纳 Go 公共库的实践项目

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 项目目的

这个项目主要是为了：

1. **学习 Bubbletea** - 实践 Elm 架构（Model-Update-View）在 Go 中的应用
2. **归纳公共库** - 整理常用的 Go 工具库（log、config、paths 等）
3. **实践 TUI 开发** - 探索终端用户界面的设计和实现

项目包含一些示例功能（待办事项、图片预览等），但这些主要是为了演示技术实现，不是核心目标。

## 技术栈

### TUI 框架
- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** - 基于 Elm 架构的 TUI 框架
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - UI 组件库（list、textinput 等）
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - 样式和布局库

### CLI 工具
- **[Cobra](https://github.com/spf13/cobra)** - CLI 命令框架
- **[Viper](https://github.com/spf13/viper)** - 配置管理

### 公共库（internal/pkg）
- **log** - 基于 slog + lumberjack 的日志系统
- **paths** - XDG 基础目录规范实现
- **version** - 构建时版本信息注入
- **sysinfo** - 系统信息收集

## 快速开始

### 前置要求

- Go 1.24+
- chafa（可选，用于图片预览）

### 运行

```bash
# 直接运行
go run .

# 或构建后运行
make build
./build/kongtools
```

### 开发模式

```bash
# 快速迭代
make dev
```

## 项目结构

```
kongtools/
├── cmd/                    # CLI 命令
├── internal/
│   ├── config/            # 配置管理（基于 Viper）
│   ├── pkg/               # 可复用的公共库
│   │   ├── log/          # 日志系统
│   │   ├── paths/        # XDG 路径
│   │   ├── sysinfo/      # 系统信息
│   │   └── version/      # 版本注入
│   └── tui/               # TUI 应用核心
│       ├── pages/         # 各个页面实现
│       ├── styles/        # Lipgloss 样式
│       ├── messages/      # 消息类型定义
│       └── keys/          # 键绑定
├── main.go
├── Makefile
└── go.mod
```

## 核心设计模式

### 1. Elm 架构

```
用户输入 → Message → Update(Model, Cmd) → View() → 渲染
```

### 2. Page 接口

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

### 3. 消息驱动

```go
// 页面切换
type SwitchPageMsg struct {
    Page string
}

// 通知
type NotificationMsg struct {
    Message string
    Level   NotificationLevel
}
```

## 学习要点

### Bubbletea 核心概念

1. **Model** - 应用状态
2. **Update** - 状态转换函数
3. **View** - 渲染函数
4. **Cmd** - 异步操作

### 最佳实践

1. **错误处理** - 不要用 panic，返回 error
2. **性能优化** - 增量更新而非全量重建
3. **代码复用** - 提取通用组件
4. **安全编程** - 超时、大小限制

## 文档

- [架构指南](docs/ARCHITECTURE.md) - 系统设计和设计决策
- [开发指南](docs/DEVELOPMENT.md) - 如何扩展
- [API 参考](docs/API.md) - 内部 API 文档
- [重构报告](docs/REFACTORING_REPORT.md) - 重构过程和成果

## 参考资料

### Bubbletea 相关
- [Bubbletea 官方教程](https://github.com/charmbracelet/bubbletea#tutorial)
- [Elm 架构](https://guide.elm-lang.org/architecture/)
- [Lipgloss 样式指南](https://github.com/charmbracelet/lipgloss)

### Go 最佳实践
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [XDG 基础目录规范](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)

## 构建和测试

```bash
# 构建
make build

# 运行测试
make test

# 清理
make clean
```

## 许可证

MIT

---

**注意**：这是一个学习项目，主要用于探索和实践 Go TUI 开发。
