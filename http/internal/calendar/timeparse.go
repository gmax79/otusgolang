package calendar

import (
	"fmt"
	"strings"
	"time"
)

type triggerParser struct {
	parsed time.Time
}

func (t *triggerParser) Parse(trigger, layout string) error {
	if trigger == "" {
		return fmt.Errorf("Error: trigger is empty")
	}
	if layout == "" {
		return fmt.Errorf("Error: trigger's layout is empty")
	}
	trigger = strings.TrimSpace(trigger)
	trigger = trigger[:len(layout)]
	var err error
	t.parsed, err = time.Parse(layout, trigger)
	return err
}
