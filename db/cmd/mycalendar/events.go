package main

import "fmt"

type SimpleEvent struct {
	event string
}

func (d *SimpleEvent) GetName() string {
	return d.event
}

func (d *SimpleEvent) Invoke() {
	fmt.Println("Event", d.event, "!!!")
}
