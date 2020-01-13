package grpccon

import (
	"context"

	"github.com/gmax79/otusgolang/microservices/api/pbcalendar"
	"github.com/gmax79/otusgolang/microservices/internal/objects"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	"google.golang.org/grpc"
)

// Client - main object for grpc client for calendar service
type Client struct {
	cancel func()
	client pbcalendar.MyCalendarClient
}

// DateToProto - convert calendar date to protobuf version
func DateToProto(s simple.Date) *pbcalendar.Date {
	var d pbcalendar.Date
	d.Year = int32(s.Year)
	d.Month = int32(s.Month)
	d.Day = int32(s.Day)
	d.Hour = int32(s.Hour)
	d.Minute = int32(s.Minute)
	d.Second = int32(s.Second)
	return &d
}

// ProtoToDate - convert from protobuf into calendar date
func ProtoToDate(p *pbcalendar.Date) simple.Date {
	var d simple.Date
	d.Year = int(p.Year)
	d.Month = int(p.Month)
	d.Day = int(p.Day)
	d.Hour = int(p.Hour)
	d.Minute = int(p.Minute)
	d.Second = int(p.Second)
	return d
}

// CreateClient - create instance of connection to service
func CreateClient(ctx context.Context, host string) (*Client, error) {

	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &Client{}
	var cancelfunc func()
	_, cancelfunc = context.WithCancel(ctx)
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
func (c *Client) CreateEvent(ctx context.Context, date simple.Date, info string) (string, error) {
	var e pbcalendar.CreateEventRequest
	e.Alerttime = DateToProto(date)
	e.Information = info
	result, err := c.client.CreateEvent(ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// DeleteEvent - call grpc to delete event
func (c *Client) DeleteEvent(ctx context.Context, date simple.Date, info string) (string, error) {
	var e pbcalendar.DeleteEventRequest
	e.Alerttime = DateToProto(date)
	e.Information = info
	result, err := c.client.DeleteEvent(ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// MoveEvent - call grpc to move event
func (c *Client) MoveEvent(ctx context.Context, date simple.Date, info string, newdate simple.Date) (string, error) {
	var e pbcalendar.MoveEventRequest
	e.Alerttime = DateToProto(date)
	e.Information = info
	e.Newdate = DateToProto(newdate)
	result, err := c.client.MoveEvent(ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// GetEventsForDay - grpc, calculate events for day
func (c *Client) GetEventsForDay(ctx context.Context, day, month, year int) (int, error) {
	var e pbcalendar.EventsForDayRequest
	e.Day = int32(day)
	e.Month = int32(month)
	e.Year = int32(year)
	result, err := c.client.EventsForDay(ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}

// GetEventsForWeek - grpc, calculate events for week
func (c *Client) GetEventsForWeek(ctx context.Context, week, year int) (int, error) {
	var e pbcalendar.EventsForWeekRequest
	e.Week = int32(week)
	e.Year = int32(year)
	result, err := c.client.EventsForWeek(ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}

// GetEventsForMonth - grpc, calculate events for month
func (c *Client) GetEventsForMonth(ctx context.Context, month, year int) (int, error) {
	var e pbcalendar.EventsForMonthRequest
	e.Month = int32(month)
	e.Year = int32(year)
	result, err := c.client.EventsForMonth(ctx, &e)
	if err != nil {
		return 0, err
	}
	return int(result.Count), nil
}

// SinceEvents - return events in next interval in seconds
func (c *Client) SinceEvents(ctx context.Context, from simple.Date) ([]objects.Event, error) {
	var r pbcalendar.SinceEventsRequest
	r.From = DateToProto(from)
	resp, err := c.client.SinceEvents(ctx, &r)
	if err != nil {
		return nil, err
	}
	count := len(resp.Events)
	events := make([]objects.Event, count)
	for i, e := range resp.Events {
		events[i].Alerttime = ProtoToDate(e.Alerttime)
		events[i].Information = e.Information
	}
	return events, nil
}
