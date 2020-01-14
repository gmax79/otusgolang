package pmetrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

// Agent - main object to send statistic to prometheus
type Agent struct {
	finished bool
}

// CreateMetricsAgent - create prometheus client
func CreateMetricsAgent() *Agent {
	var a Agent
	return &a
}

func cantCreateError(metric string) error {
	return errors.New("Can't create metric '" + metric + "', prometheus exporter closed")
}

// Shutdown - stop collect data
func (a *Agent) Shutdown() {
	a.finished = true
}

// RegisterCounterMetric - create standard counter metric
func (a *Agent) RegisterCounterMetric(name, descr string) (func(), error) {
	if a.finished {
		return func() {}, cantCreateError(name)
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
		if a.finished {
			return
		}
		c.Inc()
	}, nil
}

// RegisterGaugeMetric - create standard gauge matric
func (a *Agent) RegisterGaugeMetric(name, descr string) (func(float64), error) {
	if a.finished {
		return func(float64) {}, cantCreateError(name)
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
		if a.finished {
			return
		}
		c.Set(v)
	}, nil
}
