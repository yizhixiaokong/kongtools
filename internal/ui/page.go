package ui

import "github.com/rivo/tview"

// Pages 页面
type Pages struct {
	*tview.Pages
}

// NewPages 新建
func NewPages() *Pages {
	p := Pages{
		Pages: tview.NewPages(),
	}

	return &p
}

// // GetFrontPage 获取当前页面
// func (p *Pages) GetFrontPage() (string, tview.Primitive) {
// 	return p.Pages.GetFrontPage()
// }
