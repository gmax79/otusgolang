package calendar

// Calendar implementaion
type calendarImpl struct {
	finished   chan date
	stoptimers chan struct{}
	db         *dbProvder
	timersset  map[Date]struct{}
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
	newcalendar.timersset = make(map[Date]struct{})
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
	if _, ok := c.timersset[d]; !ok {
		err = createTimer(t, c.finished, c.stoptimers)
		if err != nil {
			return nil, err
		}
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
	delete(c.timersset, d)
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
	return c.db.FindEvents(parameters)
}
