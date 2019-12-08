package main

import (
	"context"
	"fmt"
	"net"

	"../../internal/calendar"
	"./pbcalendar"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type grpcCalendarAPI struct {
	server    *grpc.Server
	logger    *zap.Logger
	lasterror error
	calen     calendar.Calendar
}

// createGrpc - service grpc interface
func createGrpc(calen calendar.Calendar, host string, zaplog *zap.Logger) (*grpcCalendarAPI, error) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	g := &grpcCalendarAPI{}
	g.server = grpc.NewServer()
	g.logger = zaplog
	g.calen = calen
	pbcalendar.RegisterMyCalendarServer(g.server, g)

	go func() {
		g.lasterror = g.server.Serve(listen)
	}()
	return g, nil
}

func (g *grpcCalendarAPI) Shutdown() {
	g.server.Stop()
}

func (g *grpcCalendarAPI) GetLastError() error {
	return g.lasterror
}

func pbDateToCalendarDate(t *pbcalendar.Date) calendar.Date {
	var d calendar.Date
	d.Day = int(t.GetDay())
	d.Month = int(t.GetMonth())
	d.Year = int(t.GetYear())
	d.Hour = int(t.GetHour())
	d.Minute = int(t.GetMinute())
	d.Second = int(t.GetSecond())
	return d
}

func (g *grpcCalendarAPI) CreateEvent(ctx context.Context, e *pbcalendar.Event) (*pbcalendar.Result, error) {
	g.logger.Info("grpc CreateEvent", zap.String("event", e.String()))

	t := e.Alerttime
	trigger := pbDateToCalendarDate(t)
	events, err := g.calen.AddTrigger(trigger)
	if err != nil {
		g.logger.Error("grpc CreateEvent", zap.String("error", err.Error()))
		return nil, err
	}
	events.AddEvent(calendar.Event(e.Information))

	var result pbcalendar.Result
	result.Status = fmt.Sprintf("Event at %s added", trigger)
	return &result, nil
}

func (g *grpcCalendarAPI) DeleteEvent(ctx context.Context, e *pbcalendar.Event) (*pbcalendar.Result, error) {
	g.logger.Info("grpc DeleteEvent", zap.String("event", e.String()))

	var result pbcalendar.Result
	t := e.Alerttime
	trigger := pbDateToCalendarDate(t)
	events, err := g.calen.GetEvents(trigger)
	if err == nil {
		return nil, err
	}
	err = events.DeleteEvent(calendar.Event(e.Information))
	if err != nil {
		return nil, err
	}
	result.Status = fmt.Sprintf("Event %s at %s deleted", e.Information, trigger)
	return &result, nil
}

func (g *grpcCalendarAPI) MoveEvent(ctx context.Context, e *pbcalendar.MoveEvent) (*pbcalendar.Result, error) {
	g.logger.Info("grpc MoveEvent", zap.String("event", e.String()))

	trigger := pbDateToCalendarDate(e.Event.Alerttime)
	newtime := pbDateToCalendarDate(e.Newdate)
	events, err := g.calen.GetEvents(trigger)
	if err != nil {
		return nil, err
	}
	event := (calendar.Event)(e.Event.Information)
	err = events.MoveEvent(event, newtime)
	if err != nil {
		return nil, err
	}

	var result pbcalendar.Result
	result.Status = fmt.Sprintf("Event %s moved from %s to %s", e.Event.Information, trigger, newtime)
	return &result, nil
}

func (g *grpcCalendarAPI) EventsForDay(ctx context.Context, e *pbcalendar.EventsForDay) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForDay", zap.String("event", e.String()))

	var sp calendar.SearchParameters
	sp.Day = int(e.Day)
	sp.Month = int(e.Month)
	sp.Year = int(e.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForMonth(ctx context.Context, e *pbcalendar.EventsForMonth) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForMonth", zap.String("event", e.String()))

	var sp calendar.SearchParameters
	sp.Month = int(e.Month)
	sp.Year = int(e.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForWeek(ctx context.Context, e *pbcalendar.EventsForWeek) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForWeek", zap.String("event", e.String()))

	var sp calendar.SearchParameters
	sp.Week = int(e.Week)
	sp.Year = int(e.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}
