package calendar

import (
	"time"
)

type timerimpl struct {
	events   CalendarEvents
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
	duration := triggerTimer.Sub(p.parsed)

	//fmt.Println(p.parsed.String())
	//fmt.Println(triggerTimer.String())
	//fmt.Println(duration.String())

	go func(t *timerimpl, d time.Duration) {
		timer := time.NewTimer(d)
		<-timer.C
		t.timerend <- t.id
	}(timer, duration)
	return timer, nil
}
