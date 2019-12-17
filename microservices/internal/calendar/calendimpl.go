package calendar

import (
	"github.com/gmax79/otusgolang/microservices/internal/objects"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	"github.com/gmax79/otusgolang/microservices/internal/storage"
)

// Calendar implementaion
type calendarImpl struct {
	finished   chan simple.Date
	stoptimers chan struct{}
	db         *storage.DbProvider
	timersset  map[simple.Date]struct{}
}

func createCalendar(psqlConnect string) (Calendar, error) {
	db, err := storage.ConnectToDatabase(psqlConnect)
	if err != nil {
		return nil, err
	}
	var checkdb storage.DbSchema
	err = checkdb.CheckOrCreateSchema(db)
	if err != nil {
		return nil, err
	}
	newcalendar := &calendarImpl{}
	newcalendar.finished = make(chan simple.Date, 1)
	newcalendar.stoptimers = make(chan struct{})
	newcalendar.db = storage.CreateProvider(db)
	newcalendar.timersset = make(map[simple.Date]struct{})
	go func(c *calendarImpl) {
		for {
			id := <-c.finished
			c.db.Invoke(id.String()) //todo
		}
	}(newcalendar)
	return newcalendar, nil
}

func (c *calendarImpl) AddTrigger(d simple.Date) (Events, error) {
	var err error
	if err = d.Valid(); err != nil {
		return nil, err
	}
	if _, ok := c.timersset[d]; !ok {
		err = createTimer(d, c.finished, c.stoptimers)
		if err != nil {
			return nil, err
		}
	}
	return createEvents(d, c.db), nil
}

func (c *calendarImpl) GetTriggers() ([]simple.Date, error) {
	return c.db.GetTriggers()
}

func (c *calendarImpl) DeleteTrigger(d simple.Date) error {
	var err error
	if err = d.Valid(); err != nil {
		return err
	}
	delete(c.timersset, d)
	return c.db.DeleteTrigger(d)
}

func (c *calendarImpl) GetEvents(d simple.Date) (Events, error) {
	var err error
	if err = d.Valid(); err != nil {
		return nil, err
	}
	return &eventsimpl{db: c.db, d: d}, nil
}

func (c *calendarImpl) FindEvents(parameters objects.SearchParameters) ([]objects.Event, error) {
	return c.db.FindEvents(parameters)
}

func (c *calendarImpl) SinceEvents(from simple.Date) ([]objects.Event, error) {
	return c.db.SinceEvents(from)
}
