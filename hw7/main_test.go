package main

import "testing"

func TestEnvLoader(t *testing.T) {
	m, err := loadEnvVariables("envdir")
	if err != nil {
		t.Errorf("envdir parsed with error.")
		return
	}
	v, ok := m["PROPS"]
	if !ok {
		t.Errorf("PROPS in envdir not found")
	}
	const propsvalue = "props-value"
	if v != propsvalue {
		t.Errorf("PROPS with wrong value %s, expected props-value", propsvalue)
	}
	v, ok = m["TEST"]
	if !ok {
		t.Errorf("TEST in envdir not found")
	}
	const testvalue = "test-value"
	if v != testvalue {
		t.Errorf("TEST with wrong value %s, expected props-value", testvalue)
	}
}
