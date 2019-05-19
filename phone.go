package rotaryphone

import (
	"strconv"
	"time"
)

type OutputAdapter interface {
	Call(string)
	Hangup()
}

type Rotary struct {
	OutputAdapter
	digitTimeout time.Duration
	dial         Dial
	latch        Latch
}

func NewRotary(digitTimeout time.Duration) *Rotary {
	return &Rotary{
		NewOfonoAdapter(),
		digitTimeout,
		*NewDial(),
		*NewLatch(),
	}
}

func (r Rotary) Run() {
	var dialing bool
	var number string
	lastDigit := time.Now()
	for {
		select {
		case digit := <-r.dial.Digit:
			if dialing {
				number += strconv.Itoa(digit)
				lastDigit = time.Now()
				println("New number is " + number)
			}
		case dialing = <-r.latch.Active:
			if !dialing {
				println("Handset on the latch")
				r.Hangup()
				number = ""
			} else {
				println("Handset off the latch")
			}
		default:
			if time.Since(lastDigit) > r.digitTimeout && len(number) > 0 {
				println("Calling " + number)
				r.Call(number)
				lastDigit = time.Now()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
