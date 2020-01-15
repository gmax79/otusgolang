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

	log.Println("Sender init")

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

	exporter, errmetric := pmetrics.StartPrometheusExporter(config.Prometheus)
	if errmetric != nil {
		fmt.Println(errmetric)
	}

	agent := pmetrics.CreateMetricsAgent()
	counterfunc, errmetric := agent.RegisterCounterMetric("sender_messages_count", "Count messages sent by sender sevice")
	if errmetric != nil {
		fmt.Println("Can't register sender_messages_count metric", errmetric)
	}
	rpsfunc, errmetric := agent.RegisterRPSMetric("sender_messages_rps", "RPS of messages sent by sender service")
	if errmetric != nil {
		fmt.Println("Can't register sender_messages_rps metric", errmetric)
	}

	var rabbitConn *api.RmqConnection
	var datachan <-chan []byte
	connect := func() error {
		if rabbitConn != nil {
			rabbitConn.Close()
		}
		var connerr error
		for {
			rabbitConn, connerr = api.RabbitMQConnect(config.RabbitMQAddr())
			if connerr != nil {
				return connerr
			}
			if datachan, connerr = rabbitConn.Subscribe("calendar"); connerr == nil {
				break
			}
			rabbitConn.Close()
			if api.IsNoQueueError(connerr) {
				log.Println("Wait rabbitmq queue")
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}
		return connerr
	}
	if err = connect(); err != nil {
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Sender started")

loop:
	for {
		select {
		case msg, ok := <-datachan:
			if !ok || len(msg) == 0 {
				log.Println("sender:", "connection to rabbit mq lost, reconnecting")
				if err := connect(); err != nil {
					log.Println(err)
					break loop
				}
			}
			counterfunc()
			rpsfunc(1)
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
	if rabbitConn != nil {
		rabbitConn.Close()
	}
	agent.Shutdown()
	if err = exporter.Shutdown(); err != nil {
		log.Println(err)
	}
	log.Println("Sender stopped")
}
