package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type timeTriggerParser struct {
	parsed time.Time
}

var parseTimeTrigger *regexp.Regexp
var parseSmallTimeTrigger *regexp.Regexp
var parseDateTrigger *regexp.Regexp

func init() {
	parseTimeTrigger = regexp.MustCompile("^ *([0-9]{2}:[0-9]{2}:[0-9]{2})")
	parseSmallTimeTrigger = regexp.MustCompile("^ *([0-9]{2}:[0-9]{2})")
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
		clock := strings.Split(parts[1], ":")
		var err error
		tp.parsed, err = stringTime(clock)
		return err
	}
	parts = parseSmallTimeTrigger.FindStringSubmatch(timeTrigger)
	if parts != nil {
		clock := strings.Split(parts[1], ":")
		var err error
		tp.parsed, err = stringTime(clock)
		return err
	}
	return fmt.Errorf("Time is in invalid format")
}

func stringTime(parts []string) (stime time.Time, err error) {
	count := len(parts)
	var hours, minutes, seconds int
	if hours, err = strconv.Atoi(parts[0]); err != nil {
		return
	}
	if minutes, err = strconv.Atoi(parts[1]); err != nil {
		return
	}
	if count > 2 {
		if seconds, err = strconv.Atoi(parts[2]); err != nil {
			return
		}
	}
	return makeTime(hours, minutes, seconds), nil
}

func makeTime(hours, minutes, seconds int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, seconds, 0, time.UTC)
}

func (tp *timeTriggerParser) Normalize(custom time.Time) {
	tp.Parse(custom.String())
}

func (tp *timeTriggerParser) SetNormalizedNow() {
	tp.Normalize(time.Now())
}
