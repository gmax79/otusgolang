package calendar

import (
	"fmt"
	"time"
)

// DateLayout - calendar format of data with time
const DateLayout = "2006-01-02 15:04:05"

// Event - information about event
type Event string

// Date - calendar date with time
type Date struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}

func (d *Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d %d:%02d:%02d", d.Year, d.Month, d.Day, d.Hour, d.Minute, d.Second)
}

// Events - contains all events per trigger
type Events interface {
	AddEvent(e Event) error
	GetEventsCount() (int, error)
	DeleteEventIndex(index int) error
	DeleteEvent(e Event) error
	GetEvent(index int) (Event, error)
	MoveEvent(e Event, to Date) error
}

// SearchParameters - custom filters to search events
type SearchParameters struct {
	Year, Month, Week, Day int
}

// Calendar - main object, contains all triggers and objects
type Calendar interface {
	AddTrigger(date Date) (Events, error)
	DeleteTrigger(date Date) error
	GetEvents(date Date) (Events, error)
	GetTriggers() ([]Date, error)
	FindEvents(parameters SearchParameters) ([]Event, error)
}

// Create - create calendar instance
func Create(psqlConnect string) (Calendar, error) {
	return createCalendar(psqlConnect)
}

// DurationToTimeString - get time parameter, now + duration
func DurationToTimeString(d time.Duration) string {
	t := time.Now().Add(d)
	s := t.String()
	return s[:len(DateLayout)]
}

// ParseValidDate - create calendar date from string and validate
func ParseValidDate(trigger string) (Date, error) {
	var err error
	var zero Date
	var d date
	if err = d.ParseDate(trigger); err != nil {
		return zero, err
	}
	if err = d.Valid(); err != nil {
		return zero, err
	}
	return d.d, nil
}

// ParseDate - create calendar date from string
func ParseDate(trigger string) (Date, error) {
	var err error
	var zero Date
	var d date
	if err = d.ParseDate(trigger); err != nil {
		return zero, err
	}
	return d.d, nil
}
