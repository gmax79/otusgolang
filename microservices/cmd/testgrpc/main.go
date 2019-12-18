package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	tests "github.com/gmax79/otusgolang/microservices/internal/stests"
)

const host = "localhost:9090"

func s2date(stime string) simple.Date {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, stime)
	if err != nil {
		log.Fatal(err)
	}
	var d simple.Date
	d.Year = t.Year()
	d.Month = int(t.Month())
	d.Day = t.Day()
	d.Hour = t.Hour()
	d.Minute = t.Minute()
	d.Second = t.Second()
	return d
}

func assert(result string, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func assertCount(prefix string, count int, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prefix, ":", count, "events")
}

func main() {

	var err error
	cli, err := grpccon.CreateClient(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connnected at grpc host:", host)

	assert(cli.CreateEvent(s2date("2020-01-07 12:00:00"), "Exam in school"))
	assert(cli.DeleteEvent(s2date("2020-01-07 12:00:00"), "Exam in school"))
	assert(cli.CreateEvent(s2date("2020-01-09 13:00:00"), "Call Willy"))
	assert(cli.DeleteEvent(s2date("2020-01-09 15:00:00"), "Exam in school"))
	assert(cli.MoveEvent(s2date("2020-01-09 13:00:00"), "Exam in school", s2date("2020-01-09 15:00:00")))

	assert(cli.CreateEvent(s2date("2020-01-12 8:00:00"), "Pay credit"))
	assert(cli.CreateEvent(s2date("2020-01-14 10:00:00"), "Send pacel to Jack"))

	var count int
	count, err = cli.GetEventsForDay(9, 1, 2020)
	assertCount("At 2020-1-9", count, err)
	count, err = cli.GetEventsForWeek(1, 2020)
	assertCount("At week 2020-1", count, err)
	count, err = cli.GetEventsForMonth(1, 2020)
	assertCount("At month 2002-1", count, err)

	var now simple.Date
	now.SetNow()

	d1 := tests.DurationToSimpleDate(time.Second * 3)
	assert(cli.CreateEvent(d1, "Test since method #1"))
	d2 := tests.DurationToSimpleDate(time.Second * 6)
	assert(cli.CreateEvent(d2, "Test since method #2"))

	time.Sleep(time.Second * 7)

	events, err := cli.SinceEvents(now)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range events {
		fmt.Println(e)
	}

	cli.Close()
	fmt.Println("Connection at grpc host closed")
}
