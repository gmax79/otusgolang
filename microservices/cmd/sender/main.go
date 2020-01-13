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

	"github.com/gmax79/otusgolang/microservices/internal/pmetrics"

	"github.com/gmax79/otusgolang/microservices/api"
)

// SenderConfig - base parameters
type senderConfig struct {
	RmqlHost    string `json:"rabbitmq_host"`
	RmqlUser    string `json:"rabbitmq_user"`
	RmqPassword string `json:"rabbitmq_password"`
	Prometheus  string `json:"prometheus_exporter"`
}

func (s *senderConfig) RabbitMQAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s", s.RmqlUser, s.RmqPassword, s.RmqlHost)
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("sender: %v\n", err)
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

	a, errmetric := pmetrics.CreateMetricsAgent(config.Prometheus)
	if errmetric != nil {
		fmt.Println(errmetric)
	}
	defer a.Shutdown()
	counterfunc, errmetric := a.RegisterCounterMetric("sender_messages_total_sent", "Count messages sent by sender sevice")
	if errmetric != nil {
		fmt.Println("Can't register sender_messages_total_sent metric", errmetric)
	}
	rpsfunc, errmetric := a.RegisterGaugeMetric("sender_messages_rps", "RPS of messages sent by sender service")
	if errmetric != nil {
		fmt.Println("Can't register sender_messages_rps metric", errmetric)
	}
	rpsadapter := pmetrics.CreateRPSCounter(rpsfunc)

	rmqHost := config.RabbitMQAddr()
	var rabbitConn *api.RmqConnection
	var datachan <-chan []byte
	for {
		rabbitConn, err = api.RabbitMQConnect(rmqHost)
		if err != nil {
			return
		}
		if datachan, err = rabbitConn.Subscribe("calendar"); err == nil {
			break
		}
		if api.IsNoQueueError(err) {
			rabbitConn.Close()
			log.Println("Wait rabbitmq queue")
			time.Sleep(time.Second * 3)
			continue
		}
		return
	}

	defer rabbitConn.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Sender started")
loop:
	for {
		select {
		case msg, ok := <-datachan:
			if !ok {
				break loop
			}
			counterfunc()
			rpsadapter(1)
			if len(msg) == 0 {
				log.Println("sender:", "empty body from rabbit mq ???")
				continue
			}
			mq := &api.RmqMessage{}
			if err := json.Unmarshal(msg, mq); err != nil {
				log.Printf("sender: Got invalid blob: %v\n", err)
			} else {
				log.Println("sender:", "Event from calendar:", mq.Event)
			}
		case <-stop:
			break loop
		}
	}
	log.Println("Sender stopped")
}
