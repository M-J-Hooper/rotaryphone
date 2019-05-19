package rotaryphone

import (
	"fmt"
	"time"

	"github.com/brian-armstrong/gpio"
)

const dialIncrementPin = 20 //14 //=> orange wire
const dialActivePin = 21    //15    //=> brown wire
//3.3v  => white wire

type Dial struct {
	Digit chan int
}

func NewDial() *Dial {
	d := &Dial{make(chan int)}
	go d.run()
	return d
}

func (d Dial) run() {
	activeChan := make(chan interface{})
	activePinWatcher := gpio.NewWatcher()
	activePinWatcher.AddPin(dialActivePin)
	defer activePinWatcher.Close()
	go castGpioChannel(activePinWatcher.Notification, activeChan)
	activeWatcher := NewDebouncedWatcher(activeChan, 10*time.Millisecond)

	incChan := make(chan interface{})
	incPinWatcher := gpio.NewWatcher()
	incPinWatcher.AddPin(dialIncrementPin)
	defer incPinWatcher.Close()
	go castGpioChannel(incPinWatcher.Notification, incChan)
	incWatcher := NewDebouncedWatcher(incChan, 10*time.Millisecond)

	var active bool
	var count int
	for {
		select {
		case n := <-activeWatcher.Notification:
			value := n.(gpio.WatcherNotification).Value
			if value == 1 {
				active = true
			} else {
				active = false
				if count > 0 {
					if count > 9 {
						count = 0
					}
					fmt.Println("Sending digit", count)
					d.Digit <- count
					count = 0
				}
			}
		case n := <-incWatcher.Notification:
			value := n.(gpio.WatcherNotification).Value
			if value == 0 && active {
				count++
			}
		}
	}
}
