package main

import (
	"github.com/eiannone/keyboard"
)

type KeyboardEvent byte

const (
	KeyboardEventSaveCurrent = 's'
	KeyboardEventExit        = 'q'
)

func CaptureKeyboardEvents() chan KeyboardEvent {
	ch := make(chan KeyboardEvent)
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	go func() {
		defer close(ch)
		for {
			char, key, _ := keyboard.GetKey()
			if char == 'q' || char == 'Q' || key == 3 {
				ch <- KeyboardEventExit
				return
			} else if char == 's' || char == 'S' {
				ch <- KeyboardEventSaveCurrent
			}
		}
	}()
	return ch
}
