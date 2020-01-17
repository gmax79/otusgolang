package pmetrics

import (
	"net/http"
	"time"
)

type rpsvalue struct {
	counter int
	setter  GaugeFunc
}

type rpsMetricsHandler struct {
	values map[string]rpsvalue
	agent  *Agent
}

// function to return wrapper handler from main handler
func (rh *rpsMetricsHandler) Attach(labels map[string]string, h http.Handler) http.Handler {
	events := make(chan string, 1)
	go func() {
		ticker := time.NewTicker(time.Second)
		before := time.Now()
		for {
			select {
			case k := <-events:
				if v, ok := rh.values[k]; ok {
					v.counter = v.counter + 1
					rh.values[k] = v
				} else {
					labels["handler"] = k
					f, err := rh.agent.RegisterGaugeMetric("http_requests_rps", "counts rps for url requests", labels)
					if err != nil {
						return
					}
					rh.values[k] = rpsvalue{counter: 1, setter: f}
				}
			case <-ticker.C:
				now := time.Now()
				delta := now.Sub(before).Seconds()
				before = now
				for k, v := range rh.values {
					rps := float64(v.counter) / delta
					v.setter(rps)
					v.counter = 0
					rh.values[k] = v
				}
			case <-rh.agent.stop:
				return
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events <- r.URL.Path
		h.ServeHTTP(w, r)
	})
}

func createRpsMetricsHandler(a *Agent) MetricsHandler {
	var h rpsMetricsHandler
	h.values = make(map[string]rpsvalue)
	h.agent = a
	return &h
}
