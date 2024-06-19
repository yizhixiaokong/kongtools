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

	a.Menu().AddItem("Welcome", "Welcome page", rune('w'), func() {
		a.logger.Debug("switch to welcome page ...")
		a.Content.SwitchToPage("welcome")
	})

	// a.TestSwitchPagesAndContent() // test switch pages and content logic

	a.Menu().AddItem("Quit", "Press to exit", rune('q'), func() {
		a.logger.Debug("quit app ...")
		a.Application.Stop()
	})

	a.flexLayout()

	return nil
}

// TestSwitchPagesAndContent 用来测试页面和内容切换的方法
func (a *App) TestSwitchPagesAndContent() {
	// 用来测试page切换的,可删
	a.Menu().AddItem("Add Page", "Add a new page", rune('a'), func() {
		a.logger.Debug("add page ...")
		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetText("This is a new page.")
		textView.SetBorder(true)
		backButton := tview.NewButton("Back").SetSelectedFunc(func() {
			a.Main.SwitchToPage("main")
		})
		a.Main.AddPage("page", tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(textView, 0, 1, false).AddItem(backButton, 1, 1, true),
			true, false)

		a.Menu().AddItem("Switch Page", "Switch to new page", 0, func() {
			a.Main.SwitchToPage("page")
		})
	})

	// 用来测试Main的Content切换,可删
	a.Menu().AddItem("Add Content", "Add a new content", rune('c'), func() {
		a.logger.Debug("add content ...")
		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetText("This is a new content.")
		textView.SetBorder(true)
		a.Content.AddPage("content2", textView,
			true, false)
		a.Menu().AddItem("Switch Content", "Switch to new content", 0, func() {
			a.Content.SwitchToPage("content2")
		})
	})
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
