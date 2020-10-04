package handler

import (
	"time"

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

// URLQuery struct for server side rendering
type URLQuery struct {
	ItemsPerPage string `schema:"items_per_page" json:"items_per_page"`
	Page         string `schema:"page" json:"page"`
	Keyword      string `schema:"keyword" json:"keyword,omitempty"`
	StartDate    string `schema:"start_date" json:"start_date,omitempty"`
	EndDate      string `schema:"end_date" json:"end_date,omitempty"`
}

// PaginationResult for server side rendering
type PaginationResult struct {
	List         interface{} `json:"list"`
	Keyword      string      `json:"keyword,omitempty"`
	ItemsPerPage int         `json:"items_per_page"`
	TotalItems   int64       `json:"total_items"`
	Page         int         `json:"page"`
	TotalPage    int64       `json:"total_pages"`
	StartDate    *time.Time  `json:"start_date,omitempty"`
	EndDate      *time.Time  `json:"end_date,omitempty"`
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

// Gormv1Prop struct for database connection
type Gormv1Prop struct {
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
