package calendar

// Calendar implementaion
type calendarImpl struct {
	triggers []CalendarTrigger
}

func (c *calendarImpl) AddTrigger(t EventTrigger, ct CalendarTrigger) {
	c.triggers = append(c.triggers, ct)
	t.Start(ct.Invoke)
}

func (c *calendarImpl) GetTriggersCount() int {
	return len(c.triggers)
}

func (c *calendarImpl) DeleteTrigger(index int) bool {
	if index >= 0 && index < len(c.triggers) {
		c.triggers = append(c.triggers[:index], c.triggers[index+1:]...)
		return true
	}
	return false
}

func (c *calendarImpl) GetTrigger(index int) CalendarTrigger {
	if index >= 0 && index < len(c.triggers) {
		return c.triggers[index]
	}
	return nil
}

type triggerImpl struct {
	events []Event
}

func (t *triggerImpl) Invoke() {
	for _, e := range t.events {
		e.Invoke()
	}
}

func (t *triggerImpl) AddEvent(e Event) bool {
	if e == nil {
		return false
	}
	t.events = append(t.events, e)
	return true
}

func (t *triggerImpl) GetEventsCount() int {
	return len(t.events)
}

func (t *triggerImpl) DeleteEvent(index int) bool {
	if index >= 0 && index < len(t.events) {
		t.events = append(t.events[:index], t.events[index+1:]...)
		return true
	}
	return false
}

func (t *triggerImpl) GetEvent(index int) Event {
	if index >= 0 && index < len(t.events) {
		return t.events[index]
	}
	return nil
}

func (t *triggerImpl) ReplaceEvent(index int, e Event) bool {
	if e == nil {
		return false
	}
	if t.DeleteEvent(index) {
		t.AddEvent(e)
		return true
	}
	return false
}
