package calendar

import (
	"fmt"
	"time"

	"github.com/gmax79/otusgolang/rmq/internal/simple"
)

type timerimpl struct {
	timerend chan<- simple.Date
	stop     <-chan struct{}
	alert    simple.Date
	duration time.Duration
}

func createTimer(alert simple.Date, timerend chan<- simple.Date, stopch <-chan struct{}) error {
	var p simple.Date
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
