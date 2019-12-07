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
	AddEvent(e Event) bool
	GetEventsCount() int
	DeleteEvent(index int) bool
	GetEvent(index int) Event
	FindEvent(name string) int
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

// ParseDate - create calendar date from string
func ParseDate(trigger string) (dd Date, err error) {
	var d date
	if err = d.ParseDate(trigger); err != nil {
		return
	}
	if err = d.Valid(); err != nil {
		return
	}
	return d.d, nil
}
