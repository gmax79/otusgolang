package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	tests "github.com/gmax79/otusgolang/microservices/internal/testshelpers"
)

const host = "http://localhost:8888"

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Testing caledar rabbit mq pipeline. Create nearby events")
	defer fmt.Println("Tests for Messages queue are finished")
	r1 := map[string]string{
		"time":  tests.DurationToTimeString(time.Second * 5),
		"event": "RabbitMQ #1",
	}
	assert(tests.Post(host, "create_event", r1, http.StatusOK))

	r2 := map[string]string{
		"time":  tests.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.1",
	}
	assert(tests.Post(host, "create_event", r2, http.StatusOK))

	r3 := map[string]string{
		"time":  tests.DurationToTimeString(time.Second * 10),
		"event": "RabbitMQ #2.2",
	}
	assert(tests.Post(host, "create_event", r3, http.StatusOK))
}
