package modes

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/pchaussalet/lights-go/lib/model"
	"github.com/pchaussalet/lights-go/lib/model/fixtures"

	"github.com/rivo/tview"
)

type PatchMode int

const (
	Read PatchMode = iota
	Add
	Edit
	Delete
)

type patchState struct {
	mode    PatchMode
	command interface{}
	subMode int
	patch   *model.Patch
}

type AddPart int

const (
	None AddPart = iota
	Address
	Subs
	Name
)

func (part AddPart) value() int {
	return int(part)
}

type AddCommand struct {
	address int
	subs    []fixtures.ChannelRole
	name    string
}

func LoadPatch() *Mode {
	state := patchState{
		mode:  Read,
		patch: model.GetPatch(),
	}

	mode := &Mode{
		Mode:  Patch,
		Title: "Patch",
		Color: tcell.NewRGBColor(0, 0, 30),
	}

	content := tview.NewTable().SetBorders(true)

	refresh(content, state.patch)

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
						state.subMode = Address.value()
						state.command = &AddCommand{}
						commandLine.SetText("add ").SetBackgroundColor(tcell.NewRGBColor(30, 0, 0))
					}
				}
			case Add:
				command := state.command.(*AddCommand)
				switch state.subMode {
				case Name.value():
				case Subs.value():
				}
				switch event.Key() {
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					commandLine.SetText(cli[0 : len(cli)-1])
					if state.subMode == Name.value() {
						command.name = command.name[:len(command.name)-1]
					}
					if len(cli) == 1 {
						commandLine.SetBackgroundColor(tcell.NewRGBColor(0, 0, 0))
						mode.Reset()
					}
				case tcell.KeyEnter:
					if command.subs == nil || len(command.subs) == 0 {
						command.subs = []fixtures.ChannelRole{fixtures.Dimmer}
					}
					state.patch.AddFixture(command.address, command.name, command.subs...)
					refresh(content, state.patch)
					commandLine.Clear().SetBackgroundColor(tcell.NewRGBColor(0, 0, 0))
					state.mode = Read
				default:
					switch state.subMode {
					case Name.value():
						if event.Rune() != ' ' {
							command.name += string(event.Rune())
							commandLine.SetText(cli + string(event.Rune()))
						}
					default:
						switch event.Rune() {
						case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
							commandLine.SetText(cli + string(event.Rune()))
						case 's':
							commandLine.SetText(cli + "subs ")
							command.subs = []fixtures.ChannelRole{}
							state.subMode = Subs.value()
						case 'n':
							commandLine.SetText(cli + "name ")
							command.name = ""
							state.subMode = Name.value()
						case 'd':
							command.subs = append(command.subs, fixtures.Dimmer)
							commandLine.SetText(cli + "Dimmer,")
						case 'r':
							command.subs = append(command.subs, fixtures.Red)
							commandLine.SetText(cli + "Red,")
						case 'g':
							command.subs = append(command.subs, fixtures.Green)
							commandLine.SetText(cli + "Green,")
						case 'b':
							command.subs = append(command.subs, fixtures.Blue)
							commandLine.SetText(cli + "Blue,")
						case 'w':
							command.subs = append(command.subs, fixtures.White)
							commandLine.SetText(cli + "White,")
						case 'c':
							command.subs = append(command.subs, fixtures.Color)
							commandLine.SetText(cli + "Color,")
						case 'o':
							command.subs = append(command.subs, fixtures.Gobo)
							commandLine.SetText(cli + "Gobo,")
						case ' ':
							switch state.subMode {
							case Address.value():
								cliParts := strings.Split(cli, " ")
								address, _ := strconv.ParseInt(cliParts[len(cliParts)-1], 10, 32)
								command.address = int(address)
								state.subMode = None.value()
								commandLine.SetText(strings.Join(cliParts[:len(cliParts)-1], " ") + " " + strconv.Itoa(command.address) + " ")
							case Subs.value():
								state.subMode = None.value()
								commandLine.SetText(cli[:len(cli)-1] + " ")
							}
						}
					}
				}
			}
			return event

		}
	}

	mode.Reset = func() {
		state.mode = Read
	}

	mode.Exit = func() {
	}

	mode.Content = func() *tview.Flex {
		return tview.NewFlex().AddItem(content, 0, 1, false)
	}

	return mode
}

func refresh(content *tview.Table, patch *model.Patch) {
	content.Clear()
	for _, fixture := range patch.Fixtures {
		content.
			InsertRow(0).
			SetCellSimple(0, 0, strconv.Itoa(fixture.BaseAddress())).
			SetCellSimple(0, 1, fixture.Name())
	}
}
