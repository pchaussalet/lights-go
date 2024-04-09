package modes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppMode int

const (
	Live AppMode = iota
	Patch
)

type Mode struct {
	Mode            AppMode
	Title           string
	BackgroundColor tcell.Color
	Content         *tview.Box
	KeyHandler      func(commandLine *tview.TextView) func(event *tcell.EventKey) *tcell.EventKey
}
