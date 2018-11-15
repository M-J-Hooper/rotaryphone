package rotaryphone

import "github.com/brian-armstrong/gpio"

//3.3v  => white wire
const DialIncrementPin = 18 //=> orange wire
const DialTogglePin = 24    //=> brown wire
const LatchPin = 8

type GpioAdapter struct {
    digitChannel chan int
    hangupChannel chan struct{}
}

func NewGpioAdapter() *GpioAdapter {
    return &GpioAdapter{make(chan int), make(chan struct{})}
}

func (g GpioAdapter) Run() {
    watcher := gpio.NewWatcher()

    watcher.AddPin(DialIncrementPin)
    watcher.AddPin(DialTogglePin)
    watcher.AddPin(LatchPin)

    defer watcher.Close()

    active := false
    count := 0
    for {
        pin, value := watcher.Watch()
        if pin == DialTogglePin {
            if value == 1 {
                active = true
            } else {
                active = false
                if count > 0 {
                    if count >= 10 { count = 0 }
                    g.digitChannel <-count
                }
                count = 0
            }
        } else if pin == DialIncrementPin && value == 0 {
            if active {
                count++
            }
        } else if pin == LatchPin && value == 0 {
            g.hangupChannel <-struct{}{}
        }
    }
}

func (g GpioAdapter) GetDigitChannel() chan int {
    return g.digitChannel
}

func (g GpioAdapter) GetHangupChannel() chan struct{} {
    return g.hangupChannel
}

