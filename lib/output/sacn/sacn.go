package sacn

import (
	"fmt"

	sacn_lib "github.com/Hundemeier/go-sacn/sacn"
)

func NewSACNTransmitter() (chan<- []byte, error) {
	trans, err := sacn_lib.NewTransmitter("", [16]byte{12, 3, 17}, "lights-go")
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
