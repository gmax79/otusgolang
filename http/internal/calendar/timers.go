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

func createTimer(trigger string, timerend chan<- string) (*timerimpl, error) {
	var p Date
	if err := p.ParseDate(trigger); err != nil {
		return nil, err
	}
	stopch := make(chan struct{})
	timer := &timerimpl{events: createEvents(), timerend: timerend, stop: stopch, id: trigger, alerttime: p.Value()}

	if timer.alerttime.Before(p.SetNow()) {
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
