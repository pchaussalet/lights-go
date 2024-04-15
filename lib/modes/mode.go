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
	Mode       AppMode
	Title      string
	Color      tcell.Color
	KeyHandler func(commandLine *tview.TextView) func(event *tcell.EventKey) *tcell.EventKey
	Reset      func()
	Exit       func()
	Content    func() *tview.Flex
	Refresh    func()
}
