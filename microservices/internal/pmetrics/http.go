package pmetrics

import (
	"net/http"
)

// MetricsHandler - wrapper for implementation metrics handler
type MetricsHandler interface {
	Attach(labels map[string]string, h http.Handler) http.Handler
}

type responseWriterProxy struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterProxy) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
