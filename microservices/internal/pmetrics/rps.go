package pmetrics

import (
	"time"
)

// CreateRPSCounter - create adapter, which working with time. Returns incremental function to calc result value
func CreateRPSCounter(setter func(float64)) func(int) {
	var count float64
	var seconds float64
	before := time.Now()
	return func(add int) {
		now := time.Now()
		delta := now.Sub(before).Seconds()
		before = now
		seconds = seconds + delta
		count = count + float64(add)
		if seconds > 0 {
			setter(count / seconds)
			return
		}
		setter(0)
	}
}
