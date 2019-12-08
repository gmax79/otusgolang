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

func (t *eventsimpl) DeleteEventIndex(index int) error {
	return t.db.DeleteEventIndex(t.d, index)
}

func (t *eventsimpl) DeleteEvent(e Event) error {
	return t.db.DeleteEvent(t.d, e)
}

func (t *eventsimpl) GetEvent(index int) (Event, error) {
	return t.db.GetEvent(t.d, index)
}

func (t *eventsimpl) MoveEvent(e Event, to Date) error {
	var newdate date
	newdate.d = to
	return t.db.MoveEvent(t.d, e, newdate)
}
