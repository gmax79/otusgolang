package calendar

import (
	"fmt"
	"testing"
	"time"
)

type testEvent struct {
}

func (e *testEvent) Invoke() {
	fmt.Println("Event!!!")
}

func (e *testEvent) GetName() string {
	return ""
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

func TestTimeParser(t *testing.T) {
	var p Date
	if p.ParseDate("10:20:30") != nil {
		t.Fatal("Error at correct time")
	}
	if p.Hour != 10 || p.Minute != 20 || p.Second != 30 {
		t.Fatal("Invalid parsed time")
	}
	if p.ParseDate("20:40") != nil {
		t.Fatal("Error at correct time 2")
	}
	if p.Hour != 20 || p.Minute != 40 || p.Second != 0 {
		t.Fatal("Invalid parsed time 2")
	}
}

func TestBaseMethods(t *testing.T) {
	c := Create()
	_, err := c.AddTrigger("")
	if err == nil {
		t.Error("error not can be nil")
	}
	ce := c.GetEvents("")
	if ce != nil {
		t.Error("events must be nil")
	}
	tr := DurationToTimeString(time.Second * 15)
	c.AddTrigger(tr)
	tr2 := DurationToTimeString(time.Second * 25)
	c.AddTrigger(tr2)

	trs := c.GetTriggers()
	if len(trs) != 2 {
		t.Fatalf("Triggers must be 2")
	}
}

func TestTimerEvents(t *testing.T) {
	c := Create()
	tr := DurationToTimeString(time.Second * 3)
	events, err := c.AddTrigger(tr)
	if err != nil {
		t.Fatal(err)
	}
	events.AddEvent(&testEvent{})
	time.Sleep(time.Second * 4)
}
