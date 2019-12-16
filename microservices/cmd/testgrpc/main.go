package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gmax79/otusgolang/microservices/api/pbcalendar"
	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
)

const host = "localhost:9090"

type client struct {
	cli *grpccon.Client
}

func createClient() (*client, error) {
	var err error
	var c client
	c.cli, err = grpccon.CreateClient(host)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *client) CreateEvent(time, info string) {

	result, err := c.cli.CreateEvent(s2date(time))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
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

func s2date(stime string) simple.Date {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, stime)
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
