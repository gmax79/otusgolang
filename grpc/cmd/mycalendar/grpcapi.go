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
	if t == nil {
		return nil, fmt.Errorf("nil time accepted")
	}
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

func (g *grpcCalendarAPI) DeleteEvent(context.Context, *pbcalendar.Date) (*pbcalendar.Result, error) {
	var result pbcalendar.Result
	return &result, nil
}

func (g *grpcCalendarAPI) MoveEvent(ctx context.Context, req *pbcalendar.MoveEvent) (*pbcalendar.Result, error) {
	var result pbcalendar.Result
	return &result, nil
}

func (g *grpcCalendarAPI) EventsForDay(ctx context.Context, req *pbcalendar.EventsForDay) (*pbcalendar.Count, error) {
	var c pbcalendar.Count
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForMonth(ctx context.Context, req *pbcalendar.EventsForMonth) (*pbcalendar.Count, error) {
	var c pbcalendar.Count
	return &c, nil
}

func (g *grpcCalendarAPI) EventsForWeek(ctx context.Context, req *pbcalendar.EventsForWeek) (*pbcalendar.Count, error) {
	var c pbcalendar.Count
	return &c, nil
}
