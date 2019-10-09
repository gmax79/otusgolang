package main

import (
	"io"
	"testing"
)

type dummyReader struct {
	Count int
}

func (r *dummyReader) Read(p []byte) (n int, err error) {
	tocopy := len(p)
	if tocopy > r.Count {
		tocopy = r.Count
	}
	r.Count -= tocopy
	if r.Count == 0 {
		return tocopy, io.EOF
	}
	return tocopy, nil
}

type dummyWriter struct {
}

func (w *dummyWriter) Write(p []byte) (n int, err error) {
	towrite := len(p)
	return towrite, nil
}

func TestCopyIOSimple(t *testing.T) {
	r := dummyReader{Count: 10000}
	w := dummyWriter{}
	copied, err := copyio(&r, &w, 0)
	if copied != 10000 {
		t.Fatalf("Simple copy of 10000 bytes not working")
		return
	}
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestCopyIOLimited(t *testing.T) {
	r := dummyReader{Count: 10000}
	w := dummyWriter{}
	copied, err := copyio(&r, &w, 5000)
	if copied != 5000 {
		t.Fatalf("Limited copy of 5000 bytes not working")
		return
	}
	if err != nil {
		t.Fatal(err)
		return
	}
}
