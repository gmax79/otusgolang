package calendar

import (
	"time"

	"github.com/gmax79/otusgolang/rmq/internal/objects"
	"github.com/gmax79/otusgolang/rmq/internal/simple"
)

// Events - contains all events per trigger
type Events interface {
	AddEvent(e objects.Event) error
	GetEventsCount() (int, error)
	DeleteEventIndex(index int) error
	DeleteEvent(e objects.Event) error
	GetEvent(index int) (objects.Event, error)
	MoveEvent(e objects.Event, to simple.Date) error
}

// Calendar - main object, contains all triggers and objects
type Calendar interface {
	AddTrigger(date simple.Date) (Events, error)
	DeleteTrigger(date simple.Date) error
	GetEvents(date simple.Date) (Events, error)
	GetTriggers() ([]simple.Date, error)
	FindEvents(parameters objects.SearchParameters) ([]objects.Event, error)
}

// Create - create calendar instance
func Create(connect string) (Calendar, error) {
	return createCalendar(connect)
}

// DurationToTimeString - get time parameter, now + duration
func DurationToTimeString(d time.Duration) string {
	t := time.Now().Add(d)
	s := t.String()
	const layout = "2006-01-02 15:04:05"
	return s[:len(layout)]
}
