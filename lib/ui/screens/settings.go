package screens

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pchaussalet/lights-go/lib/input/midi"
	"github.com/pchaussalet/lights-go/lib/model"
	"github.com/pchaussalet/lights-go/lib/model/fixtures"
)

type SettingsScreen Screen

type SettingsSection int

const (
	MIDI SettingsSection = iota
	SACN
	FIXTURES
)

func Settings(app fyne.App, show *model.AppState) SettingsScreen {
	items := map[SettingsSection]*widget.AccordionItem{
		MIDI:     midiForm(app),
		SACN:     widget.NewAccordionItem("sACN", widget.NewForm()),
		FIXTURES: fixturesForm(show),
	}
	accordion := widget.NewAccordion()
	for _, tab := range items {
		accordion.Append(tab)
	}
	accordion.MultiOpen = false
	return SettingsScreen{
		Content: accordion,
	}
}

func (item *SettingsScreen) Open(section SettingsSection) {
	item.Content.(*widget.Accordion).Open(int(section))
}

func midiForm(app fyne.App) *widget.AccordionItem {
	inPort := widget.NewSelect(midi.ListInPorts(), func(s string) {
		app.Preferences().SetString("midi.in.port", s)
	})
	inPort.SetSelected(app.Preferences().String("midi.in.port"))
	return widget.NewAccordionItem("MIDI", widget.NewForm(
		widget.NewFormItem("Input Port", inPort),
	))
}

type FixtureFormState struct {
	selection      fixtures.Fixture
	selectionIndex int
}

func fixturesForm(show *model.AppState) *widget.AccordionItem {
	formState := FixtureFormState{}
	list := widget.NewList(
		func() int {
			return len(show.Fixtures)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(strconv.Itoa(show.Fixtures[lii].BaseAddress()))
		},
	)
	channelEntry := widget.NewEntry()
	channelEntry.OnChanged = func(s string) {
		address, _ := strconv.ParseInt(s, 10, 32)
		show.Fixtures[formState.selectionIndex] = fixtures.NewFixture(int(address))
		list.RefreshItem(formState.selectionIndex)
	}
	list.OnSelected = func(id widget.ListItemID) {
		fixture := show.Fixtures[id]
		formState.selection = *fixture
		formState.selectionIndex = id
		channelEntry.SetText(strconv.Itoa(fixture.BaseAddress()))
	}
	return widget.NewAccordionItem("Fixtures", container.New(layout.NewGridLayoutWithColumns(2),
		container.New(layout.NewVBoxLayout(),
			list,
			widget.NewButton("Add", func() {
				show.Fixtures = append(show.Fixtures, fixtures.NewFixture(0))
				list.Refresh()
			}),
		),
		widget.NewForm(
			widget.NewFormItem("Base Address", channelEntry),
		),
	))
}
