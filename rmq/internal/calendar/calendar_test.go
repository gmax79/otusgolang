package calendar

import (
	"testing"

	"github.com/gmax79/otusgolang/rmq/internal/simple"
)

func test(trigger string, d *simple.Date) error {
	t, err := simple.ParseDate(trigger)
	if err != nil {
		return err
	}
	*d = t
	return nil
}

func TestTimeParser(t *testing.T) {
	var p simple.Date
	if test("10:20:30", &p) != nil {
		t.Fatal("Error at correct time")
	}
	if p.Hour != 10 || p.Minute != 20 || p.Second != 30 {
		t.Fatal("Invalid parsed time")
	}
	if test("20:40", &p) != nil {
		t.Fatal("Error at correct time 2")
	}
	if p.Hour != 20 || p.Minute != 40 || p.Second != 0 {
		t.Fatal("Invalid parsed time 2")
	}
}
