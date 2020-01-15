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

// createGRPC - service grpc interface
func createGRPC(calen calendar.Calendar, host string, zaplog *zap.Logger) (*grpcCalendarAPI, error) {
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

func (g *grpcCalendarAPI) CreateEvent(ctx context.Context, req *pbcalendar.CreateEventRequest) (*pbcalendar.CreateEventResponse, error) {
	g.logger.Info("grpc CreateEvent", zap.String("event", req.String()))

	t := req.Alerttime
	trigger := grpccon.ProtoToDate(t)
	events, err := g.calen.AddTrigger(trigger)
	if err != nil {
		g.logger.Error("grpc CreateEvent", zap.String("error", err.Error()))
		return nil, err
	}

	var newevent objects.Event
	newevent.Alerttime = trigger
	newevent.Information = req.Information
	err = events.AddEvent(newevent)
	if err != nil {
		g.logger.Error("grpc CreateEvent", zap.String("error", err.Error()))
		return nil, err
	}

	var result pbcalendar.CreateEventResponse
	result.Status = fmt.Sprintf("Event at %s added", trigger.String())
	return &result, nil
}

func (g *grpcCalendarAPI) DeleteEvent(ctx context.Context, req *pbcalendar.DeleteEventRequest) (*pbcalendar.DeleteEventResponse, error) {
	g.logger.Info("grpc DeleteEvent", zap.String("event", req.String()))

	var result pbcalendar.DeleteEventResponse
	t := req.Alerttime
	trigger := grpccon.ProtoToDate(t)
	events, err := g.calen.GetEvents(trigger)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	var delevent objects.Event
	delevent.Alerttime = trigger
	delevent.Information = req.Information
	err = events.DeleteEvent(delevent)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}
	result.Status = fmt.Sprintf("Event %s at %s deleted", req.Information, trigger.String())
	return &result, nil
}

func (g *grpcCalendarAPI) MoveEvent(ctx context.Context, req *pbcalendar.MoveEventRequest) (*pbcalendar.MoveEventResponse, error) {
	g.logger.Info("grpc MoveEvent", zap.String("event", req.String()))

	trigger := grpccon.ProtoToDate(req.Alerttime)
	newtime := grpccon.ProtoToDate(req.Newdate)
	events, err := g.calen.GetEvents(trigger)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	var moveevent objects.Event
	moveevent.Alerttime = trigger
	moveevent.Information = req.Information
	err = events.MoveEvent(moveevent, newtime)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	var result pbcalendar.MoveEventResponse
	result.Status = fmt.Sprintf("Event %s moved from %s to %s", req.Information, trigger.String(), newtime.String())
	return &result, nil
}

func (g *grpcCalendarAPI) EventsForDay(ctx context.Context, req *pbcalendar.EventsForDayRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForDay", zap.String("event", req.String()))

	var sp objects.SearchParameters
	sp.Day = int(req.Day)
	sp.Month = int(req.Month)
	sp.Year = int(req.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForMonth(ctx context.Context, req *pbcalendar.EventsForMonthRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForMonth", zap.String("event", req.String()))

	var sp objects.SearchParameters
	sp.Month = int(req.Month)
	sp.Year = int(req.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForWeek(ctx context.Context, req *pbcalendar.EventsForWeekRequest) (*pbcalendar.Count, error) {
	g.logger.Info("grpc EventsForWeek", zap.String("event", req.String()))

	var sp objects.SearchParameters
	sp.Week = int(req.Week)
	sp.Year = int(req.Year)

	events, err := g.calen.FindEvents(sp)
	if err != nil {
		g.logger.Error("grpc", zap.Error(err))
		return nil, err
	}

	count := len(events)
	var c pbcalendar.Count
	c.Count = int32(count)
	return &c, nil
}

func (g *grpcCalendarAPI) SinceEvents(ctx context.Context, req *pbcalendar.SinceEventsRequest) (*pbcalendar.SinceEventsResponse, error) {
	g.logger.Info("grpc SinceEvents", zap.String("interval", req.String()))

	events, err := g.calen.SinceEvents(grpccon.ProtoToDate(req.From))
	if err != nil {
		g.logger.Error("sql", zap.Error(err))
		return nil, err
	}
	pbevents := make([]*pbcalendar.SinceEvent, len(events))
	for i, e := range events {
		pbevents[i] = &pbcalendar.SinceEvent{Information: e.Information, Alerttime: grpccon.DateToProto(e.Alerttime)}
	}
	var response pbcalendar.SinceEventsResponse
	response.Events = pbevents
	return &response, nil
}
