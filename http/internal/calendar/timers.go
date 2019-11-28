package calendar

import (
	"fmt"
	"time"
)

type timerimpl struct {
	events   Events
	timerend chan<- string
	stop     chan struct{}
	id       string
	alert    time.Time
}

func (t *timerimpl) Stop() {
	close(t.stop)
}

func (t *timerimpl) String() string {
	return t.alert.String()
}

func createTimer(trigger string, timerend chan<- string) (*timerimpl, error) {
	p := timeParser{}
	if err := p.Parse(trigger); err != nil {
		return nil, err
	}
	stopch := make(chan struct{})
	timer := &timerimpl{events: createEvents(), timerend: timerend, stop: stopch, id: trigger, alert: p.Value()}

	if timer.alert.Before(p.Now()) {
		return nil, fmt.Errorf("Cant set, time from past")
	}

	go func(t *timerimpl) {
		p := timeParser{}
		duration := t.alert.Sub(p.Now())
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			t.timerend <- t.id
		case <-t.stop:
		}
	}(timer)
	return timer, nil
}
