package calendar

// implements CalandarEvents
type eventsimpl struct {
	events []Event
}

func createEvents() CalendarEvents {
	return &eventsimpl{}
}

func (t *eventsimpl) Invoke() {
	for _, e := range t.events {
		e.Invoke()
	}
}

func (t *eventsimpl) AddEvent(e Event) bool {
	if e == nil {
		return false
	}
	t.events = append(t.events, e)
	return true
}

func (t *eventsimpl) GetEventsCount() int {
	return len(t.events)
}

func (t *eventsimpl) DeleteEvent(index int) bool {
	if index >= 0 && index < len(t.events) {
		t.events = append(t.events[:index], t.events[index+1:]...)
		return true
	}
	return false
}

func (t *eventsimpl) GetEvent(index int) Event {
	if index >= 0 && index < len(t.events) {
		return t.events[index]
	}
	return nil
}

func (t *eventsimpl) ReplaceEvent(index int, e Event) bool {
	if e == nil {
		return false
	}
	if t.DeleteEvent(index) {
		t.AddEvent(e)
		return true
	}
	return false
}
