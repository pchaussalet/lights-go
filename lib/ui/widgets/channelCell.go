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

type ChannelCell struct {
	widget.BaseWidget
	Index *canvas.Text
	Value *widget.Label
	Color *canvas.Rectangle
}

func NewChannelCell(intensity color.Gray) *ChannelCell {
	item := &ChannelCell{
		Index: canvas.NewText("", color.Gray{128}),
		Value: widget.NewLabel("0"),
		Color: canvas.NewRectangle(intensity),
	}
	item.Color.Resize(fyne.NewSize(item.Color.Size().Width, 10))
	item.ExtendBaseWidget(item)

	return item
}

func (item *ChannelCell) SetChannel(value int) {
	item.Index.Text = strconv.Itoa(value)
}

func (item *ChannelCell) Bind(value binding.Int) {
	item.Value.Bind(binding.IntToString(value))
	value.AddListener(binding.NewDataListener(func() {
		v, _ := value.Get()
		item.Color.FillColor = color.Gray{uint8(v)}
	}))
}

func (item *ChannelCell) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(layout.NewPaddedLayout(), container.New(layout.NewVBoxLayout(),
		item.Color,
		container.New(layout.NewGridLayoutWithColumns(2),
			container.New(layout.NewPaddedLayout(), item.Index),
			container.New(layout.NewPaddedLayout(), item.Value),
		),
	))
	// c := container.New(layout.NewVBoxLayout(),
	// 	item.Color,
	// 	container.New(layout.NewGridLayoutWithColumns(2),
	// 		item.Index,
	// 		item.Value,
	// 	),
	// )
	return widget.NewSimpleRenderer(c)
}
