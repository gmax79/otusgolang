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

// SenderConfig - base parameters
type senderConfig struct {
	RmqlHost    string `json:"rabbitmq_host"`
	RmqlUser    string `json:"rabbitmq_user"`
	RmqPassword string `json:"rabbitmq_password"`
}

func (s *senderConfig) RabbitMQAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s", s.RmqlUser, s.RmqPassword, s.RmqlHost)
}

type mqMessage struct {
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

loop:
	for {
		select {
		case msg := <-rmqchan:
			if len(msg.Body) == 0 {
				fmt.Println("empty body from rmq ???")
				msg.Ack(false)
				continue
			}
			mq := &api.RmqMessage{}
			fmt.Println(string(msg.Body))
			if err := json.Unmarshal(msg.Body, mq); err != nil {
				fmt.Printf("Got invalid blob: %v\n", err)
			} else {
				fmt.Println(mq)
			}
			msg.Ack(false)
		case <-stop:
			break loop
		}
	}
	fmt.Println("Sender stopped")
}
