package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/api/pbcalendar"
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

func (c *client) CreateEvent(e pbcalendar.CreateEventRequest) {
	result, err := c.client.CreateEvent(c.ctx, &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) DeleteEvent(e pbcalendar.DeleteEventRequest) {
	result, err := c.client.DeleteEvent(c.ctx, &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) MoveEvent(e *pbcalendar.MoveEventRequest) {
	var old pbcalendar.DeleteEventRequest
	old.Alerttime = e.Alerttime
	old.Information = e.Information
	c.DeleteEvent(old)
	result, err := c.client.MoveEvent(c.ctx, e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Status)
}

func (c *client) GetEventsForDay(e *pbcalendar.EventsForDayRequest) {
	fmt.Println("GetEventsForDay", e)
	result, err := c.client.EventsForDay(c.ctx, e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get :", result.Count, "days")
}

func (c *client) GetEventsForWeek(e *pbcalendar.EventsForWeekRequest) {
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

	var ce pbcalendar.CreateEventRequest

	ce.Alerttime = s2date("2020-01-07 12:00:00")
	ce.Information = "Exam in school"

	c.CreateEvent(ce)

	var de pbcalendar.DeleteEventRequest
	de.Alerttime = ce.Alerttime
	de.Information = ce.Information
	c.DeleteEvent(de)

	c.CreateEvent(ce)
	var nd *pbcalendar.Date
	nd = s2date("2020-01-09 15:00:00")

	var me pbcalendar.MoveEventRequest
	me.Alerttime = ce.Alerttime
	me.Information = ce.Information
	me.Newdate = nd
	c.MoveEvent(&me)

	ce.Information = "Pay credit"
	ce.Alerttime = s2date("2020-01-12 8:00:00")
	c.CreateEvent(ce)

	ce.Information = "Send pacel to Jack"
	ce.Alerttime = s2date("2020-01-14 10:00:00")
	c.CreateEvent(ce)

	var eday pbcalendar.EventsForDayRequest
	eday.Year = 2020
	eday.Month = 1
	eday.Day = 9
	c.GetEventsForDay(&eday)

	var eweek pbcalendar.EventsForWeekRequest
	eweek.Week = 1
	eweek.Year = 2020
	c.GetEventsForWeek(&eweek)

	eweek.Week = 2
	eweek.Year = 2020
	c.GetEventsForWeek(&eweek)

	c.Shutdown()
	fmt.Println("Connection at grpc host closed")
}
