package main

import (
	"context"
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

func (g *grpcCalendarAPI) AddEvent(ctx context.Context, e *pbcalendar.Event) (*pbcalendar.Result, error) {
	g.logger.Info("grpc AddEvent", zap.String("event", e.String()))
	var result pbcalendar.Result
	result.Status = "OK"
	return &result, nil
}
