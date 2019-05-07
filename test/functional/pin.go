package main

import (
	"fmt"

	"github.com/M-J-Hooper/rotaryphone"
	"github.com/brian-armstrong/gpio"
)

func main() {
	watcher := gpio.NewWatcher()
	watcher.AddPin(rotaryphone.DialActivePin)
	watcher.AddPin(rotaryphone.DialIncrementPin)
	watcher.AddPin(rotaryphone.LatchActivePin)

	for {
		p, v := watcher.Watch()
		fmt.Println("Got pin", p, v)
	}
}
