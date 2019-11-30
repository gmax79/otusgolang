package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (c *client) AddEvent(e pbcalendar.Event) (string, error) {
	result, err := c.client.AddEvent(c.ctx, &e)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

func s2date(stime string) *pbcalendar.Date {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, stime)
	var d pbcalendar.Date
	d.Year = int32(t.Year())
	d.Month = int32(t.Month())
	d.Day = int32(t.Day())
	d.Hour = int32(t.Hour())
	d.Minute = int32(t.Minute())
	d.Second = int32(t.Second())
	return &d
}

func main() {
	c, err := createClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connnected at grpc host:", host)

	var e pbcalendar.Event

	e.Alerttime = s2date("2020-01-07 12:00:00")
	e.Information = "Exam in school"

	var r string
	r, err = c.AddEvent(e)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	c.Shutdown()
	fmt.Println("Connection at grpc host closed")
}
