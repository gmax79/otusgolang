package main

import (
	"context"
	"fmt"
	"net"

	"github.com/gmax79/otusgolang/microservices/api/pbcalendar"
	"github.com/gmax79/otusgolang/microservices/internal/calendar"
	"github.com/gmax79/otusgolang/microservices/internal/grpccon"
	"github.com/gmax79/otusgolang/microservices/internal/objects"

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

func (g *grpcCalendarAPI) CreateEvent(ctx context.Context, e *pbcalendar.CreateEventRequest) (*pbcalendar.CreateEventResponse, error) {
	g.logger.Info("grpc CreateEvent", zap.String("event", e.String()))

	t := e.Alerttime
	trigger := grpccon.ProtoToDate(t)
	events, err := g.calen.AddTrigger(trigger)
	if err != nil {
		g.logger.Error("grpc CreateEvent", zap.String("error", err.Error()))
		return nil, err
	}

	var newevent objects.Event
	newevent.Alerttime = trigger
	newevent.Information = e.Information
	events.AddEvent(newevent)

	var result pbcalendar.CreateEventResponse
	result.Status = fmt.Sprintf("Event at %s added", trigger.String())
	return &result, nil
}

func (g *grpcCalendarAPI) DeleteEvent(ctx context.Context, e *pbcalendar.DeleteEventRequest) (*pbcalendar.DeleteEventResponse, error) {
	g.logger.Info("grpc DeleteEvent", zap.String("event", e.String()))

	var result pbcalendar.DeleteEventResponse
	t := e.Alerttime
	trigger := grpccon.ProtoToDate(t)
	events, err := g.calen.GetEvents(trigger)
	if err != nil {
		return nil, err
	}

	var delevent objects.Event
	delevent.Alerttime = trigger
	delevent.Information = e.Information
	err = events.DeleteEvent(delevent)
	if err != nil {
		return nil, err
	}
	result.Status = fmt.Sprintf("Event %s at %s deleted", e.Information, trigger.String())
	return &result, nil
}

func (g *grpcCalendarAPI) MoveEvent(ctx context.Context, e *pbcalendar.MoveEventRequest) (*pbcalendar.MoveEventResponse, error) {
	g.logger.Info("grpc MoveEvent", zap.String("event", e.String()))

	trigger := grpccon.ProtoToDate(e.Alerttime)
	newtime := grpccon.ProtoToDate(e.Newdate)
	events, err := g.calen.GetEvents(trigger)
	if err != nil {
		return nil, err
	}

	var moveevent objects.Event
	moveevent.Alerttime = trigger
	moveevent.Information = e.Information
	err = events.MoveEvent(moveevent, newtime)
	if err != nil {
		return nil, err
	}

	var result pbcalendar.MoveEventResponse
	result.Status = fmt.Sprintf("Event %s moved from %s to %s", e.Information, trigger.String(), newtime.String())
	return &result, nil
}

func (g *grpcCalendarAPI) EventsForDay(ctx context.Context, e *pbcalendar.EventsForDayRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForDay", zap.String("event", e.String()))

	var sp objects.SearchParameters
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

func (g *grpcCalendarAPI) EventsForMonth(ctx context.Context, e *pbcalendar.EventsForMonthRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForMonth", zap.String("event", e.String()))

	var sp objects.SearchParameters
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

func (g *grpcCalendarAPI) EventsForWeek(ctx context.Context, e *pbcalendar.EventsForWeekRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForWeek", zap.String("event", e.String()))

	var sp objects.SearchParameters
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

func (g *grpcCalendarAPI) SinceEvents(ctx context.Context, e *pbcalendar.SinceEventsRequest) (*pbcalendar.SinceEventsResponse, error) {
	g.logger.Info("grpc SinceEventsResponse", zap.String("interval", e.String()))
	var events pbcalendar.SinceEventsResponse
	return &events, nil
}
