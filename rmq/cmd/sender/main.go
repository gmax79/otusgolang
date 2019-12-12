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

	"github.com/gmax79/otusgolang/rmq/api"
)

// SenderConfig - base parameters
type senderConfig struct {
	RmqlHost    string `json:"rabbitmq_host"`
	RmqlUser    string `json:"rabbitmq_user"`
	RmqPassword string `json:"rabbitmq_password"`
}

func (s *senderConfig) RabbitMQAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s", s.RmqlUser, s.RmqPassword, s.RmqlHost)
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}()
	configFile := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	var configJSON []byte
	if configJSON, err = ioutil.ReadFile(*configFile); err != nil {
		return
	}
	var config senderConfig
	if err = json.Unmarshal(configJSON, &config); err != nil {
		return
	}

	rabbitConn, err := api.RabbitMQConnect(config.RabbitMQAddr())
	if err != nil {
		return
	}
	defer rabbitConn.Close()

	datachan, err := rabbitConn.Subscribe("calendar")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Sender started")
loop:
	for {
		select {
		case msg, ok := <-datachan:
			if !ok {
				break loop
			}
			if len(msg) == 0 {
				fmt.Println("empty body from rmq ???")
				continue
			}
			mq := &api.RmqMessage{}
			if err := json.Unmarshal(msg, mq); err != nil {
				fmt.Printf("Got invalid blob: %v\n", err)
			} else {
				fmt.Println("Event from calendar:", mq.Event)
			}
		case <-stop:
			break loop
		}
	}
	fmt.Println("Sender stopped")
}
