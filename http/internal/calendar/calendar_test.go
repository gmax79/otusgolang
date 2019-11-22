package calendar

import (
	"testing"
	"time"
)

type testEvent struct {
	value int
	summ  *int
}

func (e *testEvent) Invoke() {
	*e.summ = *e.summ + e.value
}

type testTimerEventTrigger struct {
	duration time.Duration
}

func (t *testTimerEventTrigger) Start(f func()) {
	timer := time.NewTimer(t.duration)
	go func() {
		<-timer.C
		f()
	}()
}

func TestTimerEvents(t *testing.T) {
	c := CreateCalendar()
	trigger := &testTimerEventTrigger{duration: time.Second * 2}
	ctrigger := CreateCalendarTrigger()
	result := 0
	ctrigger.AddEvent(&testEvent{value: 1, summ: &result})
	ctrigger.AddEvent(&testEvent{value: 2, summ: &result})
	c.AddTrigger(trigger, ctrigger)
	time.Sleep(time.Second * 3)
	if result != 3 {
		t.Error("Failed TestTimerEvents")
	}
	if c.GetTriggersCount() != 1 {
		t.Error("Failed GetTriggersCount")
	}
	ct := c.GetTrigger(0)
	if ct.GetEventsCount() != 2 {
		t.Error("Failed GetEventsCount")
	}
	if ct.DeleteEvent(2) {
		t.Error("Failed DeleteEvent")
	}
	if !ct.DeleteEvent(0) {
		t.Error("Failed DeleteEvent with correct index")
	}
	if ct.GetEventsCount() != 1 {
		t.Error("Failed GetEventsCount after delete event")
	}
}
