package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pchaussalet/lights-go/lib/modes"
	"github.com/rivo/tview"
)

type AppState struct {
	mode        modes.AppMode
	main        *tview.Flex
	header      *tview.TextView
	commandLine *tview.TextView
	keyHandler  func(event *tcell.EventKey) *tcell.EventKey
}

func main() {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Header")
	header.SetBorderPadding(1, 1, 1, 1)

	main := tview.NewFlex()

	commandLine := tview.NewTextView().
		SetTextAlign(tview.AlignLeft)

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(main, 1, 0, 1, 1, 1, 0, false).
		AddItem(commandLine, 2, 0, 1, 1, 0, 0, true)

	state := AppState{
		main:        main,
		header:      header,
		commandLine: commandLine,
	}

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			if commandLine.GetText(true) != "" {
				commandLine.SetText("")
			}
		default:
			switch event.Rune() {
			case 338:
				app.Stop()
			case ';':
				state.enterMode(modes.LoadPatch())
			// case 'Q', 'q':
			// 	commandLine.SetText("cue")
			// case 'R', 'r':
			// 	commandLine.SetText("record")
			// case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', '#':
			// 	commandLine.SetText(commandLine.GetText(true) + string(event.Rune()))
			default:
				if state.keyHandler != nil {
					state.keyHandler(event)
				}
			}
		}
		return event
	})

	if err := app.
		SetRoot(grid, true).
		SetFocus(grid).
		Run(); err != nil {
		panic(err)
	}
}

func (item *AppState) enterMode(mode *modes.Mode) {
	item.mode = mode.Mode
	item.header.SetText(mode.Title)
	item.main.Clear().AddItem(mode.Content, 0, 1, false)
	item.main.SetBackgroundColor(mode.BackgroundColor)
	item.keyHandler = mode.KeyHandler(item.commandLine)
}
