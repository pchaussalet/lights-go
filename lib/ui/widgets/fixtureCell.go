package widgets

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FixtureCell struct {
	widget.BaseWidget
	Index *canvas.Text
	Value *widget.Label
	Color *canvas.Rectangle
}

func NewFixtureCell() *FixtureCell {
	item := &FixtureCell{
		Index: canvas.NewText("", color.Gray{128}),
		Value: widget.NewLabel("0"),
		Color: canvas.NewRectangle(color.Black),
	}
	item.Color.SetMinSize(fyne.NewSize(item.Color.Size().Width, 10))
	item.ExtendBaseWidget(item)

	return item
}

func (item *FixtureCell) SetChannel(value int) {
	item.Index.Text = strconv.Itoa(value)
}

func (item *FixtureCell) BindColor(colorBytes binding.Bytes) {
	colorBytes.AddListener(binding.NewDataListener(func() {
		value, _ := colorBytes.Get()
		item.Color.FillColor = color.RGBA{value[0], value[1], value[2], 255}
		item.Color.Refresh()
	}))
}

func (item *FixtureCell) BindDimmer(dimmer binding.Int) {
	item.Value.Bind(binding.IntToString(dimmer))
}

func (item *FixtureCell) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(),
		item.Color,
		container.New(layout.NewGridLayoutWithColumns(2),
			container.New(layout.NewPaddedLayout(), item.Index),
			container.New(layout.NewPaddedLayout(), item.Value),
		),
	))
	return widget.NewSimpleRenderer(c)
}
