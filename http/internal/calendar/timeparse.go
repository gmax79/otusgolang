package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type timeParser struct {
	year, month, day     int
	hour, minute, second int
}

var parseTime *regexp.Regexp
var parseDate *regexp.Regexp
var days = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func init() {
	parseTime = regexp.MustCompile("(?:^|\\s)([0-9]{1,2})(:[0-9]{2})?(:[0-9]{2})?(?:$|\\s)")
	parseDate = regexp.MustCompile("(?:^|\\s)([0-9]{4})(-[0-9]{2})?(-[0-9]{2})?(?:$|\\s)")
}

func (tp *timeParser) Parse(timeTrigger string) error {
	if timeTrigger == "" {
		return fmt.Errorf("Time or date doesn't declared")
	}
	atoi := func(s string, skip0 bool) int {
		if s == "" {
			return 0
		}
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
		tp.hour = atoi(p[1], false)
		tp.minute = atoi(p[2], true)
		tp.second = atoi(p[3], true)
		if tp.hour < 0 || tp.hour > 23 || tp.minute < 0 || tp.minute > 59 || tp.second < 0 || tp.second > 59 {
			return fmt.Errorf("Time with invalid value")
		}
	}
	if d != nil {
		tp.year = atoi(d[1], false)
		tp.month = atoi(d[2], true)
		tp.day = atoi(d[3], true)
		if tp.year < 0 || tp.month < 1 || tp.month > 12 || tp.day < 1 {
			return fmt.Errorf("Date with invalid value")
		}
		d := days[tp.month]
		if (tp.year % 4) == 0 {
			d++
		}
		if tp.day > d {
			return fmt.Errorf("Date with invalid value")
		}
	}
	return nil
}

func (tp *timeParser) Value() time.Time {
	month := time.Month(tp.month)
	return time.Date(tp.year, month, tp.day, tp.hour, tp.minute, tp.second, 0, time.UTC)
}

func (tp *timeParser) Now() time.Time {
	p := timeParser{}
	t := time.Now()
	p.hour = t.Hour()
	p.minute = t.Minute()
	p.second = t.Second()
	p.year = t.Year()
	p.month = int(t.Month())
	p.day = t.Day()
	return p.Value()
}
