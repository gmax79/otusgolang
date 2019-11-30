package calendar

type eventsimpl struct {
	events []Event
}

func createEvents() Events {
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

func (t *eventsimpl) FindEvent(name string) int {
	for i, e := range t.events {
		if e.GetName() == name {
			return i
		}
	}
	return -1
}
