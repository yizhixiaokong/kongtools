package ui

import (
	"github.com/rivo/tview"
	"github.com/sagikazarmark/slog-shim"
)

// App 应用界面
type App struct {
	*tview.Application

	Main  *Pages
	views map[string]tview.Primitive

	logger *slog.Logger
}

// NewApp 新建
func NewApp(logger *slog.Logger) *App {
	a := App{
		Application: tview.NewApplication(),
		Main:        NewPages(logger),
		logger:      logger.With("module", "ui-app"),
	}

	a.views = map[string]tview.Primitive{
		"menu": NewMenu(logger),
	}
	return &a
}

// Init 初始化
func (a *App) Init() {
	a.logger.Debug("init app start ...")
	defer a.logger.Debug("init app end ...")
	a.setupApp()

	a.SetRoot(a.Main, true).EnableMouse(true)
}

func (a *App) setupApp() {
	a.bindKeys()
	a.setupStyles()
}

func (a *App) bindKeys() {
}

func (a *App) setupStyles() {
}

// Views Views
func (a *App) Views() map[string]tview.Primitive {
	return a.views
}

// Menu 菜单
func (a *App) Menu() *Menu {
	return a.views["menu"].(*Menu)
}
