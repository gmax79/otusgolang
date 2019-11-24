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
	p := &triggerParser{}
	if err := p.Parse(trigger); err != nil {
		return nil, err
	}
	timer := &timerimpl{events: createEvents(), timerend: timerend, id: trigger}

	triggerTimer := p.parsed
	p.NormalizeNow()
	duration := triggerTimer.Sub(p.parsed)

	//fmt.Println(p.parsed.String())
	//fmt.Println(triggerTimer.String())
	//fmt.Println(duration.String())

	go func(t *timerimpl) {
		timer := time.NewTimer(duration)
		<-timer.C
		t.timerend <- t.id
	}(timer)
	return timer, nil
}
