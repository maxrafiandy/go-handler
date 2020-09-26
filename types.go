package handler

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Response struct hold default server response
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error inherits error interface
type Error struct {
	error
	Description string `json:"description"`
	Errors      error  `json:"error,ommitempty"`
}

// Error implementation of error interface
func (m *Error) Error() string {
	return m.Description
}

// URLQuery struct for Server Side Rendering
type URLQuery struct {
	ItemsPerPage string `schema:"items_per_page"`
	Page         string `schema:"page"`
	Keyword      string `schema:"keyword"`
}

// Validate funtion
func (u URLQuery) Validate() error {
	isNumeric := regexp.MustCompile("^[0-9]+$")
	isAlphanumeric := regexp.MustCompile("^[a-zA-Z0-9\\s.]+$")
	return validation.ValidateStruct(&u,
		validation.Field(&u.ItemsPerPage, validation.Match(isNumeric)),
		validation.Field(&u.Page, validation.Match(isAlphanumeric)),
		validation.Field(&u.Keyword, validation.Match(isAlphanumeric)),
	)
}

// GormProp struct for database connection
type GormProp struct {
	Host             string
	Port             string
	User             string
	Pass             string
	Database         string
	Driver           string
	ConnectionString string
}

// RedisProp struct for cached connection
type RedisProp struct {
	Host     string
	Port     string
	Pass     string
	Database int
}
