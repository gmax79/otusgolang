package main

import (
	"strings"
	"testing"
)

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

func TestErrEnvLoader(t *testing.T) {
	_, err := loadEnvVariables("errenv")
	if err == nil {
		t.Errorf("errenv parsed without error, expected error")
		return
	}
}

type testWriter struct {
	data []byte
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func TestAppRuner(t *testing.T) {
	stdout := &testWriter{}
	stderr := &testWriter{}
	vars, _ := loadEnvVariables("envdir")
	err := runApp("env", []string{}, vars, stdout, stderr)
	if err != nil {
		t.Error(err)
	}
	out := string(stdout.data)
	if !strings.Contains(out, "PROPS=props-value") {
		t.Error("envdir/PROPS!=props-value")
	}
	if !strings.Contains(out, "TEST=test-value") {
		t.Error("envdir/PROPS!=props-value")
	}
}
