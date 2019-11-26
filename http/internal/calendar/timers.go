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
}

func (t *timerimpl) Stop() {
	close(t.stop)
}

func createTimer(trigger string, timerend chan<- string) (*timerimpl, error) {
	p := &timeTriggerParser{}
	if err := p.Parse(trigger); err != nil {
		return nil, err
	}
	stopch := make(chan struct{})
	timer := &timerimpl{events: createEvents(), timerend: timerend, stop: stopch, id: trigger}

	triggerTimer := p.parsed
	p.SetNormalizedNow()
	if triggerTimer.Before(p.parsed) {
		return nil, fmt.Errorf("Cant set, time from past")
	}

	duration := triggerTimer.Sub(p.parsed)
	fmt.Println("Added trigger with duration:", duration.String())
	//fmt.Println(p.parsed.String())
	//fmt.Println(triggerTimer.String())

	go func(t *timerimpl, d time.Duration, stop <-chan struct{}) {
		timer := time.NewTimer(d)
		select {
		case <-timer.C:
			t.timerend <- t.id
		case <-stop:
		}
	}(timer, duration, stopch)
	return timer, nil
}
