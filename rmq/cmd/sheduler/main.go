package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"../../api"
	"github.com/streadway/amqp"
)

// ShedulerConfig - base parameters
type shedulerConfig struct {
	RmqlHost    string `json:"rabbitmq_host"`
	RmqlUser    string `json:"rabbitmq_user"`
	RmqPassword string `json:"rabbitmq_password"`
}

func (s *shedulerConfig) RabbitMQAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s", s.RmqlUser, s.RmqPassword, s.RmqlHost)
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error: %v\nUse --help option to read usage information", err)
		}
	}()
	configFile := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	var configJSON []byte
	if configJSON, err = ioutil.ReadFile(*configFile); err != nil {
		return
	}
	var config shedulerConfig
	if err = json.Unmarshal(configJSON, &config); err != nil {
		return
	}
	var rabbitConn *amqp.Connection
	if rabbitConn, err = amqp.Dial(config.RabbitMQAddr()); err != nil {
		return
	}
	defer rabbitConn.Close()

	var rabbitChan *amqp.Channel
	if rabbitChan, err = rabbitConn.Channel(); err != nil {
		return
	}
	defer rabbitChan.Close()

	if _, err = rabbitChan.QueueDeclare(
		"calendar", // queue name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		return
	}

	rmqchan, err := rabbitChan.Consume(
		"calendar", // queue name
		"",         // consumer
		false,      // auto ask
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil)        // arguments
	if err != nil {
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var m api.RmqMessage
	m.Event = "test"
	rmqchan <- m

loop:
	for {
		select {
		case <-stop:
			break loop
		}
	}
	fmt.Println("Sheduler stopped")

}
