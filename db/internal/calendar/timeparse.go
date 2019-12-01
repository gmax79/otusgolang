package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Date - calendar date with time
type Date struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}

var parseTime *regexp.Regexp
var parseDate *regexp.Regexp
var days = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func init() {
	parseTime = regexp.MustCompile("(?:^|\\s)([0-9]{1,2})(:[0-9]{2})?(:[0-9]{2})?(?:$|\\s)")
	parseDate = regexp.MustCompile("(?:^|\\s)([0-9]{4})(-[0-9]{2})?(-[0-9]{2})?(?:$|\\s)")
}

// ParseDate - parse string for date and time
func (tp *Date) ParseDate(dateString string) error {
	if dateString == "" {
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

	p := parseTime.FindStringSubmatch(dateString)
	d := parseDate.FindStringSubmatch(dateString)
	if p == nil && d == nil {
		return fmt.Errorf("Time and date in invalid format")
	}

	if p != nil {
		tp.Hour = atoi(p[1], false)
		tp.Minute = atoi(p[2], true)
		tp.Second = atoi(p[3], true)
	}
	if d != nil {
		tp.Year = atoi(d[1], false)
		tp.Month = atoi(d[2], true)
		tp.Day = atoi(d[3], true)
	}
	return nil
}

// Valid - validate date
func (tp *Date) Valid() error {
	if tp.Hour < 0 || tp.Hour > 23 || tp.Minute < 0 || tp.Minute > 59 || tp.Second < 0 || tp.Second > 59 {
		return fmt.Errorf("Time with invalid value")
	}
	if tp.Year < 0 || tp.Month < 1 || tp.Month > 12 || tp.Day < 1 {
		return fmt.Errorf("Date with invalid value")
	}
	d := days[tp.Month-1]
	if tp.Month == 2 && (tp.Year%4) == 0 {
		d++
	}
	if tp.Day > d {
		return fmt.Errorf("Date with invalid value")
	}
	return nil
}

// Value - return time.Time
func (tp *Date) Value() time.Time {
	month := time.Month(tp.Month)
	return time.Date(tp.Year, month, tp.Day, tp.Hour, tp.Minute, tp.Second, 0, time.UTC)
}

func (tp *Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d %d:%02d:02d", tp.Year, tp.Month, tp.Day, tp.Hour, tp.Minute, tp.Second)
}

// SetNow - set and return Now time
func (tp *Date) SetNow() time.Time {
	t := time.Now()
	tp.Hour = t.Hour()
	tp.Minute = t.Minute()
	tp.Second = t.Second()
	tp.Year = t.Year()
	tp.Month = int(t.Month())
	tp.Day = t.Day()
	return tp.Value()
}
