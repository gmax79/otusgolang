package simple

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var parseTime *regexp.Regexp
var parseDate *regexp.Regexp
var days = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func init() {
	parseTime = regexp.MustCompile("(?:^|\\s)([0-9]{1,2})(:[0-9]{2})?(:[0-9]{2})?(?:$|\\s)")
	parseDate = regexp.MustCompile("(?:^|\\s)([0-9]{4})(-[0-9]{2})?(-[0-9]{2})?(?:$|\\s)")
}

// Date -  date with time without timezone
type Date struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}

func (da *Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d %d:%02d:%02d", da.Year, da.Month, da.Day, da.Hour, da.Minute, da.Second)
}

// ParseDate - parse string for date and time
func (da *Date) ParseDate(date string) error {
	if date == "" {
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

	p := parseTime.FindStringSubmatch(date)
	d := parseDate.FindStringSubmatch(date)
	if p == nil && d == nil {
		return fmt.Errorf("Time and date in invalid format")
	}

	if p != nil {
		da.Hour = atoi(p[1], false)
		da.Minute = atoi(p[2], true)
		da.Second = atoi(p[3], true)
	}
	if d != nil {
		da.Year = atoi(d[1], false)
		da.Month = atoi(d[2], true)
		da.Day = atoi(d[3], true)
	}
	return nil
}

// Valid - validate date and time
func (da *Date) Valid() error {
	if da.Hour < 0 || da.Hour > 23 || da.Minute < 0 || da.Minute > 59 || da.Second < 0 || da.Second > 59 {
		return fmt.Errorf("Time with invalid value")
	}
	if da.Year < 0 || da.Month < 1 || da.Month > 12 || da.Day < 1 {
		return fmt.Errorf("Date with invalid value")
	}
	d := days[da.Month-1]
	if da.Month == 2 && (da.Year%4) == 0 {
		d++
	}
	if da.Day > d {
		return fmt.Errorf("Date with invalid value")
	}
	return nil
}

// Value - return time.Time
func (da *Date) Value() time.Time {
	month := time.Month(da.Month)
	return time.Date(da.Year, month, da.Day, da.Hour, da.Minute, da.Second, 0, time.UTC)
}

// SetNow - set and return Now time
func (da *Date) SetNow() {
	t := time.Now()
	da.Hour = t.Hour()
	da.Minute = t.Minute()
	da.Second = t.Second()
	da.Year = t.Year()
	da.Month = int(t.Month())
	da.Day = t.Day()
}

// ParseValidDate - create calendar date from string and validate
func ParseValidDate(trigger string) (Date, error) {
	var err error
	var zero Date
	var d Date
	if err = d.ParseDate(trigger); err != nil {
		return zero, err
	}
	if err = d.Valid(); err != nil {
		return zero, err
	}
	return d, nil
}

// ParseDate - create calendar date from string
func ParseDate(trigger string) (Date, error) {
	var err error
	var zero Date
	var d Date
	if err = d.ParseDate(trigger); err != nil {
		return zero, err
	}
	return d, nil
}
