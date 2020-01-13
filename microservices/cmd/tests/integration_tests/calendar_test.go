package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
	tests "github.com/gmax79/otusgolang/microservices/internal/testshelpers"
)

const host = "http://localhost:8888"

var resultCode int
var resultCount int

type result struct {
	Count int `json:"result"`
}

func parseResultCount(data []byte) error {
	resultCount = -1
	content := string(data)
	content = strings.ReplaceAll(content, "\\", "")
	var r result
	if err := json.Unmarshal([]byte(content), &r); err != nil {
		return err
	}
	resultCount = r.Count
	return nil
}

func iCreateEventAtWith(arg1, arg2 string) error {
	r := map[string]string{
		"time":  arg1,
		"event": arg2,
	}
	resp, err := tests.PostRequest(host, "create_event", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resultCode = resp.StatusCode
	return nil
}

func responseCodeShouldBe(arg1 int) error {
	if resultCode != arg1 {
		return fmt.Errorf("response not equal. Got code %d should be code %d", resultCode, arg1)
	}
	return nil
}

func iGetEventsAt(arg1 string) error {
	data, err := tests.GetContent(host, "events_for_day?day="+arg1, http.StatusOK)
	if err != nil {
		return err
	}
	return parseResultCount(data)
}

func iGetWeekEventsAt(arg1 string) error {
	data, err := tests.GetContent(host, "events_for_week?week="+arg1, http.StatusOK)
	if err != nil {
		fmt.Println("#$$", err)
		return err
	}
	return parseResultCount(data)
}

func iGetMonthEventsAt(arg1 string) error {
	data, err := tests.GetContent(host, "events_for_month?month="+arg1, http.StatusOK)
	if err != nil {
		return err
	}
	return parseResultCount(data)
}

func countShouldBe(arg1 int) error {
	if resultCount != arg1 {
		return fmt.Errorf("Count not equal. Got %d should be %d", resultCount, arg1)
	}
	return nil
}

func iDeleteEventAt(arg1, arg2 string) error {
	r := map[string]string{
		"time":  arg2,
		"event": arg1,
	}
	return tests.Post(host, "delete_event", r, http.StatusOK)
}

func iMoveEventAtTo(arg1, arg2, arg3 string) error {
	rmove := map[string]string{
		"time":    arg2,
		"event":   arg1,
		"newtime": arg3,
	}
	return tests.Post(host, "move_event", rmove, http.StatusOK)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create event at "([^"]*)" with "([^"]*)"$`, iCreateEventAtWith)
	s.Step(`^response code should be (\d+)$`, responseCodeShouldBe)
	s.Step(`^I get events at "([^"]*)"$`, iGetEventsAt)
	s.Step(`^count should be (\d+)$`, countShouldBe)
	s.Step(`^I delete event "([^"]*)" at "([^"]*)"$`, iDeleteEventAt)
	s.Step(`^I move event "([^"]*)" at "([^"]*)" to "([^"]*)"$`, iMoveEventAtTo)
	s.Step(`^I get week events at "([^"]*)"$`, iGetWeekEventsAt)
	s.Step(`^I get month events at "([^"]*)"$`, iGetMonthEventsAt)
}
