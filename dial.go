package rotaryphone

import "time"

//3.3v  => white wire
const DialIncrementPin = 14 //=> orange wire
const DialActivePin = 15    //=> brown wire

type Dial struct {
    digit chan int
}

func NewDial() *Dial {
    d := &Dial{make(chan int)}
    go d.Run()
    return d
}

func (d Dial) Run() {
    watcher := NewDebouncedWatcher(10 * time.Millisecond)
    watcher.AddPin(DialIncrementPin)
    watcher.AddPin(DialActivePin)
    defer watcher.Close()

    var active bool
    var count int
    for {
        println("Calling dial watch")
        pin, value := watcher.Watch()
        println("After dial watch")
        if pin == DialActivePin {
            if value == 1 {
                active = true
            } else {
                active = false
                if count > 0 {
                    if count >= 10 { count = 0 }
                    d.digit <-count
                }
                count = 0
            }
        } else if pin == DialIncrementPin && value == 0 {
            if active {
                count++
            }
        }
    }
}
