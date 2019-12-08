package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Date - calendar date with time
type date struct {
	d Date
}

var parseTime *regexp.Regexp
var parseDate *regexp.Regexp
var days = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func init() {
	parseTime = regexp.MustCompile("(?:^|\\s)([0-9]{1,2})(:[0-9]{2})?(:[0-9]{2})?(?:$|\\s)")
	parseDate = regexp.MustCompile("(?:^|\\s)([0-9]{4})(-[0-9]{2})?(-[0-9]{2})?(?:$|\\s)")
}

// ParseDate - parse string for date and time
func (tp *date) ParseDate(dateString string) error {
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

	ref := &tp.d
	p := parseTime.FindStringSubmatch(dateString)
	d := parseDate.FindStringSubmatch(dateString)
	if p == nil && d == nil {
		return fmt.Errorf("Time and date in invalid format")
	}

	if p != nil {
		ref.Hour = atoi(p[1], false)
		ref.Minute = atoi(p[2], true)
		ref.Second = atoi(p[3], true)
	}
	if d != nil {
		ref.Year = atoi(d[1], false)
		ref.Month = atoi(d[2], true)
		ref.Day = atoi(d[3], true)
	}
	return nil
}

// Valid - validate date and time
func (tp *date) Valid() error {
	ref := &tp.d
	if ref.Hour < 0 || ref.Hour > 23 || ref.Minute < 0 || ref.Minute > 59 || ref.Second < 0 || ref.Second > 59 {
		return fmt.Errorf("Time with invalid value")
	}
	if ref.Year < 0 || ref.Month < 1 || ref.Month > 12 || ref.Day < 1 {
		return fmt.Errorf("Date with invalid value")
	}
	d := days[ref.Month-1]
	if ref.Month == 2 && (ref.Year%4) == 0 {
		d++
	}
	if ref.Day > d {
		return fmt.Errorf("Date with invalid value")
	}
	return nil
}

// Value - return time.Time
func (tp *date) Value() time.Time {
	ref := &tp.d
	month := time.Month(ref.Month)
	return time.Date(ref.Year, month, ref.Day, ref.Hour, ref.Minute, ref.Second, 0, time.UTC)
}

func (tp *date) String() string {
	return tp.d.String()
}

// SetNow - set and return Now time
func (tp *date) SetNow() time.Time {
	t := time.Now()
	ref := &tp.d
	ref.Hour = t.Hour()
	ref.Minute = t.Minute()
	ref.Second = t.Second()
	ref.Year = t.Year()
	ref.Month = int(t.Month())
	ref.Day = t.Day()
	return tp.Value()
}
