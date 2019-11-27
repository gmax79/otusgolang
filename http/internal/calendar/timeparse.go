package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type timeParser struct {
	year, month, day     int
	hour, minute, second int
	parsed               time.Time
}

var parseTimeTrigger *regexp.Regexp
var parseSmallTimeTrigger *regexp.Regexp
var parseDateTrigger *regexp.Regexp

var parseTime *regexp.Regexp
var parseDate *regexp.Regexp

func init() {
	parseTimeTrigger = regexp.MustCompile("^ *([0-9]{2}:[0-9]{2}:[0-9]{2})")
	parseSmallTimeTrigger = regexp.MustCompile("^ *([0-9]{2}:[0-9]{2})")
	parseDateTrigger = regexp.MustCompile("^ *([0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})")

	parseTime = regexp.MustCompile("(?:^|\\s)([0-9]{1,2})(:[0-9]{2})?(:[0-9]{2})?(?:$|\\s)")
	parseDate = regexp.MustCompile("(?:^|\\s)([0-9]{4})(-[0-9]{2})?(-[0-9]{2})?(?:$|\\s)")
}

func (tp *timeParser) Parse(timeTrigger string) error {
	if timeTrigger == "" {
		return fmt.Errorf("Time or date doesn't declared")
	}
	atoi := func(s string, skip0 bool) int {
		from := 1
		if !skip0 {
			from = 0
		}
		i, _ := strconv.Atoi(s[from:])
		return i
	}

	p := parseTime.FindStringSubmatch(timeTrigger)
	d := parseDate.FindStringSubmatch(timeTrigger)
	if p == nil && d == nil {
		return fmt.Errorf("Time and date in invalid format")
	}

	if p != nil {
		count := len(p) - 1 // skip index 0
		tp.hour = atoi(p[1], false)
		if count > 1 {
			tp.minute = atoi(p[2], true)
		}
		if count == 3 {
			tp.second = atoi(p[3], true)
		}
	}

	if d != nil {
		count := len(d) - 1 // skip index 0
		tp.year = atoi(d[1], false)
		if count > 1 {
			tp.month = atoi(p[2], true)
		}
		if count == 3 {
			tp.day = atoi(p[3], true)
		}
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

func (tp *timeParser) Normalize(custom time.Time) {
	tp.Parse(custom.String())
}

func (tp *timeParser) SetNormalizedNow() {
	tp.Normalize(time.Now())
}
