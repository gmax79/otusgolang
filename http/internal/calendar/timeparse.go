package calendar

import (
	"fmt"
	"regexp"
	"time"
)

type timeTriggerParser struct {
	parsed time.Time
}

var parseTimeTrigger *regexp.Regexp
var parseDateTrigger *regexp.Regexp

func init() {
	parseTimeTrigger = regexp.MustCompile("^ *[0-9]{2}:[0-9]{2}:[0-9]{2}")
	parseDateTrigger = regexp.MustCompile("^ *([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})")
}

func (tp *timeTriggerParser) Parse(timeTrigger string) error {
	if timeTrigger == "" {
		return fmt.Errorf("Time doesn't declared")
	}
	// variant date with time
	parts := parseDateTrigger.FindStringSubmatch(timeTrigger)
	if parts != nil {
		var err error
		tp.parsed, err = time.Parse(DateLayout, parts[0])
		return err
	}
	// variant with time only
	parts = parseTimeTrigger.FindStringSubmatch(timeTrigger)
	if parts != nil {
		var err error
		tp.parsed, err = time.Parse(TimeLayout, parts[0])
		return err
	}
	return fmt.Errorf("Time is in invalid format")
}

func (tp *timeTriggerParser) Normalize(custom time.Time) {
	tp.Parse(custom.String())
}

func (tp *timeTriggerParser) SetNormalizedNow() {
	tp.Normalize(time.Now())
}
