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
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/grpccon"

	"github.com/gmax79/otusgolang/microservices/api"
)

// ShedulerConfig - base parameters
type shedulerConfig struct {
	RmqlHost    string `json:"rabbitmq_host"`
	RmqlUser    string `json:"rabbitmq_user"`
	RmqPassword string `json:"rabbitmq_password"`
	GrpcHost    string `json:"grpc_host"`
}

func (s *shedulerConfig) RabbitMQAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s", s.RmqlUser, s.RmqPassword, s.RmqlHost)
}

func (s *shedulerConfig) ApplicationAddr() string {
	return s.GrpcHost
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

	var con *grpccon.Client
	con, err = grpccon.CreateClient(config.ApplicationAddr())
	if err != nil {
		return
	}
	defer con.Close()

	//finishedEvents := make(chan string)
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
	ticker := time.NewTicker(time.Second * 2)
	fmt.Println("Sheduler started")

loop:
	for {
		select {
		case <-ticker.C:
			/*if err = db.ReadEvents(); err != nil {
				fmt.Println(err)
				continue
			}*/
		//case e := <-finishedEvents:
		//	sendEventToRabbit(rabbitConn, e)
		case <-stop:
			break loop
		}
	}
	ticker.Stop()
	fmt.Println("Sheduler stopped")
}

func sendEventToRabbit(conn *api.RmqConnection, event string) error {
	fmt.Println("Send to RabbitMQ:", event)
	var m api.RmqMessage
	m.Event = event
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = conn.SendMessage("calendar", data)
	return err
}
