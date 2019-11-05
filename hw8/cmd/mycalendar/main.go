package main

import (
	"fmt"
	"time"

	"github.com/gmax79/otusgolang/hw8/internal/calendar"
)

// DateTimeEventTrigger - trigger via date and time
type DateTimeEventTrigger struct {
	delay time.Duration
}

// Start - function to staring the trigger
func (t *DateTimeEventTrigger) Start(handler func()) {
	go func() {
		trigger := time.NewTimer(t.delay)
		<-trigger.C
		handler()
	}()
}

// Action - object with some action
type Action struct {
}

// Invoke - function, which call by calendar
func (a *Action) Invoke() {
	fmt.Println("Action !!!")
}

// Action2 - another object with action
type Action2 struct {
}

// Invoke - action happens
func (a *Action2) Invoke() {
	fmt.Println("Action 2 !!!")
}

func main() {
	c := calendar.CreateCalendar()

	t := calendar.CreateCalendarTrigger()
	t.AddEvent(&Action{})
	timeTrigger := &DateTimeEventTrigger{delay: time.Second * 4}
	c.AddTrigger(timeTrigger, t)

	t2 := calendar.CreateCalendarTrigger()
	t2.AddEvent(&Action2{})
	timeTrigger2 := &DateTimeEventTrigger{delay: time.Second * 1}
	c.AddTrigger(timeTrigger2, t2)

	time.Sleep(time.Second * 5)
}
