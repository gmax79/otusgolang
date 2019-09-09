package main

import (
	"strconv"
	"strings"
	"unicode"
)

func unpackString(str string) (string, error) {

	result := []string{} //make([]rune, 0, len(str))
	digits := []rune{}
	for _, sym := range str {
		if unicode.IsDigit(sym) {
			digits = append(digits, sym)
		} else {
			s := string(sym)
			if len(digits) == 0 {
				result = append(result, s)
			} else {
				count, err := strconv.Atoi(string(digits))
				if err != nil {
					return "", err
				}
				result = append(result, strings.Repeat(s, count))
			}
			digits = digits[:0]
		}

	}

	return strings.Join(result, ""), nil
}
