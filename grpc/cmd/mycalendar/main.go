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

	nlog "github.com/gmax79/otusgolang/grpc/internal/log"
	"go.uber.org/zap"
)

type CalendarServiceConfig struct {
	ListenHTTP string `json:"host"`
}

func (cc *CalendarServiceConfig) Check() error {
	if cc.ListenHTTP == "" {
		return fmt.Errorf("Host address not declared")
	}
	return nil
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
	config := CalendarServiceConfig{}
	if err = json.Unmarshal(configJSON, &config); err != nil {
		return
	}
	if err = config.Check(); err != nil {
		return
	}

	var logger *zap.Logger
	if logger, err = nlog.CreateLogger(configJSON); err != nil {
		return
	}
	server := createServer(config.ListenHTTP, logger)
	logger.Info("Calendar service started")
	logger.Info("Go in browser at host ", zap.String("url", config.ListenHTTP))

	server.ListenAndServe()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	server.Shutdown()
	if httperr := server.GetLastError(); httperr != nil {
		logger.Error("error", zap.Error(httperr))
	}
	logger.Info("Calendar service stopped")
}
