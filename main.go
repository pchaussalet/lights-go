package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/pchaussalet/lights-go/lib/modes"
	"github.com/rivo/tview"
)

type KeyHandler func(event *tcell.EventKey) *tcell.EventKey

type AppState struct {
	main        *tview.Flex
	header      *tview.TextView
	commandLine *tview.TextView
	keyHandler  KeyHandler
	modes       []*modes.Mode
}

func main() {
	logFile, _ := os.Create("/tmp/lights_go.log")
	log.Default().SetOutput(logFile)

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("")
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

	state := NewState(main, header, commandLine, modes.LoadLive())

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			if commandLine.GetText(true) != "" {
				commandLine.SetText("")
				if state.currentMode().Reset != nil {
					state.currentMode().Reset()
				}
			} else if len(state.modes) > 1 {
				state.exitMode()
			}
		default:
			switch event.Rune() {
			case 338:
				app.Stop()
			case ';':
				state.enterMode(modes.LoadPatch(), true)
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

func NewState(main *tview.Flex, header, commandLine *tview.TextView, defaultMode *modes.Mode) *AppState {
	state := AppState{
		main:        main,
		header:      header,
		commandLine: commandLine,
		modes:       []*modes.Mode{},
	}
	state.enterMode(defaultMode, true)
	return &state
}

func (state *AppState) currentMode() *modes.Mode {
	return state.modes[len(state.modes)-1]
}

func (state *AppState) exitMode() {
	state.currentMode().Exit()
	state.main.Clear()
	state.modes = state.modes[:len(state.modes)-1]
	state.enterMode(state.currentMode(), false)
}

func (state *AppState) enterMode(mode *modes.Mode, addToStack bool) {
	if addToStack {
		state.modes = append(state.modes, mode)
	}
	state.header.SetText(mode.Title)
	if mode.Content != nil && mode.Content() != nil {
		state.main.Clear().AddItem(mode.Content(), 0, 1, false)
	}
	state.main.SetBorderColor(mode.Color)
	if mode.KeyHandler != nil {
		state.keyHandler = mode.KeyHandler(state.commandLine)
	} else {
		state.keyHandler = nil
	}
	if mode.Refresh != nil {
		mode.Refresh()
	}
}
