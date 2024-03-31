package midi

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiListener struct {
	onControlChange func(channel uint8, control uint8, value uint8)
	onNoteStart     func(channel uint8, note uint8, velocity uint8)
	onNoteEnd       func(channel uint8, note uint8)
}

func ListPorts() {
	fmt.Printf("MIDI IN Ports\n")
	fmt.Println(midi.GetInPorts())
	fmt.Printf("\n\nMIDI OUT Ports\n")
	fmt.Println(midi.GetOutPorts())
	fmt.Printf("\n\n")
}

func NewMidiListener(portName string) *MidiListener {
	item := MidiListener{nil, nil, nil}

	in, err := midi.FindInPort(portName)
	if err != nil {
		log.Fatalf("Cannot find MIDI In port matching %s\n", portName)
	}
	log.Printf("Opened MIDI In port %s\n", in.String())

	_, err = midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch uint8
		var key, vel uint8
		var cc, val uint8
		switch {
		case msg.GetControlChange(&ch, &cc, &val):
			if item.onControlChange != nil {
				item.onControlChange(ch, cc, val)
			}
		case msg.GetNoteStart(&ch, &key, &vel):
			if item.onNoteStart != nil {
				item.onNoteStart(ch, key, vel)
			}
		case msg.GetNoteEnd(&ch, &key):
			if item.onNoteEnd != nil {
				item.onNoteEnd(ch, key)
			}
		default:
			// ignore
		}
	})
	if err != nil {
		log.Fatalf("Cannot listen to MIDI events on In port %s (%v)\n", in.String(), err)
	}
	log.Printf("Listening on MIDI In port %s\n", in.String())

	return &item
}

func (item *MidiListener) Close() {
	midi.CloseDriver()
}

func (item *MidiListener) OnControlChange(handler func(uint8, uint8, uint8)) {
	item.onControlChange = handler
}

func (item *MidiListener) OnNoteStart(handler func(uint8, uint8, uint8)) {
	item.onNoteStart = handler
}

func (item *MidiListener) OnNoteEnd(handler func(uint8, uint8)) {
	item.onNoteEnd = handler
}
