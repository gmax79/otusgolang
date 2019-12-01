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

	"github.com/gmax79/otusgolang/db/internal/calendar"
	nlog "github.com/gmax79/otusgolang/db/internal/log"
	"go.uber.org/zap"
)

// CalendarServiceConfig - base parameters
type CalendarServiceConfig struct {
	ListenHTTP string `json:"host"`
	GrpcHost   string `json:"grpc"`
}

// Check - check paraameters
func (cc *CalendarServiceConfig) Check() error {
	if cc.ListenHTTP == "" {
		return fmt.Errorf("Host address not declared")
	}
	if cc.GrpcHost == "" {
		return fmt.Errorf("Grpc host address not declared")
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

	calen := calendar.Create()
	server := createServer(calen, config.ListenHTTP, logger)
	logger.Info("Calendar service started")
	logger.Info("Caledar api ", zap.String("host", config.ListenHTTP))
	logger.Info("Calendar grpc api ", zap.String("host", config.GrpcHost))

	grpc, err := createGrpc(calen, config.GrpcHost, logger)
	if err != nil {
		logger.Error("Grpc", zap.String("error", err.Error()))
		return
	}

	server.ListenAndServe()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	grpc.Shutdown()
	server.Shutdown()
	if httperr := server.GetLastError(); httperr != nil {
		logger.Error("error", zap.Error(httperr))
	}
	logger.Info("Calendar service stopped")
}
