package calendar

import (
	"github.com/gmax79/otusgolang/microservices/internal/objects"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	"github.com/gmax79/otusgolang/microservices/internal/storage"
)

type eventsimpl struct {
	d  simple.Date
	db *storage.DbProvider
}

func createEvents(alert simple.Date, dbc *storage.DbProvider) Events {
	return &eventsimpl{db: dbc, d: alert}
}

func (t *eventsimpl) AddEvent(e objects.Event) error {
	return t.db.AddEvent(t.d, string(e))
}

func (t *eventsimpl) GetEventsCount() (int, error) {
	return t.db.GetEventsCount(t.d)
}

func (t *eventsimpl) DeleteEventIndex(index int) error {
	return t.db.DeleteEventIndex(t.d, index)
}

func (t *eventsimpl) DeleteEvent(e objects.Event) error {
	return t.db.DeleteEvent(t.d, e)
}

func (t *eventsimpl) GetEvent(index int) (objects.Event, error) {
	return t.db.GetEvent(t.d, index)
}

func (t *eventsimpl) MoveEvent(e objects.Event, to simple.Date) error {
	return t.db.MoveEvent(t.d, e, to)
}
