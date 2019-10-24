package main

// EventTrigger - key, when event happends
type EventTrigger struct {
	date string
}

// Event - interface for invoke event, and serialize it for restore later
type Event interface {
	Invoke() string
	Serialize() string
}

// Calendar - main object
type Calendar interface {
	AddEvent(t EventTrigger, e Event) bool
	GetEventsCount(t EventTrigger) int
	DeleteAllEvents(t EventTrigger) bool
	DeleteEvent(t EventTrigger, index int) bool
	GetEvent(t EventTrigger, index int) Event
}

// Calendar implementaion
type calendImpl struct {
	events map[EventTrigger][]Event
}

func (c *calendImpl) Create() {
	c.events = make(map[EventTrigger][]Event)
}

func (c *calendImpl) AddEvent(t EventTrigger, e Event) bool {
	if e == nil {
		return false
	}
	events, ok := c.events[t]
	if !ok {
		c.events[t] = []Event{e}
	} else {
		c.events[t] = append(events, e)
	}
	return true
}

func (c *calendImpl) GetEventsCount(t EventTrigger) int {
	events, ok := c.events[t]
	if !ok {
		return 0
	}
	return len(events)
}

func (c *calendImpl) DeleteEvent(t EventTrigger, index int) bool {
	events, ok := c.events[t]
	if !ok || index < 0 || index >= len(events) {
		return false
	}
	c.events[t] = append(events[:index], events[index+1:]...)
	return true
}

func (c *calendImpl) DeleteAllEvents(t EventTrigger) bool {
	if _, ok := c.events[t]; !ok {
		return false
	}
	delete(c.events, t)
	return true
}

func (c *calendImpl) GetEvent(t EventTrigger, index int) Event {
	events, ok := c.events[t]
	if !ok || index < 0 || index >= len(events) {
		return nil
	}
	return events[index]
}
