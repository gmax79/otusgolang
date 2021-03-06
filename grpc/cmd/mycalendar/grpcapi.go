package main

import (
	"context"
	"fmt"
	"net"

	"github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar"
	"github.com/gmax79/otusgolang/grpc/internal/calendar"
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

func pbTimeToString(t *pbcalendar.Date) string {
	return fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d", t.Year, t.Month, t.Day, t.Hour, t.Minute, t.Second)
}

func (g *grpcCalendarAPI) CreateEvent(ctx context.Context, e *pbcalendar.Event) (*pbcalendar.Result, error) {
	g.logger.Info("grpc CreateEvent", zap.String("event", e.String()))
	t := e.Alerttime
	trigger := pbTimeToString(t)
	events, err := g.calen.AddTrigger(trigger)
	if err != nil {
		g.logger.Error("grpc CreateEvent", zap.String("error", err.Error()))
		return nil, err
	}
	events.AddEvent(&SimpleEvent{event: e.Information})

	var result pbcalendar.Result
	result.Status = fmt.Sprintf("Event at %s added", trigger)
	return &result, nil
}

func (g *grpcCalendarAPI) DeleteEvent(ctx context.Context, e *pbcalendar.Event) (*pbcalendar.Result, error) {
	g.logger.Info("grpc DeleteEvent", zap.String("event", e.String()))
	var result pbcalendar.Result
	t := e.Alerttime
	trigger := pbTimeToString(t)
	events := g.calen.GetEvents(trigger)
	if events == nil {
		return nil, fmt.Errorf("Trigger %s not found", trigger)
	}
	index := events.FindEvent(e.Information)
	if index == -1 {
		return nil, fmt.Errorf("Event %s in trigger %s not found", e.Information, trigger)
	}
	events.DeleteEvent(index)
	result.Status = fmt.Sprintf("Event %s at %s deleted", e.Information, trigger)
	return &result, nil
}

func (g *grpcCalendarAPI) MoveEvent(ctx context.Context, e *pbcalendar.MoveEvent) (*pbcalendar.Result, error) {
	g.logger.Info("grpc MoveEvent", zap.String("event", e.String()))
	t := e.Event.Alerttime
	trigger := pbTimeToString(t)
	events := g.calen.GetEvents(trigger)
	if events == nil {
		return nil, fmt.Errorf("Trigger %s not found", trigger)
	}
	index := events.FindEvent(e.Event.Information)
	if index == -1 {
		return nil, fmt.Errorf("Event %s in trigger %s not found", e.Event.Information, trigger)
	}
	newtime := pbTimeToString(e.Newdate)
	newevents, err := g.calen.AddTrigger(newtime)
	if err != nil {
		return nil, err
	}
	foundevent := events.GetEvent(index)
	events.DeleteEvent(index)
	newevents.AddEvent(foundevent)
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
