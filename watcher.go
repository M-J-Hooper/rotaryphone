package rotaryphone

import (
	"log"
	"time"
)

type Event struct {
	Time    time.Time
	Payload interface{}
}

type DebouncedWatcher struct {
	BounceTime   time.Duration
	Notification chan interface{}
	report       chan Event
	notified     Event
	reported     Event
}

func NewDebouncedWatcher(n chan interface{}, b time.Duration) *DebouncedWatcher {
	emptyEvent := Event{
		Time:    time.Now(),
		Payload: nil,
	}
	dbw := DebouncedWatcher{
		Notification: n,
		BounceTime:   b,
		report:       make(chan Event),
		notified:     emptyEvent,
		reported:     emptyEvent,
	}
	go dbw.run()
	return &dbw
}

func Debounce(n chan interface{}, b time.Duration) chan interface{} {
	dbw := NewDebouncedWatcher(n, b)
	return dbw.Notification
}

func (dbw *DebouncedWatcher) Watch() interface{} {
	e := <-dbw.report
	return e.Payload
}

func (dbw *DebouncedWatcher) run() {
	for {
		n := dbw.notified
		select {
		case payload := <-dbw.Notification:
			log.Println("Received notification with", payload)
			if n.Payload != payload {
				dbw.notified = Event{
					Time:    time.Now(),
					Payload: payload,
				}
			}
		default:
			r := dbw.reported
			if n.Payload != nil && r.Payload != n.Payload {
				if time.Since(n.Time) > dbw.BounceTime {
					log.Println("Reporting debounced at", time.Now())
					dbw.reported = n
					dbw.report <- n
				}
			}
		}
	}
}
