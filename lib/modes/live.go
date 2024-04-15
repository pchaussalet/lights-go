package modes

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/pchaussalet/lights-go/lib/model"
	"github.com/rivo/tview"
)

func LoadLive() *Mode {
	mode := &Mode{
		Mode:  Live,
		Title: "Live",
		Color: tcell.NewRGBColor(10, 10, 10),
	}

	patch := model.GetPatch()

	content := tview.NewGrid()

	mode.Refresh = func() {
		content.AddItem(tview.NewTextView().SetText("Foo "+strconv.Itoa(len(patch.Fixtures))), 0, 0, 1, 1, 0, 0, false)

		for _, fixture := range patch.Fixtures {
			cell := tview.NewFlex().
				AddItem(tview.NewTextView().SetText(fixture.Name()), 0, 1, false)
			content.AddItem(cell, 0, 0, 1, 1, 0, 0, false)
		}
	}

	mode.Exit = func() {
	}

	mode.Content = func() *tview.Flex {
		return tview.NewFlex().AddItem(content, 0, 1, false)
	}

	mode.Refresh()

	return mode
}
