package rotaryphone

import (
    "time"
    "strconv"
)

type OutputAdapter interface {
    Call(string)
    Hangup()
    Debug()
}

type Rotary struct {
    digitTimeout time.Duration
    output OutputAdapter
    dial Dial
    latch Latch
}

func NewRotary(digitTimeout time.Duration) *Rotary {
    r := &Rotary{
        digitTimeout,
        NewOfonoAdapter(),
        *NewDial(),
        *NewLatch(),
    }
    return r
}

func (r Rotary) Run() {
    var dialing bool
    var number string
    lastDigit := time.Now()
    for {
        select {
        case digit := <-r.dial.digit:
            if dialing {
                number += strconv.Itoa(digit)
                lastDigit = time.Now()
                println("New number is " + number)
            }
        case dialing = <-r.latch.active:
            if !dialing {
                println("Handset on the latch")
                r.output.Hangup()
                number = ""
            } else {
                println("Handset off the latch")
            }
        default:
            if time.Since(lastDigit) > r.digitTimeout && len(number) > 0 {
                println("Calling " + number)
                r.output.Call(number)
                lastDigit = time.Now()
            }
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func (r Rotary) Debug() {
    r.output.Debug()
}
