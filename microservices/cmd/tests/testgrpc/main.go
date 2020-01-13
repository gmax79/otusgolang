package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	tests "github.com/gmax79/otusgolang/microservices/internal/testshelpers"
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
	if result != "" {
		fmt.Println(result)
	}
}

func assertCount(prefix string, awaitcount, count int, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prefix, ":", count, "events")
	if awaitcount != count {
		log.Fatal(fmt.Errorf("Returned (%d) invalid count of events. Must %d", count, awaitcount))
	}
}

func main() {

	fmt.Println("Testing calendar grpc interface app")
	defer fmt.Println("Tests via grpc interface finished")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var err error
	cli, err := grpccon.CreateClient(ctx, host)
	assert("", err)

	createEvent := func(date, info string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return cli.CreateEvent(ctx, s2date(date), info)
	}

	deleteEvent := func(date, info string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return cli.DeleteEvent(ctx, s2date(date), info)
	}

	moveEvent := func(date, info, nextdate string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return cli.MoveEvent(ctx, s2date(date), info, s2date(nextdate))
	}

	fmt.Println("Connnecting at grpc host:", host)

	assert(createEvent("2020-04-07 12:00:00", "Exam in school"))
	assert(deleteEvent("2020-04-07 12:00:00", "Exam in school"))
	assert(createEvent("2020-04-09 13:00:00", "Call Willy"))
	assert(deleteEvent("2020-04-09 15:00:00", "Exam in school"))
	assert(moveEvent("2020-04-09 13:00:00", "Exam in school", "2020-04-09 15:00:00"))

	assert(createEvent("2020-04-12 8:00:00", "Pay credit"))
	assert(createEvent("2020-04-14 10:00:00", "Send pacel to Jack"))

	var count int
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	count, err = cli.GetEventsForDay(ctx, 9, 4, 2020)
	assertCount("At 2020-4-9", 1, count, err)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	count, err = cli.GetEventsForWeek(ctx, 16, 2020)
	assertCount("At week 2020-16", 1, count, err)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	count, err = cli.GetEventsForMonth(ctx, 4, 2020)
	assertCount("At month 2020-4", 3, count, err)

	var now simple.Date
	now.SetNow()

	d1 := tests.DurationToSimpleDate(time.Second * 3)
	assert(createEvent(d1.String(), "Test since method #1"))
	d2 := tests.DurationToSimpleDate(time.Second * 6)
	assert(createEvent(d2.String(), "Test since method #2"))

	time.Sleep(time.Second * 7)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	events, err := cli.SinceEvents(ctx, now)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range events {
		fmt.Println(e)
	}

	cli.Close()
}
