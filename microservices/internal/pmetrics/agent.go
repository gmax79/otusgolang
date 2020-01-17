package pmetrics

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Agent - main object to send statistic to prometheus
type Agent struct {
	finished bool
	stop     chan struct{}
}

// CounterFunc - function to increment counter prometheus metric
type CounterFunc func()

// GaugeFunc - function to set value of gauge prometheus metric
type GaugeFunc func(float64)

// CreateMetricsAgent - create prometheus client
func CreateMetricsAgent() *Agent {
	var a Agent
	a.stop = make(chan struct{})
	return &a
}

func cantCreateError(metric string) error {
	return errors.New("Can't create metric '" + metric + "', prometheus exporter closed")
}

// Shutdown - stop collect data
func (a *Agent) Shutdown() {
	a.finished = true
	close(a.stop)
}

// RegisterCounterMetric - create standard counter metric
func (a *Agent) RegisterCounterMetric(name, descr string) (CounterFunc, error) {
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
func (a *Agent) RegisterGaugeMetric(name, descr string, labels map[string]string) (GaugeFunc, error) {
	if a.finished {
		return func(float64) {}, cantCreateError(name)
	}
	var opt prometheus.GaugeOpts

	opt.Help = descr
	opt.Name = name
	opt.ConstLabels = labels
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

// RegisterRPSMetric - create metric to autocalculate value per second from simple set value
func (a *Agent) RegisterRPSMetric(name, descr string) (GaugeFunc, error) {
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

	valchan := make(chan float64, 1)
	var valacc float64
	go func() {
		ticker := time.NewTicker(time.Second)
		before := time.Now()
		for {
			select {
			case <-a.stop:
				return
			case <-ticker.C:
				now := time.Now()
				delta := now.Sub(before).Seconds()
				before = now
				c.Set(valacc / delta)
				valacc = 0
			case v := <-valchan:
				valacc += v
			}
		}
	}()

	return func(v float64) {
		if a.finished {
			return
		}
		valchan <- v
	}, nil
}

// CreateReturnCodesMetricsHandler - create handler to calculate return codes from http request
func (a *Agent) CreateReturnCodesMetricsHandler() MetricsHandler {
	return createReturnCodesMetricsHandler(a)
}

// CreateRpsMetricsHandler - create handler to calculate rps of http requests
func (a *Agent) CreateRpsMetricsHandler() MetricsHandler {
	return createRpsMetricsHandler(a)
}
