package grpccon

import (
	"context"

	"github.com/gmax79/otusgolang/microservices/api/pbcalendar"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	"google.golang.org/grpc"
)

// Client - main object for grpc client for calendar service
type Client struct {
	cancel func()
	ctx    context.Context
	client pbcalendar.MyCalendarClient
}

func d2pb(s simple.Date) *pbcalendar.Date {
	var d pbcalendar.Date
	d.Year = int32(s.Year)
	d.Month = int32(s.Month)
	d.Day = int32(s.Day)
	d.Hour = int32(s.Hour)
	d.Minute = int32(s.Minute)
	d.Second = int32(s.Second)
	return &d
}

// CreateClient - create instance of connection to service
func CreateClient(host string) (*Client, error) {

	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &Client{}
	var cancelfunc func()
	c.ctx, cancelfunc = context.WithCancel(context.Background())
	c.cancel = func() {
		cancelfunc()
		clientCon.Close()
	}
	c.client = pbcalendar.NewMyCalendarClient(clientCon)
	return c, nil
}

// Close - close grpc connection
func (c *Client) Close() {
	c.cancel()
}

// CreateEvent - call grpc to create event
func (c *Client) CreateEvent(date simple.Date, info string) (string, error) {
	var e pbcalendar.CreateEventRequest
	e.Alerttime = d2pb(date)
	e.Information = info
	result, err := c.client.CreateEvent(c.ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// DeleteEvent - call grpc to delete event
func (c *Client) DeleteEvent(date simple.Date, info string) (string, error) {
	var e pbcalendar.DeleteEventRequest
	e.Alerttime = d2pb(date)
	e.Information = info
	result, err := c.client.DeleteEvent(c.ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// MoveEvent - call grpc to move event
func (c *Client) MoveEvent(date simple.Date, info string, newdate simple.Date) (string, error) {
	var e pbcalendar.MoveEventRequest
	e.Alerttime = d2pb(date)
	e.Information = info
	e.Newdate = d2pb(newdate)
	result, err := c.client.MoveEvent(c.ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// GetEventsForDay - grpc, calculate events for day
func (c *Client) GetEventsForDay(day, month, year int) (int, error) {
	var e pbcalendar.EventsForDayRequest
	e.Day = int32(day)
	e.Month = int32(month)
	e.Year = int32(year)
	result, err := c.client.EventsForDay(c.ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}

// GetEventsForWeek - grpc, calculate events for week
func (c *Client) GetEventsForWeek(week, year int) (int, error) {
	var e pbcalendar.EventsForWeekRequest
	e.Week = int32(week)
	e.Year = int32(year)
	result, err := c.client.EventsForWeek(c.ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}

// GetEventsForMonth - grpc, calculate events for month
func (c *Client) GetEventsForMonth(month, year int) (int, error) {
	var e pbcalendar.EventsForMonthRequest
	e.Month = int32(month)
	e.Year = int32(year)
	result, err := c.client.EventsForMonth(c.ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}
