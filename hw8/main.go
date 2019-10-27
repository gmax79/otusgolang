|package main

type EventTrigger interface {
	SetEventCallback(f func())
}

// Event - interface for invoke event
type Event interface {
	Invoke()
}

type CalendarTrigger interface {
	AddEvent(e Event) bool
	GetEventsCount() int
	DeleteEvent(index int) bool
	GetEvent(index int) Event
	ReplaceEvent(index int, e Event) bool
}

// Calendar - main object
type Calendar interface {
	AddTrigger(t EventTrigger)
	FindTrigger(t EventTrigger) CalendarTrigger
	DeleteTrigger(t EventTrigger)
}

// Calendar implementaion
type calendarImpl struct {
	triggers map[EventTrigger]triggerImpl
}

type triggerImpl struct {
	events []Event
}

// CreateCalendar - create calendar instance
func CreateCalendar() Calendar {
	return &calendarImpl{ triggers: make(map[EventTrigger]triggerImpl)}
}

func (c *calendarImpl) AddTrigger(t EventTrigger) {
	callback := func() {
	}

	c.trigger[t] = triggerImpl{}
}

func (c *calendarImpl) FindTrigger(t EventTrigger) CalendarTrigger {
	return nil
}

func (c *calendarImpl) DeleteTrigger(t EventTrigger) {
	delete(c.triggers, t)
}

func (t *triggerImpl) AddEvent(e Event) bool {
	if e == nil {
		return false
	}
	t.events = append(t.events, e)
	return true
}

func (t *triggerImpl) GetEventsCount() int {
	return len(t.events)
}

func (t *triggerImpl) DeleteEvent(index int) bool {
	if index >= 0 && index < len(t.events) {
		t.events = append(t.events[:index], t.events[index+1:]...)
		return true
	}
	return false
}

func (t *triggerImpl) GetEvent(index int) Event {
	if index >= 0 && index < len(t.events) {
		return t.events[index]
	}
	return nil
}

func (t *triggerImpl) ReplaceEvent(index int, e Event) bool {
	if e == nil {
		return false
	}
	if t.DeleteEvent(index) {
		t.AddEvent(e)
		return true
	}
	return false
}
