package calendar

import (
	"fmt"
	"time"
)

type timerimpl struct {
	timerend chan<- date
	stop     <-chan struct{}
	alert    date
	duration time.Duration
}

func createTimer(alert date, timerend chan<- date, stopch <-chan struct{}) error {
	var p date
	p.SetNow()
	a := alert.Value()
	if a.Before(p.Value()) {
		return fmt.Errorf("Cant set, time from past")
	}
	timer := &timerimpl{timerend: timerend, stop: stopch, alert: alert}
	timer.duration = a.Sub(p.Value())
	go func(t *timerimpl) {
		timer := time.NewTimer(t.duration)
		select {
		case <-timer.C:
			t.timerend <- t.alert
		case <-t.stop:
		}
	}(timer)
	return nil
}
