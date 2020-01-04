package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
	tests "github.com/gmax79/otusgolang/microservices/internal/stests"
)

const host = "http://localhost:8888"

var resultCode int
var resultCount int

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

type result struct {
	Count int `json:"result"`
}

func iGetEventsAt(arg1 string) error {
	data, err := tests.GetContent(host, "events_for_day?day="+arg1, http.StatusOK)
	if err != nil {
		return err
	}
	content := string(data)
	content = strings.ReplaceAll(content, "\\", "")
	var r result
	if err = json.Unmarshal([]byte(content), &r); err != nil {
		return err
	}
	resultCount = r.Count
	return nil
}

func countShouldBe(arg1 int) error {
	if resultCount != arg1 {
		return fmt.Errorf("Count not equal. Got %d should be %d", resultCode, arg1)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create event at "([^"]*)" with "([^"]*)"$`, iCreateEventAtWith)
	s.Step(`^response code should be (\d+)$`, responseCodeShouldBe)
	s.Step(`^I get events at "([^"]*)"$`, iGetEventsAt)
	s.Step(`^count should be (\d+)$`, countShouldBe)
}
