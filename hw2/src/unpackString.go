package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func unpackString(str string) (string, error) {

	result := []string{}
	digits := []rune{}
	slashMask := false
	for i, sym := range str {
		if slashMask {
			slashMask = false
			result = append(result, string(sym))
			continue
		}
		lastdigit := false
		if unicode.IsDigit(sym) {
			digits = append(digits, sym)
			if i == len(str)-1 {
				lastdigit = true
			} else {
				continue
			}
		}
		if len(digits) > 0 {
			count, err := strconv.Atoi(string(digits))
			if err != nil {
				return "", errors.New("very big number, overflow") // atoi with only numbers (IsDigit)
			}
			if len(result) == 0 {
				return "", errors.New("invalid input string")
			}
			if count > 1 {
				lastSymbolIndex := len(result) - 1
				result = append(result, strings.Repeat(result[lastSymbolIndex], count-1))
			}
			digits = digits[:0]
		}
		if sym == '\\' {
			slashMask = true
		} else if !lastdigit {
			result = append(result, string(sym))
		}
	}
	return strings.Join(result, ""), nil
}
