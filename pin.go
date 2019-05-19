package rotaryphone

import (
	"time"

	"github.com/M-J-Hooper/debounce"
	"github.com/brian-armstrong/gpio"
)

func debouncedPin(pin uint, t time.Duration) chan gpio.WatcherNotification {
	c := make(chan interface{})
	p := make(chan gpio.WatcherNotification)

	w := gpio.NewWatcher()
	w.AddPin(pin)
	defer w.Close()

	go convertFromPinChannel(w.Notification, c)
	db := debounce.Channel(c, t)
	go convertToPinChannel(db, p)
	return p
}

func convertFromPinChannel(from chan gpio.WatcherNotification, to chan interface{}) {
	for n := range from {
		to <- interface{}(n)
	}
}

func convertToPinChannel(from chan interface{}, to chan gpio.WatcherNotification) {
	for n := range from {
		to <- n.(gpio.WatcherNotification)
	}
}
