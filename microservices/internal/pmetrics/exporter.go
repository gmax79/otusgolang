package pmetrics

import (
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/gracefully"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
)

// AttachMiddlewareHandler - wrapping handler by prometheus
func AttachMiddlewareHandler(service string, handler http.Handler) http.Handler {
	metricsConfig := metrics.Config{}
	mdlw := middleware.New(middleware.Config{
		Service:  service,
		Recorder: metrics.NewRecorder(metricsConfig),
	})
	return mdlw.Handler("", handler)
}

// Exporter - main object for prometheus exporter
type Exporter struct {
	server *gracefully.HTTPServer
}

// StartPrometheusExporter - starts exporter web server
func StartPrometheusExporter(host string) (*Exporter, error) {
	var e Exporter
	e.server = gracefully.CreateHTTPServer(host, promhttp.Handler())
	timetostart := time.Millisecond * 200
	err := e.server.ListenAndServe(timetostart)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetLastError - returns error within finished promethus exporter
func (e *Exporter) GetLastError() error {
	return e.server.GetLastError()
}

// Shutdown stop prometheus exporter http server
func (e *Exporter) Shutdown() error {
	return e.server.Shutdown(time.Second * 2)
}
