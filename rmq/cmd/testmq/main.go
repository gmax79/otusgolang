package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/rmq/internal/calendar"
	tests "github.com/gmax79/otusgolang/rmq/internal/stests"
)

const host = "http://localhost:8080"

func main() {
	fmt.Println("Testing rabbit mq pipeline. Create nearby events")
	r1 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 5),
		"event": "RabbitMQ #1",
	}
	tests.Post(host, "create_event", r1, http.StatusOK)

	r2 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.1",
	}
	tests.Post(host, "create_event", r2, http.StatusOK)

	r3 := map[string]string{
		"time":  calendar.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.2",
	}
	tests.Post(host, "create_event", r3, http.StatusOK)
}
