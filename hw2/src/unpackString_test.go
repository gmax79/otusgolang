package main

import (
	"errors"
	"fmt"
	"testing"
)

func checkUnpackString(src, good string) error {
	unpacked, err := unpackString(src)
	if err != nil {
		return err
	}
	if good != unpacked {
		errlog := fmt.Sprintf("'%s' unpacked into '%s', required: '%s'", src, unpacked, good)
		return errors.New(errlog)
	}
	return nil
}

func TestUnpackString(t *testing.T) {
	tests := [][2]string{
		{"a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd"},
		{"qwe\\4\\5", "qwe45"},
		{"a\\42\\53", "a44555"},
		{"a\\4b2\\511", "a4bb55555555555"},
	}
	for i, test := range tests {
		if err := checkUnpackString(test[0], test[1]); err != nil {
			indexStr := fmt.Sprintf("[%d] ", i)
			t.Error(indexStr, err)
		}
	}
}

func checkFailUnpackString(src string) error {
	_, err := unpackString(src)
	if err == nil {
		errlog := fmt.Sprintf("'%s' unpacked without errors, but it is incorrect string", src)
		return errors.New(errlog)
	}
	return nil
}

func TestErrorsUnpackString(t *testing.T) {
	tests := []string{
		"45",
		"123b",
		"1a2b",
	}
	for i, test := range tests {
		if err := checkFailUnpackString(test); err != nil {
			indexStr := fmt.Sprintf("[%d] ", i)
			t.Error(indexStr, err)
		}
	}
}
