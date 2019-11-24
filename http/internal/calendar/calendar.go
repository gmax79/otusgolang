package calendar

// Event - interface for invoke event
type Event interface {
	Invoke()
}

// CalendarEvents - contains all events per trigger
type CalendarEvents interface {
	AddEvent(e Event) bool
	GetEventsCount() int
	DeleteEvent(index int) bool
	GetEvent(index int) Event
	ReplaceEvent(index int, e Event) bool
	Invoke()
}

// Calendar - main object, contains all triggers and objects
type Calendar interface {
	AddTrigger(trigger string) (CalendarEvents, error)
	DeleteTrigger(trigger string) bool
	GetEvents(trigger string) CalendarEvents
	GetTriggers() []string
}

// Create - create calendar instance
func Create() Calendar {
	return createCalendar()
}
