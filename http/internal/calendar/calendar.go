package calendar

// Event - interface for invoke event
type Event interface {
	Invoke()
}

// EventTrigger - interface for objects, implements some event
// trigger must call handler() when event happends
type EventTrigger interface {
	Start(handler func())
}

// CalendarTrigger - contains all events per EventTrigger
type CalendarTrigger interface {
	AddEvent(e Event) bool
	GetEventsCount() int
	DeleteEvent(index int) bool
	GetEvent(index int) Event
	ReplaceEvent(index int, e Event) bool
	Invoke()
}

// Calendar - main object, contains all triggers and objects
type Calendar interface {
	AddTrigger(t EventTrigger, ct CalendarTrigger)
	GetTriggersCount() int
	DeleteTrigger(index int) bool
	GetTrigger(index int) CalendarTrigger
}

// CreateCalendar - create calendar instance
func CreateCalendar() Calendar {
	return &calendarImpl{}
}

// CreateCalendarTrigger - holder of events collection
func CreateCalendarTrigger() CalendarTrigger {
	return &triggerImpl{}
}
