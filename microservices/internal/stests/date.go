package stests // support for tests

import (
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/simple"
)

// DurationToTimeString - get time parameter, now + duration
func DurationToTimeString(d time.Duration) string {
	t := time.Now().Add(d)
	s := t.String()
	const layout = "2006-01-02 15:04:05"
	return s[:len(layout)]
}

// DurationToSimpleDate - calculate date from Now + time.Duration
func DurationToSimpleDate(d time.Duration) simple.Date {
	t := DurationToTimeString(d)
	var date simple.Date
	err := date.ParseDate(t)
	if err != nil {
		log.Fatal(err)
	}
	return date
}
