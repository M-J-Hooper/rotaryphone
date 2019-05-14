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
	pinWatcher := gpio.NewWatcher()
	pinWatcher.AddPin(dialIncrementPin)
	pinWatcher.AddPin(dialActivePin)
	defer pinWatcher.Close()

	watcher := NewDebouncedWatcher(pinWatcher, 10*time.Millisecond)

	var active bool
	var count int
	for {
		pin, value := watcher.Watch()
		fmt.Println("Dial got stable signal", pin, value)
		if pin == dialActivePin {
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
		} else if pin == dialIncrementPin && value == 0 {
			if active {
				count++
			}
		}
	}
}
