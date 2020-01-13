package pmetrics

import "time"

// RPScounter - counts rps from incremental information
type RPScounter struct {
	last time.Time
}

// CreateRPSCounter - create adapter, which working with time. Returns incremental function to calc result value
func CreateRPSCounter(setter func(float64)) func(int) {
	var c RPScounter
	c.last = time.Now()

	return func(add int) {
		var rps float64
		setter(rps)
	}
}
