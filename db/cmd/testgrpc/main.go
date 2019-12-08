package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"../mycalendar/pbcalendar"
	"google.golang.org/grpc"
)

const host = "localhost:9090"

type client struct {
	cancel func()
	ctx    context.Context
	client pbcalendar.MyCalendarClient
}

func createClient() (*client, error) {

	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &client{}
	var cancelfunc func()
	c.ctx, cancelfunc = context.WithCancel(context.Background())
	c.cancel = func() {
		cancelfunc()
		clientCon.Close()
	}
	c.client = pbcalendar.NewMyCalendarClient(clientCon)
	return c, nil
}

func (c *client) Shutdown() {
	c.cancel()
}

func (c *client) CreateEvent(e pbcalendar.Event) {
	result, err := c.client.CreateEvent(c.ctx, &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) DeleteEvent(e pbcalendar.Event) {
	result, err := c.client.DeleteEvent(c.ctx, &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) MoveEvent(e *pbcalendar.Event, nd *pbcalendar.Date) {
	var old pbcalendar.Event
	old.Alerttime = nd
	old.Information = e.Information
	c.DeleteEvent(old)
	var ne pbcalendar.MoveEvent
	ne.Newdate = nd
	ne.Event = e
	result, err := c.client.MoveEvent(c.ctx, &ne)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) GetEventsForDay(e *pbcalendar.EventsForDay) {
	fmt.Println("GetEventsForDay", e)
	result, err := c.client.EventsForDay(c.ctx, e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get :", result.Count, "days")
}

func (c *client) GetEventsForWeek(e *pbcalendar.EventsForWeek) {
	fmt.Println("GetEventsForWeek", e)
	result, err := c.client.EventsForWeek(c.ctx, e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get :", result.Count, "days")
}

func s2date(stime string) *pbcalendar.Date {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, stime)
	var d pbcalendar.Date
	d.Year = int32(t.Year())
	d.Month = int32(t.Month())
	d.Day = int32(t.Day())
	d.Hour = int32(t.Hour())
	d.Minute = int32(t.Minute())
	d.Second = int32(t.Second())
	return &d
}

func main() {
	c, err := createClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connnected at grpc host:", host)

	var e pbcalendar.Event

	e.Alerttime = s2date("2020-01-07 12:00:00")
	e.Information = "Exam in school"

	c.CreateEvent(e)
	c.DeleteEvent(e)

	c.CreateEvent(e)
	var nd *pbcalendar.Date
	nd = s2date("2020-01-09 15:00:00")
	c.MoveEvent(&e, nd)

	e.Information = "Pay credit"
	e.Alerttime = s2date("2020-01-12 8:00:00")
	c.CreateEvent(e)

	e.Information = "Send pacel to Jack"
	e.Alerttime = s2date("2020-01-14 10:00:00")
	c.CreateEvent(e)

	var eday pbcalendar.EventsForDay
	eday.Year = 2020
	eday.Month = 1
	eday.Day = 9
	c.GetEventsForDay(&eday)

	var eweek pbcalendar.EventsForWeek
	eweek.Week = 1
	eweek.Year = 2020
	c.GetEventsForWeek(&eweek)

	eweek.Week = 2
	eweek.Year = 2020
	c.GetEventsForWeek(&eweek)

	c.Shutdown()
	fmt.Println("Connection at grpc host closed")
}
