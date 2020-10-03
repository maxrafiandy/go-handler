package handler

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Response struct hold default server response
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error inherits error interface
type Error struct {
	error
	Description string `json:"description"`
	Errors      error  `json:"error,omitempty"`
}

// Validator interface
type Validator interface {
	Validate() error
}

// Error implementation of error interface
func (e *Error) Error() string {
	return e.Description
}

// URLQuery struct for Server Side Rendering
type URLQuery struct {
	ItemsPerPage string `schema:"items_per_page" json:"items_per_page"`
	Page         string `schema:"page" json:"page"`
	Keyword      string `schema:"keyword" json:"keyword,omitempty"`
	StartDate    string `schema:"start_date" json:"start_date,omitempty"`
	EndDate      string `schema:"end_date" json:"end_date,omitempty"`
}

// Validate implements Validatior Validate
func (u URLQuery) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ItemsPerPage, isNumeric),
		validation.Field(&u.Page, isNumeric),
		validation.Field(&u.Keyword, isAlphaNum),
		validation.Field(&u.StartDate, isDate),
		validation.Field(&u.EndDate, isDate),
	)
}

// GormProp struct for database connection
type GormProp struct {
	host             string
	port             string
	user             string
	pass             string
	database         string
	driver           string
	connectionString string
}

// RedisProp struct for cached connection
type RedisProp struct {
	host     string
	port     string
	pass     string
	database int
}
