package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"go.uber.org/zap"
)

// Settings - parameters to start program
type Settings struct {
	LogFile      string `json:"logfile"`
	Encoding     string `json:"encoding"`
	VerboseLevel string `json:"level"`
	Debugmode    bool   `json:"debug"`
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
	zapconfig.Development = s.Debugmode
	if s.Encoding != "console" || s.Encoding == "json" {
		zapconfig.Encoding = s.Encoding
	} else if s.Encoding != "" {
		fmt.Printf("Encoding type %s not supported./n", s.Encoding)
	}
	switch s.VerboseLevel {
	case "debug":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "error":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warning":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "info":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zapconfig.Build()
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Application can't start, error: %v", err)
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

	logger.Info("Info")
	logger.Warn("Warn")
	logger.Debug("Debug")
	logger.Error("Error")
}
