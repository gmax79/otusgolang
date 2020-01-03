package main

import (
	"github.com/DATA-DOG/godog"
	tests "github.com/gmax79/otusgolang/microservices/internal/stests"
	"net/http"
)

const host = "http://localhost:8888"

var resultCode int
var response string

func iCreateEventAtWith(arg1, arg2 string) error {
	r := map[string]string{
		"time":  arg1,
		"event": arg2,
	}
	resp, err := tests.PostRequest(host, "create_event", r)
	if err != nil {
		return err
	}
	resultCode = resp.StatusCode
	return nil
}

func responseCodeShouldBe(arg1 int) error {
	if resultCode != arg1 {
		return godog.ErrPending
	}
	return nil
}

func iGetEventAt(arg1 string) error {
	return godog.ErrPending
}

func responseShouldBe(arg1 string) error {
	if response != arg1 {
		r := map[string]string{
			"time": arg1,
		}
		resp, err := tests.Post(host, "get_event", r, http.StatusOK)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return nil
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create event at "([^"]*)" with "([^"]*)"$`, iCreateEventAtWith)
	s.Step(`^response code should be (\d+)$`, responseCodeShouldBe)
	s.Step(`^I get event at "([^"]*)"$`, iGetEventAt)
	s.Step(`^response should be "([^"]*)"$`, responseShouldBe)
}
