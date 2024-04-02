package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/pchaussalet/lights-go/lib/input/midi"
	"github.com/pchaussalet/lights-go/lib/model"
	"github.com/pchaussalet/lights-go/lib/output/sacn"
	"github.com/pchaussalet/lights-go/lib/ui/screens"
)

func main() {
	outputs := map[string]chan<- []byte{}

	universeB := binding.NewIntList()
	for i := 0; i < 512; i++ {
		universeB.Append(0)
	}

	ch, err := sacn.NewSACNTransmitter()
	if err != nil {
		return
	}
	defer close(ch)
	outputs["sacn"] = ch

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu(
			"Settings",
			fyne.NewMenuItem("MIDI", func() {}),
		),
	)

	a := app.NewWithID("con.github.pchaussalet.lights-go")
	appTheme := a.Settings().Theme()
	prefs := a.Preferences()
	lastFile := prefs.String("lastFile")
	w := a.NewWindow("Lights, GO!")
	w.Resize(fyne.NewSize(1024, 600))
	// w.SetFullScreen(true)
	w.SetMainMenu(mainMenu)

	var tabs *container.AppTabs
	appState, err := loadLastFile(lastFile, a, tabs)
	if err != nil {
		warningDialog := dialog.NewError(err, w)
		warningDialog.Show()
		prefs.RemoveValue("lastFile")
		appState, _ = loadLastFile("", a, tabs)
	}
	defer appState.Close()

	screens := []screens.Screen{
		screens.NewFixtures(appState),
		screens.Universe(appState.Fixtures),
		screens.Screen(screens.Settings(a, appState)),
	}

	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Fixtures", appTheme.Icon(theme.IconNameHome), screens[0].Content),
		container.NewTabItemWithIcon("Universe", appTheme.Icon(theme.IconNameSearch), screens[1].Content),
		container.NewTabItemWithIcon("Settings", appTheme.Icon(theme.IconNameSettings), screens[2].Content),
	)

	content := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			widget.NewButtonWithIcon("", appTheme.Icon(theme.IconNameDocumentSave), func() {
				saveAs := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
					saveShowToWriter(appState, uc)
					a.Preferences().SetString("lastFile", uc.URI().Path())
				}, w)
				if lastFile == "" {
					saveAs.Show()
				} else {
					err := saveShow(appState, lastFile)
					if err != nil {
						dialog.NewConfirm("Unable to save", fmt.Sprintf("%s\nDo you want to save to another file?", err.Error()), func(save bool) {
							if save {
								saveAs.Show()
							}
						}, w)
					}
				}
			}),
		),
		tabs,
	)

	w.SetContent(content)
	w.Show()

	a.Run()
}

func loadLastFile(lastFile string, a fyne.App, c *container.AppTabs) (*model.AppState, error) {
	show := &model.AppState{}
	if lastFile != "" {
		showJson, err := os.ReadFile(lastFile)
		if err != nil {
			return show, fmt.Errorf("cannot open last open file %s\n%v", lastFile, err)
		}
		show, err = model.LoadShow(showJson)
		if err != nil {
			return show, fmt.Errorf("error occurred when reading last open file %s\n%v", lastFile, err)
		}
	}
	midiListener := loadMidi(a.Preferences(), show)
	if midiListener == nil {
		displaySettings(a, c, "midi")
		return show, nil
	}
	show.Add(midiListener)
	return show, nil
}

func saveShow(show *model.AppState, lastFile string) error {
	if lastFile != "" {
		showJson, err := show.ToJSON()
		if err != nil {
			return fmt.Errorf("error occurred while preparing show for save (%v)", err)
		}
		err = os.WriteFile(lastFile, showJson, 0644)
		if err != nil {
			return fmt.Errorf("unable to save show to file %s (%v)", lastFile, err)
		}
	}
	return nil
}

func saveShowToWriter(show *model.AppState, writer fyne.URIWriteCloser) error {
	showJson, err := show.ToJSON()
	if err != nil {
		return fmt.Errorf("error occurred while preparing show for save (%v)", err)
	}
	_, err = writer.Write(showJson)
	if err != nil {
		return fmt.Errorf("unable to save show to file %s (%v)", writer.URI(), err)
	}
	return nil
}

func displaySettings(a fyne.App, c *container.AppTabs, tab string) {
	for _, tab := range c.Items {
		if tab.Text != "Settings" {
			c.DisableItem(tab)
		} else {
			c.Select(tab)
			tab.Content.(*widget.Accordion).CloseAll()
			tab.Content.(*widget.Accordion).Open(0)
		}
	}
}

func loadMidi(prefs fyne.Preferences, appState *model.AppState) *midi.MidiListener {
	midiPort := prefs.String("midi.in.port")
	var midiListener *midi.MidiListener
	if midiPort != "" {
		midiListener = midi.NewMidiListener(midiPort)

		midiListener.OnControlChange(func(channel, control, value uint8) {
			dmxVal := byte((float32(value) / 127) * 255)
			fixtureChannel, ok := appState.Mappings[channel][control]
			if ok && fixtureChannel != nil {
				fixtureChannel.Set(dmxVal)
			}
		})
	}
	return midiListener
}
