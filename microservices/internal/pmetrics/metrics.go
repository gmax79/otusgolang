package pmetrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Agent - main object to send statistic to prometheus
type Agent struct {
	lastError error
	client    *http.Server
	finished  bool
}

// CreateMetricsAgent - create prometheus client
func CreateMetricsAgent(listen string) (*Agent, error) {
	var a Agent
	if listen == "" {
		a.finished = true
		return &a, errors.New("Metrics will not collected, listen host for exporter not declared")
	}
	a.client = &http.Server{Addr: listen, Handler: promhttp.Handler()}
	wait := make(chan struct{})
	go func() {
		close(wait)
		a.lastError = a.client.ListenAndServe()
		a.finished = true
	}()
	<-wait
	time.Sleep(time.Millisecond * 100)
	if a.lastError != nil {
		var dummy Agent
		dummy.finished = true
		return &dummy, errors.New("Metrics will not collected. " + a.lastError.Error())
	}
	return &a, nil
}

// Shutdown - stopping prometheus client
func (a *Agent) Shutdown() {
	a.finished = true
	if a.client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	a.client.Shutdown(ctx)
}

// RegisterCounterMetric -
func (a *Agent) RegisterCounterMetric(name, descr string) (func(), error) {
	if a.finished {
		return func() {}, nil // return dummy function (not collecting)
	}
	var opt prometheus.CounterOpts
	opt.Help = descr
	opt.Name = name
	c := prometheus.NewCounter(opt)
	err := prometheus.Register(c)
	if err != nil {
		return nil, err
	}
	return func() {
		c.Inc()
	}, nil
}

// RegisterGaugeMetric -
func (a *Agent) RegisterGaugeMetric(name, descr string) (func(float64), error) {
	if a.finished {
		return func(float64) {}, nil
	}
	var opt prometheus.GaugeOpts
	opt.Help = descr
	opt.Name = name
	c := prometheus.NewGauge(opt)
	err := prometheus.Register(c)
	if err != nil {
		return nil, err
	}
	return func(v float64) {
		c.Set(v)
	}, nil
}
