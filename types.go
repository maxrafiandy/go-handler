package handler

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type (
	// Response struct hold default server response
	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	// Error inherits error interface
	Error struct {
		error
		Description string `json:"description"`
		Errors      error  `json:"error,omitempty"`
	}

	// Validator interface
	Validator interface {
		Validate() error
	}

	// URLQuery struct for server side rendering
	URLQuery struct {
		ItemsPerPage string `schema:"items_per_page" json:"items_per_page"`
		Page         string `schema:"page" json:"page"`
		Keyword      string `schema:"keyword" json:"keyword,omitempty"`
		StartDate    string `schema:"start_date" json:"start_date,omitempty"`
		EndDate      string `schema:"end_date" json:"end_date,omitempty"`
	}

	// PaginationResult for server side rendering
	PaginationResult struct {
		List         interface{} `json:"list"`
		Keyword      string      `json:"keyword,omitempty"`
		ItemsPerPage int         `json:"items_per_page"`
		TotalItems   int64       `json:"total_items"`
		Page         int         `json:"page"`
		TotalPage    int64       `json:"total_pages"`
		StartDate    *time.Time  `json:"start_date,omitempty"`
		EndDate      *time.Time  `json:"end_date,omitempty"`
	}

	// gormv1Config struct for database connection
	gormv1Config struct {
		host             string
		port             string
		user             string
		pass             string
		database         string
		driver           string
		connectionString string
	}

	// RedisOptions inherits redis.Options
	RedisOptions struct {
		redis.Options
	}

	// GormConfig inherits gorm.Config
	GormConfig struct {
		gorm.Config
		dataSource string
	}

	// Model inherit from gorm model
	Model struct {
		ID        uint           `gorm:"primarykey" json:"id"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	}

	// GormDB type of gormdb
	GormDB *gorm.DB
)

// Error implementation of error interface
func (e *Error) Error() string {
	return e.Description
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
