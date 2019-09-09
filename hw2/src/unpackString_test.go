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
	}

	for i, test := range tests {
		if err := checkUnpackString(test[0], test[1]); err != nil {
			indexStr := fmt.Sprintf("[%d] ", i)
			t.Error(indexStr, err)
		}
	}
}
