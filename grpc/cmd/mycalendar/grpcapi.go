package main

import (
	"context"
	"net"

	"github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar"

	//"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/reflection"
	//"google.golang.org/grpc/status"
	//pbcalendar "./pbcalendar"
)

type grpcCalendarAPI struct {
	server    *grpc.Server
	lasterror error
}

// createGrpc - service grpc interface
func createGrpc(host string) (*grpcCalendarAPI, error) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}
	g := &grpcCalendarAPI{}
	g.server = grpc.NewServer()
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

func (g *grpcCalendarAPI) AddEvent(context.Context, *pbcalendar.Event) (*pbcalendar.Result, error) {
	var result pbcalendar.Result
	result.Status = "OK"
	return &result, nil
}
