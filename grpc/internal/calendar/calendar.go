package calendar

import "time"

// DateLayout - calendar format of data with time
const DateLayout = "2006-01-02 15:04:05"

// Event - interface for invoke event
type Event interface {
	Invoke()
	GetName() string
}

// Events - contains all events per trigger
type Events interface {
	AddEvent(e Event) bool
	GetEventsCount() int
	DeleteEvent(index int) bool
	GetEvent(index int) Event
	FindEvent(name string) int
	Invoke()
}

// SearchParameters - custom filters to search events
type SearchParameters struct {
	Year  int
	Month int
	Week  int
	Day   int
}

// Calendar - main object, contains all triggers and objects
type Calendar interface {
	AddTrigger(trigger string) (Events, error)
	DeleteTrigger(trigger string) bool
	GetEvents(trigger string) Events
	GetTriggers() []string
	GetTriggerAlert(trigger string) (time.Time, bool)
	FindEvents(parameters SearchParameters) ([]Event, error)
}

// Create - create calendar instance
func Create() Calendar {
	return createCalendar()
}

// DurationToTimeString - get time parameter, now + duration
func DurationToTimeString(d time.Duration) string {
	t := time.Now().Add(d)
	s := t.String()
	return s[:len(DateLayout)]
}
