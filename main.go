package main

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"math/big"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/Hundemeier/go-sacn/sacn"

	"github.com/pchaussalet/lights-go/lib/midi"
	"github.com/pchaussalet/lights-go/lib/ui/widgets"
)

func main() {
	// var err error
	universe := []int{}
	for i := 0; i < 512; i++ {
		universe = append(universe, 0)
	}
	universeB := binding.BindIntList(&universe)

	// updateRandomValues(universeB)

	ch, err := loadSACN()
	if err != nil {
		return
	}
	defer close(ch)

	var midiPort string
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "list":
			midi.ListPorts()
			return
		default:
			midiPort = os.Args[1]
		}
	}

	midiListener := midi.NewMidiListener(midiPort)
	defer midiListener.Close()

	midiListener.OnControlChange(func(channel, control, value uint8) {
		dmxVal := int((float32(value) / 127) * 255)
		universeB.SetValue(int(control-1), dmxVal)
	})

	a := app.New()
	w := a.NewWindow("Lights, GO!")
	w.Resize(fyne.NewSize(1024, 600))

	tabs := container.NewAppTabs(
		container.NewTabItem("Universe", widget.NewGridWrap(
			func() int {
				return universeB.Length()
			},
			func() fyne.CanvasObject {
				return widgets.NewChannelCell(color.Gray{0})
			},
			func(gwii widget.GridWrapItemID, co fyne.CanvasObject) {
				cell := co.(*widgets.ChannelCell)
				item, _ := universeB.GetItem(gwii)
				cell.SetChannel(gwii)
				cell.Bind(item.(binding.Int))
			},
		)),
	)
	w.SetContent(tabs)
	w.ShowAndRun()
}

func updateRandomValues(universeB binding.IntList) {
	for i := 0; i < universeB.Length(); i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(255))
		universeB.SetValue(i, int(v.Int64()))
	}
	time.AfterFunc(time.Duration(2_000_000_000), func() {
		updateRandomValues(universeB)
	})
}

func loadSACN() (chan<- []byte, error) {
	trans, err := sacn.NewTransmitter("", [16]byte{12, 3, 17}, "lights-go")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}
	trans.SetMulticast(1, true)
	ch, err := trans.Activate(1)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	return ch, nil
}
