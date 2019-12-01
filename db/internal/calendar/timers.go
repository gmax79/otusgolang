package calendar

import (
	"fmt"
	"time"
)

type timerimpl struct {
	events    Events
	timerend  chan<- string
	stop      chan struct{}
	id        string
	alerttime time.Time
}

func (t *timerimpl) Stop() {
	close(t.stop)
}

func (t *timerimpl) String() string {
	return t.alerttime.String()
}

func createTimer(date Date, timerend chan<- string) (*timerimpl, error) {
	stopch := make(chan struct{})
	timer := &timerimpl{events: createEvents(), timerend: timerend, stop: stopch, id: date.String(), alerttime: date.Value()}

	if timer.alerttime.Before(date.SetNow()) {
		return nil, fmt.Errorf("Cant set, time from past")
	}

	go func(t *timerimpl) {
		var p Date
		duration := t.alerttime.Sub(p.SetNow())
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			t.timerend <- t.id
		case <-t.stop:
		}
	}(timer)
	return timer, nil
}
