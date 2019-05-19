package rotaryphone

import (
	"fmt"
	"time"
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
	incChan := debouncedPin(dialIncrementPin, 10*time.Millisecond)
	activeChan := debouncedPin(dialActivePin, 10*time.Millisecond)

	var active bool
	var count int
	for {
		select {
		case n := <-activeChan:
			if n.Value == 1 {
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
		case n := <-incChan:
			if n.Value == 0 && active {
				count++
			}
		}
	}
}
