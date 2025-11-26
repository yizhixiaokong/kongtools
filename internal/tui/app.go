package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

// App Bubbletea 应用
type App struct {
	program *tea.Program
	model   *Model
	logger  *slog.Logger
}

// NewApp 创建新的 Bubbletea 应用
func NewApp(logger *slog.Logger, cfg Config) *App {
	model := NewModel(logger, cfg)

	app := &App{
		model:  model,
		logger: logger.With("module", "tui-app"),
	}

	return app
}

// Run 运行应用
func (a *App) Run() error {
	a.logger.Debug("run bubbletea app start ...")
	defer a.logger.Debug("run bubbletea app end ...")

	// 创建程序
	a.program = tea.NewProgram(
		a.model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// 运行程序
	if _, err := a.program.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

// Quit 退出应用
func (a *App) Quit() {
	if a.program != nil {
		a.program.Quit()
	}
}
