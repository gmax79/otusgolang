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

// Settings - parameters to start program
type Settings struct {
	LogFile      string `json:"log_file"`
	VerboseLevel string `json:"log_level"`
	Encoding     string `json:"encoding"`
	Debugmode    int    `json:"debug"`
	ListenHTTP   string `json:"http_listen"`
}

func readSettings() (*Settings, error) {
	settings := &Settings{}
	settingsFile := flag.String("settings", "settings.json", "path to settings file")
	flag.Parse()
	data, err := ioutil.ReadFile(*settingsFile)
	if err == nil {
		err = json.Unmarshal(data, &settings)
	}
	return settings, err
}

func createLogger(s *Settings) (*zap.Logger, error) {
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
			log.Fatalf("Error: %v", err)
		}
	}()

	var settings *Settings
	if settings, err = readSettings(); err != nil {

		return
	}
	var logger *zap.Logger
	if logger, err = createLogger(settings); err != nil {
		return
	}

	s := &CalendarService{logger: logger}
	http.HandleFunc("/", s.httpRoot)
	http.HandleFunc("/hello", s.httpHello)
	logger.Info("Calendar service started")
	logger.Info("Go in browser at host ", zap.String("url", settings.ListenHTTP))
	httperr := http.ListenAndServe(settings.ListenHTTP, nil)
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
