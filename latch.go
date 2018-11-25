package rotaryphone

import "time"

//3.3v  => white wire
const LatchActivePin = 21 // => dark green wire

type Latch struct {
    active chan bool
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
        pin, value := watcher.Watch()
        if pin == LatchActivePin {
            if value == 1 {
                l.active <-true
            } else {
                l.active <-false
            }
        }
    }
}
