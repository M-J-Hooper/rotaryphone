package rotaryphone

import "time"

//3.3v  => white wire
const LatchActivePin = 18 // => dark green wire

type Latch struct {
    Active chan bool
}

func NewLatch() *Latch {
    l := &Latch{make(chan bool)}
    go l.Run()
    return l
}

func (l Latch) Run() {
    watcher := NewDebouncedWatcher(100 * time.Millisecond)
    watcher.AddPin(LatchActivePin)
    defer watcher.Close()

    for {
        println("Calling latch watch")
        pin, value := watcher.Watch()
        println("After latch watch")
        if pin == LatchActivePin {
            if value == 1 {
                l.Active <-true
            } else {
                l.Active <-false
            }
        }
    }
}
