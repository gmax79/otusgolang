package calendar

import (
	"fmt"
	"time"
)

type timerimpl struct {
	timerend  chan<- string
	stop      <-chan struct{}
	id        string
	alerttime time.Time
}

func createTimer(id string, alert date, timerend chan<- string, stopch <-chan struct{}) (*timerimpl, error) {
	a := alert.Value()
	if a.Before(alert.SetNow()) {
		return nil, fmt.Errorf("Cant set, time from past")
	}
	timer := &timerimpl{timerend: timerend, stop: stopch, id: id, alerttime: a}
	go func(t *timerimpl) {
		var p date
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
