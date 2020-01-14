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

	"github.com/gmax79/otusgolang/microservices/internal/calendar"
	nlog "github.com/gmax79/otusgolang/microservices/internal/log"
	"github.com/gmax79/otusgolang/microservices/internal/pmetrics"
	"go.uber.org/zap"
)

// CalendarServiceConfig - base parameters
type CalendarServiceConfig struct {
	ListenHTTP   string `json:"host"`
	GRPCHost     string `json:"grpc"`
	PsqlHost     string `json:"postgres_host"`
	PsqlUser     string `json:"postgres_user"`
	PsqlPassword string `json:"postgres_password"`
	PsqlDatabase string `json:"postgres_database"`
	PromExporter string `json:"prometheus_exporter"`
}

// Check - check paraameters
func (cc *CalendarServiceConfig) Check() error {
	if cc.ListenHTTP == "" {
		return fmt.Errorf("Host address not declared")
	}
	if cc.GRPCHost == "" {
		return fmt.Errorf("GRPC host address not declared")
	}
	if cc.PsqlHost == "" || cc.PsqlUser == "" || cc.PsqlPassword == "" || cc.PsqlDatabase == "" {
		return fmt.Errorf("Postgres not fully configurated")
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
	logger.Info("Calendar service started")
	logger.Info("Caledar api ", zap.String("host", config.ListenHTTP))
	logger.Info("Calendar grpc api ", zap.String("host", config.GRPCHost))
	logger.Info("Calendar db ", zap.String("host", config.PsqlHost), zap.String("database", config.PsqlDatabase))
	if config.PromExporter == "" {
		logger.Warn("Calendar prometheus exporter not configured. Service will not monitored")
	} else {
		logger.Info("Calendar prometheus exporter ", zap.String("host", config.PromExporter))
	}

	connection := fmt.Sprintf("postgresql://%s:%s@%s/%s", config.PsqlUser, config.PsqlPassword, config.PsqlHost, config.PsqlDatabase)
	calen, dberr := calendar.Create(connection)
	if dberr != nil {
		logger.Error("Postgres", zap.String("error", dberr.Error()))
		return
	}
	server := createServer(calen, config.ListenHTTP, logger)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("HTTP", zap.String("error", err.Error()))
		return
	}

	grpc, err := createGRPC(calen, config.GRPCHost, logger)
	if err != nil {
		logger.Error("GRPC", zap.String("error", err.Error()))
		return
	}

	var exporter *pmetrics.Exporter
	if config.PromExporter != "" {
		exporter, err = pmetrics.StartPrometheusExporter(config.PromExporter)
		if err != nil {
			logger.Warn("Monitoring", zap.Error(err))
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	grpc.Shutdown()
	server.Shutdown()
	if exporter != nil {
		exporter.Shutdown()
	}
	if experr := exporter.GetLastError(); experr != nil {
		logger.Error("error", zap.Error(experr))
	}
	if httperr := server.GetLastError(); httperr != nil {
		logger.Error("error", zap.Error(httperr))
	}
	logger.Info("Calendar service stopped")
}
