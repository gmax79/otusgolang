package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gmax79/otusgolang/http/internal"
	"go.uber.org/zap"
)

type CalendarServiceConfig struct {
	ListenHTTP string `json:"host"`
}

// CalendarService - main struct of service
type CalendarService struct {
	logger *zap.Logger
	config *CalendarServiceConfig
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
	configJSON, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return
	}

	s := &CalendarService{config: &CalendarServiceConfig{}}
	if err = json.Unmarshal(configJSON, &s.config); err != nil {
		return
	}
	if s.logger, err = internal.CreateLogger(configJSON); err != nil {
		return
	}

	http.HandleFunc("/", s.httpRoot)
	http.HandleFunc("/hello", s.httpHello)

	s.logger.Info("Calendar service started")
	s.logger.Info("Go in browser at host ", zap.String("url", s.config.ListenHTTP))
	httperr := http.ListenAndServe(s.config.ListenHTTP, nil)
	if httperr != nil {
		s.logger.Error("error", zap.Error(httperr))
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
