package rotaryphone

import (
    "time"

    "github.com/brian-armstrong/gpio"
)

type DebouncedWatcher struct {
    *gpio.Watcher
    bounceTime time.Duration
    lastValue map[uint]uint
    lastTime map[uint]time.Time
}

func NewDebouncedWatcher(bounceTime time.Duration) *DebouncedWatcher {
    return &DebouncedWatcher{
        gpio.NewWatcher(),
        bounceTime,
        make(map[uint]uint),
        make(map[uint]time.Time),
    }
}

func (dbw *DebouncedWatcher) Watch() (p uint, v uint) {
    for n := range dbw.Notification {
        if dbw.lastValue[n.Pin] != n.Value {
            dbw.lastValue[n.Pin] = n.Value
            if dbw.lastTime[n.Pin].Add(dbw.bounceTime).Before(time.Now()) {
                dbw.lastTime[n.Pin] = time.Now()
                return n.Pin, n.Value
            }
        }
    }
    return 0, 0
}

