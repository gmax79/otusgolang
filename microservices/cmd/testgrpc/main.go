package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
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

func main() {

	var err error
	cli, err := grpccon.CreateClient(host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connnected at grpc host:", host)

	var result string
	result, err = cli.CreateEvent(s2date("2020-01-07 12:00:00"), "Exam in school")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	result, err = cli.DeleteEvent(s2date("2020-01-07 12:00:00"), "Exam in school")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	result, err = cli.CreateEvent(s2date("2020-01-09 13:00:00"), "Call Willy")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	result, err = cli.DeleteEvent(s2date("2020-01-09 15:00:00"), "Exam in school")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	result, err = cli.MoveEvent(s2date("2020-01-09 13:00:00"), "Exam in school", s2date("2020-01-09 15:00:00"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	result, err = cli.CreateEvent(s2date("2020-01-12 8:00:00"), "Pay credit")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	result, err = cli.CreateEvent(s2date("2020-01-14 10:00:00"), "Send pacel to Jack")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	var count int
	count, err = cli.GetEventsForDay(9, 1, 2020)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("At 2020-1-9 :", count, "events")

	count, err = cli.GetEventsForWeek(1, 2020)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("At week 2020-1 :", count, "events")

	count, err = cli.GetEventsForWeek(2, 2020)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("At week 2020-2 :", count, "events")

	count, err = cli.GetEventsForMonth(1, 2020)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("At month 2002-1 :", count, "events")

	cli.Close()
	fmt.Println("Connection at grpc host closed")
}
