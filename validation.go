package handler

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// these variable should be constant
var (
	isNumeric  = validation.Match(regexp.MustCompile("^[0-9]+$"))
	isAlphaNum = validation.Match(regexp.MustCompile("^[a-zA-Z0-9\\s.]+$"))
	isDate     = validation.Date("2006-01-02")
)
