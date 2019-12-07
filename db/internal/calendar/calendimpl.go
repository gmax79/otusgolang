package calendar

import (
	"time"
)

// Calendar implementaion
type calendarImpl struct {
	finished   chan date
	stoptimers chan struct{}
	db         *dbProvder
}

func createCalendar(psqlConnect string) (Calendar, error) {
	db, err := connectToDatabase(psqlConnect)
	if err != nil {
		return nil, err
	}
	var checkdb dbSchema
	err = checkdb.CheckOrCreateSchema(db)
	if err != nil {
		return nil, err
	}
	newcalendar := &calendarImpl{}
	newcalendar.finished = make(chan date, 1)
	newcalendar.stoptimers = make(chan struct{})
	newcalendar.db = getProvider(db)
	go func(c *calendarImpl) {
		for {
			id := <-c.finished
			c.db.Invoke(id.String()) //todo
		}
	}(newcalendar)
	return newcalendar, nil
}

func valid(d Date) (date, error) {
	var t date
	t.d = d
	return t, t.Valid()
}

func (c *calendarImpl) AddTrigger(d Date) (Events, error) {
	var err error
	var t date
	if t, err = valid(d); err != nil {
		return nil, err
	}
	err = c.db.AddTrigger(t)
	if err != nil {
		return nil, err
	}
	err = createTimer(t, c.finished, c.stoptimers)
	if err != nil {
		return nil, err
	}
	return createEvents(t, c.db), nil
}

func (c *calendarImpl) GetTriggers() ([]Date, error) {
	return c.db.GetTriggers()
}

func (c *calendarImpl) DeleteTrigger(d Date) error {
	var err error
	var t date
	if t, err = valid(d); err != nil {
		return err
	}
	return c.db.DeleteTrigger(t)
}

func (c *calendarImpl) GetEvents(d Date) (Events, error) {
	var err error
	var t date
	if t, err = valid(d); err != nil {
		return nil, err
	}
	return &eventsimpl{db: c.db, d: t}, nil
}

func (c *calendarImpl) FindEvents(parameters SearchParameters) ([]Event, error) {
	events := make([]Event, 0, 10)
	/*for _, t := range c.triggers {
		if checkSearchParameters(t.alerttime, parameters) {
			//todo optimize
			count := t.events.GetEventsCount()
			for i := 0; i < count; i++ {
				events = append(events, t.events.GetEvent(i))
			}
		}
	}*/
	return events, nil
}

func checkSearchParameters(t time.Time, p SearchParameters) bool {
	if p.Year <= 0 {
		return false
	}
	if p.Week > 0 {
		if p.Month == 0 && p.Day == 0 {
			year, week := t.ISOWeek()
			if year == p.Year && week == p.Week {
				return true
			}
		}
		return false
	}
	if p.Month > 0 {
		if p.Day == 0 {
			if p.Month == int(t.Month()) {
				return true
			}
		}
		if p.Day > 0 {
			if p.Year == t.Year() && p.Month == int(t.Month()) && p.Day == t.Day() {
				return true
			}
		}
	}
	return false
}
