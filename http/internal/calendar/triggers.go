package calendar

import (
	"time"
)

type TimeEventTrigger struct {
	endtime time.Time
}

func CreateTimeEventTrigger(t time.Time) EventTrigger {
	return &TimeEventTrigger{endtime: t}
}

func (t *TimeEventTrigger) Start(handler func()) {
	if handler == nil {
		return
	}
	now := time.Now()
	if t.endtime.After(now) {
		duration := t.endtime.Sub(now)
		timer := time.NewTimer(duration)
		<-timer.C
		handler()
	}
}
