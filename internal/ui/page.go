package ui

import (
	"github.com/rivo/tview"
	"github.com/sagikazarmark/slog-shim"
)

// Pages 页面
type Pages struct {
	*tview.Pages

	logger *slog.Logger
}

// NewPages 新建
func NewPages(logger *slog.Logger) *Pages {
	p := Pages{
		Pages:  tview.NewPages(),
		logger: logger.With("module", "ui-pages"),
	}

	return &p
}

// // GetFrontPage 获取当前页面
// func (p *Pages) GetFrontPage() (string, tview.Primitive) {
// 	return p.Pages.GetFrontPage()
// }
