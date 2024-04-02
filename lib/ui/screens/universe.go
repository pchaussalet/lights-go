package screens

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/pchaussalet/lights-go/lib/model/fixtures"
	"github.com/pchaussalet/lights-go/lib/ui/widgets"
)

func Universe(fixturesList []*fixtures.Fixture) Screen {
	return Screen{Content: widget.NewGridWrap(
		func() int {
			return 512
		},
		func() fyne.CanvasObject {
			return widgets.NewChannelCell(color.Gray{0})
		},
		func(gwii widget.GridWrapItemID, co fyne.CanvasObject) {
			cell := co.(*widgets.ChannelCell)
			address := gwii + 1
			cell.SetChannel(address)
			channel := fixtures.FixtureChannelByAddress(fixturesList, address)
			if channel != nil {
				cell.Bind(channel.Value())
			}
		},
	)}
}
