package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/pchaussalet/lights-go/lib/model"
	"github.com/pchaussalet/lights-go/lib/ui/widgets"
)

func NewFixtures(state *model.AppState) Screen {
	fixturesList := state.Fixtures

	return Screen{Content: widget.NewGridWrap(
		func() int {
			return len(fixturesList)
		},
		func() fyne.CanvasObject {
			return widgets.NewFixtureCell()
		},
		func(gwii widget.GridWrapItemID, co fyne.CanvasObject) {
			cell := co.(*widgets.FixtureCell)
			fixture := fixturesList[gwii]
			cell.SetChannel(gwii + 1)
			cell.BindColor(fixture.Color())
			cell.BindDimmer(fixture.Dimmer())
		},
	)}
}
