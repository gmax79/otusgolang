package calendar

type eventsimpl struct {
	d  date
	db *dbProvder
}

func createEvents(alert date, dbc *dbProvder) Events {
	return &eventsimpl{db: dbc, d: alert}
}

func (t *eventsimpl) AddEvent(e Event) error {
	return t.db.AddEvent(t.d, string(e))
}

func (t *eventsimpl) GetEventsCount() (int, error) {
	return t.db.GetEventsCount(t.d)
}

func (t *eventsimpl) DeleteEvent(index int) bool {
	/*if index >= 0 && index < len(t.events) {
		t.events = append(t.events[:index], t.events[index+1:]...)
		return true
	}*/
	return false
}

func (t *eventsimpl) GetEvent(index int) Event {
	/*if index >= 0 && index < len(t.events) {
		return t.events[index]
	}*/
	return ""
}

func (t *eventsimpl) FindEvent(name string) int {
	/*for i, e := range t.events {
		if e.GetName() == name {
			return i
		}
	}*/
	return -1
}
