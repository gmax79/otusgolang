package calendar

import (
	"fmt"
	"time"
)

// Calendar implementaion
type calendarImpl struct {
	triggers map[string]*timerimpl
	finished chan string
}

func createCalendar() Calendar {
	newcalendar := &calendarImpl{}
	newcalendar.triggers = make(map[string]*timerimpl)
	newcalendar.finished = make(chan string, 1)
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
	return newcalendar
}

func (c *calendarImpl) AddTrigger(trigger string) (Events, error) {
	timer, ok := c.triggers[trigger]
	if !ok {
		newtimer, err := createTimer(trigger, c.finished)
		if err != nil {
			return nil, err
		}
		c.triggers[trigger] = newtimer
		return newtimer.events, nil
	}
	return timer.events, nil
}

func (c *calendarImpl) GetTriggers() []string {
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
	t, ok := c.triggers[trigger]
	if ok {
		t.Stop()
		delete(c.triggers, trigger)
	}
	return ok
}

func (c *calendarImpl) GetEvents(trigger string) Events {
	e, ok := c.triggers[trigger]
	if !ok {
		return nil
	}
	return e.events
}

func (c *calendarImpl) GetTriggerAlert(trigger string) (t time.Time, ok bool) {
	e, ok := c.triggers[trigger]
	if !ok {
		return
	}
	return e.alerttime, true
}

func (c *calendarImpl) FindEvents(parameters SearchParameters) ([]Event, error) {

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
	return true
}
