package calendar

import (
	"fmt"
	"strings"
	"time"
)

type triggerParser struct {
	parsed time.Time
}

func (t *triggerParser) Parse(trigger string) error {
	if trigger == "" {
		return fmt.Errorf("Time doesn't declared")
	}
	trigger = strings.TrimSpace(trigger)
	if len(trigger) < len(DateLayout) {
		return fmt.Errorf("Time is in invalid format")
	}

	trigger = trigger[:len(DateLayout)]
	var err error
	t.parsed, err = time.Parse(DateLayout, trigger)
	return err
}

func (t *triggerParser) Normalize(custom time.Time) {
	t.Parse(custom.String())
}

func (t *triggerParser) NormalizeNow() {
	t.Normalize(time.Now())
}
