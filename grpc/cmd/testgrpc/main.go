package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar"
	"google.golang.org/grpc"
)

const host = "localhost:9090"

type client struct {
	cancel func()
	ctx    context.Context
	client pbcalendar.MyCalendarClient
}

func createClient() (*client, error) {

	clientCon, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &client{}
	var cancelfunc func()
	c.ctx, cancelfunc = context.WithCancel(context.Background())
	c.cancel = func() {
		cancelfunc()
		clientCon.Close()
	}
	c.client = pbcalendar.NewMyCalendarClient(clientCon)
	return c, nil
}

func (c *client) Shutdown() {
	c.cancel()
}

func (c *client) Test() {

}

func main() {
	c, err := createClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connnected at grpc host:", host)

	c.Test()

	c.Shutdown()
	fmt.Println("Connection at grpc host closed")
}
