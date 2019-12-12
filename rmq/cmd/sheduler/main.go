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

	"github.com/gmax79/otusgolang/rmq/api"
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

func (s *shedulerConfig) PostgresAddr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s", s.PsqlUser, s.PsqlPassword, s.PsqlHost, s.PsqlDatabase)
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

	var db *dbMonitor
	if db, err = connectToDatabase(config.PostgresAddr()); err != nil {
		return
	}
	defer db.Close()
	if err = db.ReadEvents(); err != nil {
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
	ticker := time.NewTicker(time.Second * 10)

loop:
	for {
		select {
		case <-ticker.C:
			if err = db.ReadEvents(); err != nil {
				fmt.Println(err)
				continue
			}
			event, ok := db.GetNearestEvent()
			if ok {
				duration := event.Sub(time.Now())
				fmt.Println("Next event after", duration.String())
			}
		case <-stop:
			break loop
		}
	}
	ticker.Stop()
	fmt.Println("Sheduler stopped")
}

func sendEventToRabbit(conn *api.RmqConnection, event string) error {
	var m api.RmqMessage
	m.Event = event
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = conn.SendMessage("calendar", data)
	return err
}
