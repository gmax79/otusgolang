package calendar

import (
	"time"
)

type timerimpl struct {
	events   CalendarEvents
	timerend chan<- string
	id       string
}

const dateLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

func createTimer(trigger string, timerend chan<- string) (*timerimpl, error) {
	p := &triggerParser{}
	if err := p.Parse(trigger, dateLayout); err != nil {
		return nil, err
	}
	timer := &timerimpl{events: createEvents(), timerend: timerend, id: trigger}
	//fmt.Println(p.parsed.String())
	//fmt.Println(time.Now())
	duration := time.Until(p.parsed)
	//fmt.Println(duration.String())
	go func(t *timerimpl) {
		timer := time.NewTimer(duration)
		<-timer.C
		t.timerend <- t.id
	}(timer)
	return timer, nil
}
