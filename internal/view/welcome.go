package view

import "github.com/rivo/tview"

// WelcomeMsg 欢迎页
var WelcomeMsg = []string{
	` _    _  _____  _      _____  _____ ___  ___ _____  _ `,
	`| |  | ||  ___|| |    /  __ \|  _  ||  \/  ||  ___|| |`,
	`| |  | || |__  | |    | /  \/| | | || .  . || |__  | |`,
	`| |/\| ||  __| | |    | |    | | | || |\/| ||  __| | |`,
	`\  /\  /| |___ | |____| \__/\\ \_/ /| |  | || |___ |_|`,
	` \/  \/ \____/ \_____/ \____/ \___/ \_|  |_/\____/ (_)`,
}

// Welcome 欢迎页
type Welcome struct {
	*tview.TextView
}

// NewWelcome 新建
func NewWelcome() *Welcome {
	w := Welcome{
		TextView: tview.NewTextView(),
	}

	w.SetBorder(true)
	w.SetTextAlign(tview.AlignCenter)
	w.SetDynamicColors(true)
	w.SetWordWrap(true)
	w.SetWrap(false)
	w.SetTitle("Welcome")
	w.refreshWelcome()

	return &w
}

func (w *Welcome) refreshWelcome() {
	w.Clear()

	for _, line := range WelcomeMsg {
		w.Write([]byte(line + "\n"))
	}
}
