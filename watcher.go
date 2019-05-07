package rotaryphone

import (
    "time"
    "fmt"
    "math/rand"

    "github.com/brian-armstrong/gpio"
)

type Watcher interface {
    Watch() (uint, uint)
    AddPin(uint)
    Close()
}

type DebouncedWatcher struct {
    Watcher
    BounceTime time.Duration
    CurrentValue map[uint]uint
}

func NewDebouncedWatcher(BounceTime time.Duration) *DebouncedWatcher {
    return &DebouncedWatcher{
        gpio.NewWatcher(),
        BounceTime,
        make(map[uint]uint),
    }
}

func (dbw *DebouncedWatcher) Watch() (uint, uint) {
    var waiting bool
    stop := make(chan struct{})
    stable := make(chan gpio.WatcherNotification)

    notify := make(chan gpio.WatcherNotification)
    go func (notify chan gpio.WatcherNotification) {
        for {
            p, v := dbw.Watcher.Watch()
            notify <-gpio.WatcherNotification{p, v}
        }
    }(notify)

    for {
        select {
        case n := <-notify:
            fmt.Println("Notification", n.Pin, n.Value)
            // Stop goroutine waiting for stability
            if waiting {
                stop <-struct{}{}
                waiting = false
            }
            if n.Value != dbw.CurrentValue[n.Pin] {
                go func(n gpio.WatcherNotification, stable chan gpio.WatcherNotification, stop chan struct{}) {
                    id := rand.Intn(10000)
                    fmt.Println(id, "Starting goroutine")
                    start := time.Now()
                    for {
                        select {
                        case <-stop:
                            fmt.Println(id, "Got cancel")
                            return
                        default:
                            if time.Since(start) > dbw.BounceTime {
                                fmt.Println(id, "Sending stable", n.Pin, n.Value)
                                stable <-n
                                return
                            }
                        }
                    }
                }(n, stable, stop)
                waiting = true
            }
        case n := <-stable:
            fmt.Println("Got stable", n.Pin, n.Value)
            dbw.CurrentValue[n.Pin] = n.Value
            return n.Pin, n.Value
        }
    }
    return 0, 0
}
