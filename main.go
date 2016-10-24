package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AllenDang/w32"
	"github.com/simulatedsimian/joystick"
)

const (
	xbox360StartButton = 128
	keyboardRButton    = 82
	keyEventDown       = 0
	keyEventUp         = 0x0002
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	js, err := joystick.Open(0)
	exitOnError(err)

	ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

	for {
		<-ticker.C
		state, err := js.Read()
		exitOnError(err)
		if state.Buttons == xbox360StartButton {
			pressKey(keyboardRButton, keyEventDown)
		} else {
			pressKey(keyboardRButton, keyEventUp)
		}
	}
}

func pressKey(vk uint16, flags int) {
	w32.SendInput([]w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki: w32.KEYBDINPUT{
				WVk:         vk,
				WScan:       0,
				DwFlags:     uint32(flags),
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	})
}
