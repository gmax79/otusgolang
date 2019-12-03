package calendar

import (
	"fmt"
	"sync"
	"time"
)

// Calendar implementaion
type calendarImpl struct {
	triggers map[string]*timerimpl
	finished chan string
	m        *sync.Mutex
	db       *dbHandler
}

func createCalendar(psqlConnect string) (Calendar, error) {
	db, err := dbconnect(psqlConnect)
	if err != nil {
		return nil, err
	}
	err = db.CheckOrCreateSchema()
	if err != nil {
		return nil, err
	}
	newcalendar := &calendarImpl{}
	newcalendar.triggers = make(map[string]*timerimpl)
	newcalendar.finished = make(chan string, 1)
	newcalendar.m = &sync.Mutex{}
	newcalendar.db = db
	go func(c *calendarImpl) {
		for {
			id := <-c.finished
			t, ok := c.triggers[id]
			if !ok {
				fmt.Printf("Error: trigger %s not found\n", id)
			} else {
				fmt.Printf("Processed %s trigger\n", id)
				t.events.Invoke()
			}
		}
	}(newcalendar)
	return newcalendar, nil
}

func (c *calendarImpl) AddTrigger(trigger string) (Events, error) {
	c.m.Lock()
	defer c.m.Unlock()

	var p Date
	if err := p.ParseDate(trigger); err != nil {
		return nil, err
	}

	/*err := c.db.AddEvent(p, newtrigger)

	if err := c.db.FindEvent(p); err != nil {
		return nil, err
	}*/

	/*timer, ok := c.triggers[trigger]
	if !ok {
		//c.db.AddEvent()
		newtimer, err := createTimer(p, c.finished)
		if err != nil {
			return nil, err
		}
		c.triggers[trigger] = newtimer
		return newtimer.events, nil
	}
	return timer.events, nil*/
	return nil, nil
}

func (c *calendarImpl) GetTriggers() []string {
	c.m.Lock()
	defer c.m.Unlock()
	count := len(c.triggers)
	list := make([]string, count)
	i := 0
	for _, t := range c.triggers {
		list[i] = t.id
		i++
	}
	return list
}

func (c *calendarImpl) DeleteTrigger(trigger string) bool {
	c.m.Lock()
	defer c.m.Unlock()
	t, ok := c.triggers[trigger]
	if ok {
		t.Stop()
		delete(c.triggers, trigger)
	}
	return ok
}

func (c *calendarImpl) GetEvents(trigger string) Events {
	c.m.Lock()
	defer c.m.Unlock()
	e, ok := c.triggers[trigger]
	if !ok {
		return nil
	}
	return e.events
}

func (c *calendarImpl) GetTriggerAlert(trigger string) (t time.Time, ok bool) {
	c.m.Lock()
	defer c.m.Unlock()
	e, ok := c.triggers[trigger]
	if !ok {
		return
	}
	return e.alerttime, true
}

func (c *calendarImpl) FindEvents(parameters SearchParameters) ([]Event, error) {
	c.m.Lock()
	defer c.m.Unlock()
	events := make([]Event, 0, 10)
	for _, t := range c.triggers {
		if checkSearchParameters(t.alerttime, parameters) {
			//todo optimize
			count := t.events.GetEventsCount()
			for i := 0; i < count; i++ {
				events = append(events, t.events.GetEvent(i))
			}
		}
	}
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
			//fmt.Println(t.Year(), int(t.Month()), t.Day(), p)
			if p.Year == t.Year() && p.Month == int(t.Month()) && p.Day == t.Day() {
				return true
			}
		}
	}
	return false
}
