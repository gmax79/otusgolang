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
)

// ShedulerConfig - base parameters
type shedulerConfig struct {
	RmqlHost     string `json:"rabbitmq_host"`
	RmqlUser     string `json:"rabbitmq_user"`
	RmqPassword  string `json:"rabbitmq_password"`
	PsqlHost     string `json:"postgres_host"`
	PsqlUser     string `json:"postgres_user"`
	PsqlPassword string `json:"postgres_password"`
	PsqlDatabase string `json:"postgres_database"`
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

	rabbitConn, err := api.RabbitMQConnect(config.RabbitMQAddr())
	if err != nil {
		return
	}
	defer rabbitConn.Close()

	if err = rabbitConn.DeclareQueue("calendar"); err != nil {
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var m api.RmqMessage
	m.Event = "test"
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	err = rabbitConn.SendMessage("calendar", data)
	if err != nil {
		return
	}
loop:
	for {
		select {
		case <-stop:
			break loop
		}
	}
	fmt.Println("Sheduler stopped")
}
