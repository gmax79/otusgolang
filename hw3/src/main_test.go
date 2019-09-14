package main

import (
	"testing"
)

func TestGetWord(t *testing.T) {
	tests := [][2]string{
		{"привет,", "привет"},
		{"веревка", "веревка"},
		{"'Шкаф", "Шкаф"},
		{"#абзац#", "абзац"},
		{"-", ""},
	}
	for i, test := range tests {
		result := getWord(test[0])
		if result != test[1] {
			t.Errorf("[%d] GetWord from %s retutns %s, required %s", i, test[0], result, test[1])
		}
	}
}
