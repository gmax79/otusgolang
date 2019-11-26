package log

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// Parameters to create logger
type loggerConfig struct {
	LogFile      string `json:"log_file"`
	VerboseLevel string `json:"log_level"`
	LogEncoding  string `json:"log_encoding"`
	Debugmode    int    `json:"log_debug"`
}

// CreateLogger - create logger with parameters
func CreateLogger(jsondata []byte) (*zap.Logger, error) {
	c := &loggerConfig{}
	err := json.Unmarshal(jsondata, &c)
	if err != nil {
		return nil, err
	}
	zapconfig := zap.NewDevelopmentConfig()
	zapconfig.DisableStacktrace = true
	zapconfig.DisableCaller = true
	if c.Debugmode != 0 {
		zapconfig.Development = true
		zapconfig.DisableStacktrace = false
		zapconfig.DisableCaller = false
	}
	if c.LogEncoding == "console" || c.LogEncoding == "json" {
		zapconfig.Encoding = c.LogEncoding
	} else if c.LogEncoding != "" {
		return nil, fmt.Errorf("Encoding type '%s' not supported. Supported console, json", c.LogEncoding)
	}
	if c.LogFile != "" {
		zapconfig.OutputPaths = append(zapconfig.OutputPaths, c.LogFile)
	}
	switch c.VerboseLevel {
	case "debug":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warning":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("Verbose level '%s' not supported. Supported debug, info, warning, error", c.VerboseLevel)
	}
	return zapconfig.Build()
}
