package modes

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/pchaussalet/lights-go/lib/model"

	"github.com/rivo/tview"
)

type PatchMode int

const (
	Read PatchMode = iota
	Add
	Edit
	Delete
)

type localState struct {
	mode  PatchMode
	patch *model.Patch
}

func LoadPatch() *Mode {
	state := localState{
		mode:  Read,
		patch: model.GetPatch(),
	}

	mode := &Mode{
		Mode:            Patch,
		Title:           "Patch",
		BackgroundColor: tcell.NewRGBColor(10, 0, 0),
		Content:         nil,
	}

	refresh(mode, state.patch)

	mode.KeyHandler = func(commandLine *tview.TextView) func(event *tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			cli := commandLine.GetText(true)
			switch state.mode {
			case Read:
				switch event.Key() {
				default:
					switch event.Rune() {
					case 'A', 'a':
						state.mode = Add
						commandLine.SetText("add ").SetBackgroundColor(tcell.NewRGBColor(30, 0, 0))
					}
				}
			case Add:
				switch event.Key() {
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					commandLine.SetText(cli[0 : len(cli)-1])
					if len(cli) == 1 {
						state.mode = Read
						commandLine.SetBackgroundColor(tcell.NewRGBColor(0, 0, 0))
					}
				case tcell.KeyEnter:
					// address, _ := strconv.ParseInt(cli[strings.LastIndex(cli, " "):], 10, 32)
					state.patch, _ = state.patch.AddDimmer(int(13))
					commandLine.Clear().SetBackgroundColor(tcell.NewRGBColor(0, 0, 0))
					refresh(mode, state.patch)
					// commandLine.SetText(fmt.Sprintf("%v %v %v %v %v", bp, bf, len(state.patch.Fixtures), fixtureCells.GetRowCount(), address))
				default:
					switch event.Rune() {
					case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
						commandLine.SetText(cli + string(event.Rune()))
					}
				}
			}
			return event

		}
	}

	return mode
}

func refresh(mode *Mode, patch *model.Patch) {
	fixtureCells := tview.NewTable().
		SetBorders(true)

	for _, fixture := range patch.Fixtures {
		fixtureCells.InsertRow(0)
		fixtureCells.SetCellSimple(0, 0, strconv.Itoa(fixture.BaseAddress()))
	}

	mode.Content = fixtureCells.Box
}
