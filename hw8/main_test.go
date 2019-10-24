package main

import (
	"strconv"
	"testing"
)

type testEvent struct {
	invoke string
}

func (e *testEvent) Invoke() string {
	return e.invoke
}

func (e *testEvent) Serialize() string {
	return "testEvent"
}

func createCalendar() Calendar {
	cimpl := &calendImpl{}
	cimpl.Create()
	return cimpl
}

func TestSimpleAddGet(t *testing.T) {
	c := createCalendar()
	c.AddEvent(EventTrigger{date: "11:00"}, &testEvent{invoke: "1"})
	c.AddEvent(EventTrigger{date: "12:00"}, &testEvent{invoke: "2"})

	e0 := c.GetEvent(EventTrigger{date: "10:00"}, 0)
	if e0 != nil {
		t.Error("Failed Get not existing event")
	}
	e1 := c.GetEvent(EventTrigger{date: "11:00"}, 0)
	if e1 == nil {
		t.Error("Not found event at 11:00")
	}
	if e1.Invoke() != "1" {
		t.Error("Invalid invoke at 11:00")
	}
	e2 := c.GetEvent(EventTrigger{date: "12:00"}, 0)
	if e2 == nil {
		t.Error("Not found event at 12:00")
	}
	if e2.Invoke() != "2" {
		t.Error("Invalid invoke at 12:00")
	}
}

func TestAddGetAtOneTime(t *testing.T) {
	c := createCalendar()
	c.AddEvent(EventTrigger{date: "13:00"}, &testEvent{invoke: "1"})
	c.AddEvent(EventTrigger{date: "13:00"}, &testEvent{invoke: "2"})
	c.AddEvent(EventTrigger{date: "13:00"}, &testEvent{invoke: "3"})

	count := c.GetEventsCount(EventTrigger{date: "13:00"})
	if count != 3 {
		t.Error("Invalid events count")
	}
	for i := 0; i < count; i++ {
		e := c.GetEvent(EventTrigger{date: "13:00"}, i)
		if e == nil {
			t.Errorf("Not found event at 13:00, index %d", i)
		}
		if e.Invoke() != strconv.Itoa(i+1) {
			t.Errorf("Invalid invoke at 13:00, index %d", i)
		}
	}
}

func TestAddGetDelete(t *testing.T) {
	c := createCalendar()
	c.AddEvent(EventTrigger{date: "11:00"}, &testEvent{invoke: "1"})
	c.AddEvent(EventTrigger{date: "12:00"}, &testEvent{invoke: "2"})
	c.AddEvent(EventTrigger{date: "13:00"}, &testEvent{invoke: "3.1"})
	c.AddEvent(EventTrigger{date: "13:00"}, &testEvent{invoke: "3.2"})

	if c.GetEventsCount(EventTrigger{date: "13:00"}) != 2 {
		t.Error("Invalid events count at 13:00")
	}
	if c.DeleteAllEvents(EventTrigger{date: "10:00"}) != false {
		t.Error("Failed test DeleteAllEvents at 10:00")
	}
	if c.DeleteEvent(EventTrigger{date: "11:00"}, 1) != false {
		t.Error("Failed test DeleteEvent at 11:00, index 1")
	}
	if c.DeleteEvent(EventTrigger{date: "11:00"}, 0) != true {
		t.Error("Failed test DeleteEvent at 11:00, index 0")
	}
	if c.GetEventsCount(EventTrigger{date: "11:00"}) != 0 {
		t.Error("Invalid events count at 11:00")
	}
	if c.DeleteAllEvents(EventTrigger{date: "13:00"}) != true {
		t.Error("Failed test DeleteAllEvents at 13:00")
	}
	if c.GetEventsCount(EventTrigger{date: "13:00"}) != 0 {
		t.Error("Invalid events count at 13:00, attempt 2")
	}
	if c.GetEventsCount(EventTrigger{date: "12:00"}) != 1 {
		t.Error("Invalid events count at 12:00")
	}
}
