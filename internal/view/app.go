package view

import (
	"kongtools/internal/ui"
	"log/slog"

	"github.com/rivo/tview"
)

// App 应用视图
type App struct {
	*ui.App
	Content *ui.Pages

	logger *slog.Logger
}

// NewApp 新建
func NewApp(logger *slog.Logger) *App {
	a := App{
		App:     ui.NewApp(logger),
		Content: ui.NewPages(logger),
		logger:  logger.With("module", "view-app"),
	}

	a.Views()["welcome"] = NewWelcome(logger)

	return &a
}

// Init 初始化
func (a *App) Init() error {
	a.logger.Debug("init app start ...")
	defer a.logger.Debug("init app end ...")

	a.App.Init()

	// a.Menu().AddItem("Test1", "Press to test1", rune('a'+0), nil) //! test
	// a.Menu().AddItem("Test2", "Press to test2", rune('a'+1), nil) //! test
	a.Menu().AddItem("Quit", "Press to exit", rune('q'), func() {
		a.logger.Debug("quit app ...")
		a.Application.Stop()
	})
	// a.Menu().SetSelectedFunc(func(i int, mainText, secondaryText string, r rune) {
	// 	a.Welcome().SetTitle(mainText) //! test
	// })

	a.flexLayout()

	return nil
}

// Run 运行
func (a *App) Run() error {
	a.logger.Debug("run app start ...")
	defer a.logger.Debug("run app end ...")

	a.Content.AddPage("welcome", a.Welcome(), true, true)
	a.Main.SwitchToPage("main")

	return a.Application.Run()
}

// flexLayout app flex布局
func (a *App) flexLayout() {
	main := tview.NewFlex().SetDirection(tview.FlexColumn)

	main.AddItem(a.Menu(), 0, 1, true)  // SideBar
	main.AddItem(a.Content, 0, 3, true) // Body

	a.Main.AddPage("main", main, true, true)
}

// buildContent 内容
func (a *App) buildContent() tview.Primitive {
	content := tview.NewFlex()
	content.SetDirection(tview.FlexColumn)

	content.AddItem(a.Menu(), 0, 1, true)  // SideBar
	content.AddItem(a.Content, 0, 3, true) // Body

	return content
}

// Welcome 欢迎页
func (a *App) Welcome() *Welcome {
	return a.Views()["welcome"].(*Welcome)
}
