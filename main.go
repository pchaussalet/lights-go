package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"

	"github.com/Hundemeier/go-sacn/sacn"
)

func main() {
	var err error

	trans, err := sacn.NewTransmitter("", [16]byte{12, 3, 17}, "lights-go")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	trans.SetMulticast(1, true)
	ch, err := trans.Activate(1)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	defer close(ch)

	defer midi.CloseDriver()
	universe := []binding.Int{}
	for i := 0; i < 512; i++ {
		addressValue := binding.NewInt()
		universe = append(universe, addressValue)
	}
	for i := 0; i < 512; i++ {
		universe[i].AddListener(binding.NewDataListener(func() {
			status := []byte{}
			for i := 0; i < 512; i++ {
				newValue, _ := universe[i].Get()
				status = append(status, byte(newValue))
			}
			println(status)
			ch <- status
		}))
	}

	var in drivers.In

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "list":
			fmt.Printf("MIDI IN Ports\n")
			fmt.Println(midi.GetInPorts())
			fmt.Printf("\n\nMIDI OUT Ports\n")
			fmt.Println(midi.GetOutPorts())
			fmt.Printf("\n\n")
			return
		default:
			in, err = midi.FindInPort(os.Args[1])
			if err != nil {
				fmt.Printf("can't find %v\n", os.Args[1])
				return
			}
		}
	}

	_, err = midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch uint8
		var key, vel uint8
		var cc, val uint8
		switch {
		case msg.GetControlChange(&ch, &cc, &val):
			newVal := (float32(val) / 127) * 255
			universe[cc].Set(int(newVal))
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
		default:
			// ignore
		}
	})

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	a := app.New()
	w := a.NewWindow("Lights GO")
	w.Resize(fyne.NewSize(800, 600))

	status := container.NewGridWrap(fyne.NewSize(75, 50), channel(universe, 0))
	for i := 1; i < len(universe); i++ {
		status.Add(channel(universe, i))
	}

	content := container.NewGridWithColumns(2, status)
	content.Resize(fyne.NewSize(800, 600))
	w.SetContent(content)
	w.ShowAndRun()
}

func channel(universe []binding.Int, index int) *fyne.Container {
	border := color.Gray{128}
	colorRectangle := canvas.NewRectangle(color.Gray{0})
	colorRectangle.StrokeColor = border
	colorRectangle.StrokeWidth = 1
	colorRectangle.Resize(fyne.NewSize(colorRectangle.Size().Width, 2))
	universe[index].AddListener(binding.NewDataListener(func() {
		value, _ := universe[index].Get()
		colorRectangle.FillColor = color.Gray{uint8(value)}
	}))
	valueLabel := widget.NewLabelWithData(binding.IntToString(universe[index]))

	return container.NewGridWithRows(3, colorRectangle, canvas.NewText(strconv.Itoa(index), color.Gray{192}), valueLabel)
}
