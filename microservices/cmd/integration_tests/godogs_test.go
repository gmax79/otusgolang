package main

import (
	"github.com/DATA-DOG/godog"
	tests "github.com/gmax79/otusgolang/microservices/internal/stests"
)

const host = "http://localhost:8888"

var resultCode int

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

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create event at "([^"]*)" with "([^"]*)"$`, iCreateEventAtWith)
	s.Step(`^response code should be (\d+)$`, responseCodeShouldBe)
}
