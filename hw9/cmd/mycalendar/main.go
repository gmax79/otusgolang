package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// ServiceConfig - parameters to start service
type ServiceConfig struct {
	LogFile      string `json:"log_file"`
	VerboseLevel string `json:"log_level"`
	Encoding     string `json:"encoding"`
	Debugmode    int    `json:"debug"`
	ListenHTTP   string `json:"http_listen"`
}

func readServiceConfig() (*ServiceConfig, error) {
	config := &ServiceConfig{}
	configFile := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	data, err := ioutil.ReadFile(*configFile)
	if err == nil {
		err = json.Unmarshal(data, &config)
	}
	return config, err
}

func createLogger(s *ServiceConfig) (*zap.Logger, error) {
	zapconfig := zap.NewProductionConfig()
	zapconfig.DisableStacktrace = true
	if s.Debugmode != 0 {
		zapconfig.Development = true
		zapconfig.DisableStacktrace = false
	}
	if s.Encoding == "console" || s.Encoding == "json" {
		zapconfig.Encoding = s.Encoding
	} else if s.Encoding != "" {
		return nil, fmt.Errorf("Encoding type '%s' not supported", s.Encoding)
	}
	if s.LogFile != "" {
		zapconfig.OutputPaths = append(zapconfig.OutputPaths, s.LogFile)
	}
	switch s.VerboseLevel {
	case "debug":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warning":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("Verbose level '%s' not supported", s.VerboseLevel)
	}
	return zapconfig.Build()
}

// CalendarService - main struct of service
type CalendarService struct {
	logger *zap.Logger
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error: %v\nUse --help option to read usage information", err)
		}
	}()

	var config *ServiceConfig
	if config, err = readServiceConfig(); err != nil {

		return
	}
	var logger *zap.Logger
	if logger, err = createLogger(config); err != nil {
		return
	}

	s := &CalendarService{logger: logger}
	http.HandleFunc("/", s.httpRoot)
	http.HandleFunc("/hello", s.httpHello)
	logger.Info("Calendar service started")
	logger.Info("Go in browser at host ", zap.String("url", config.ListenHTTP))
	httperr := http.ListenAndServe(config.ListenHTTP, nil)
	if httperr != nil {
		logger.Error("error", zap.Error(httperr))
	}
}

func (s *CalendarService) logRequest(r *http.Request) {
	s.logger.Info("request", zap.String("url", r.URL.Path))
}

func (s *CalendarService) httpRoot(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	fmt.Fprint(w, "Calendar app")
}

func (s *CalendarService) httpHello(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	fmt.Fprint(w, "Hello !")
}
