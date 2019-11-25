package calendar

import (
	"fmt"
	"time"
)

type timerimpl struct {
	events   Events
	timerend chan<- string
	id       string
}

func createTimer(trigger string, timerend chan<- string) (*timerimpl, error) {
	p := &timeTriggerParser{}
	if err := p.Parse(trigger); err != nil {
		return nil, err
	}
	timer := &timerimpl{events: createEvents(), timerend: timerend, id: trigger}

	triggerTimer := p.parsed
	p.SetNormalizedNow()
	if triggerTimer.Before(p.parsed) {
		return nil, fmt.Errorf("Cant set, time from past")
	}

	duration := triggerTimer.Sub(p.parsed)
	fmt.Println("Added trigger with duration:", duration.String())
	//fmt.Println(p.parsed.String())
	//fmt.Println(triggerTimer.String())

	go func(t *timerimpl, d time.Duration) {
		timer := time.NewTimer(d)
		<-timer.C
		t.timerend <- t.id
	}(timer, duration)
	return timer, nil
}
