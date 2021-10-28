package gfiat

import (
	"golang.org/x/text/currency"
)

func ParseFiat(s string) (string, error) {
	_, err := currency.ParseISO(s)
	if err != nil {
		return "", err
	}
	return s, nil
}
