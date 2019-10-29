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
	Debugmode    int    `json:"debug"`
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
	if s.Debugmode != 0 {
		zapconfig.Development = true
	}
	if s.Encoding == "console" || s.Encoding == "json" {
		zapconfig.Encoding = s.Encoding
	} else if s.Encoding != "" {
		return nil, fmt.Errorf("Encoding type '%s' not supported", s.Encoding)
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

	logger.Info("Info")
	logger.Warn("Warn")
	logger.Debug("Debug")
	logger.Error("Error")
}
