package handler

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	gormDBs map[string]*gorm.DB = make(map[string]*gorm.DB)
)

// GetGormDB returns database instance
func GetGormDB(alias string) *gorm.DB {
	return gormDBs[alias]
}

// ConnectMysql open connection to mysql server
// Example datasource "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func ConnectMysql(alias, dataSource string, config *gorm.Config) *gorm.DB {
	var (
		err error
	)

	if gormDBs[alias] == nil {

		if config == nil {
			config = &gorm.Config{}
		}

		gormDBs[alias], err = gorm.Open(mysql.Open(dataSource), config)

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}

	return gormDBs[alias]
}

// ConnectPostgres open connection to postgres server
// "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func ConnectPostgres(alias, dataSource string, config *gorm.Config) *gorm.DB {
	var (
		err error
	)

	if gormDBs[alias] == nil {

		if config == nil {
			config = &gorm.Config{}
		}

		gormDBs[alias], err = gorm.Open(postgres.Open(dataSource), config)

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}

	return gormDBs[alias]
}

// ConnectMssql open connection to postgres server
// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
func ConnectMssql(alias, dataSource string, config *gorm.Config) *gorm.DB {
	var (
		err error
	)

	if gormDBs[alias] == nil {

		if config == nil {
			config = &gorm.Config{}
		}

		gormDBs[alias], err = gorm.Open(sqlserver.Open(dataSource), config)

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}

	return gormDBs[alias]
}

// Pagination set offset and limit of query.
// Default value of limit is 10, and offset of page 1.
// Note: dateColumn is an optional parameter and the function
// only use dateColumn[0]. Call this function before PageResult
func Pagination(db *gorm.DB, urlQuery URLQuery, dateColumn ...string) *gorm.DB {
	var (
		limit   int    = 10
		page    int    = 0
		between string = "%s = ?"
	)

	if len(urlQuery.ItemsPerPage) != 0 {
		limit, _ = strconv.Atoi(urlQuery.ItemsPerPage)
	}
	db = db.Limit(limit)

	if len(urlQuery.Page) != 0 {
		page, _ = strconv.Atoi(urlQuery.Page)
		page = (page - 1) * limit
	}
	db = db.Offset(page)

	if dateColumn != nil {
		switch {
		// if start and end sets
		case len(urlQuery.StartDate) > 0 && len(urlQuery.EndDate) > 0:
			between = fmt.Sprintf("%s between ? and ?", dateColumn[0])
			return db.Where(between, urlQuery.StartDate, urlQuery.EndDate)
		case len(urlQuery.StartDate) > 0:
			between = fmt.Sprintf(between, dateColumn[0])
			return db.Where(between, urlQuery.StartDate)
		case len(urlQuery.EndDate) > 0:
			between = fmt.Sprintf("%s between ? and ?", dateColumn[0])
			return db.Where(between, time.Now().Format(formatDate), urlQuery.EndDate)
		default:
			between = fmt.Sprintf(between, dateColumn[0])
			return db.Where(between, time.Now().Format(formatDate))
		}
	}

	return db
}

// PageResult handle detail of BuildQuery and returns PaginationResult.
// This should called after Pagination() function to generate offset limit.
// Warning: must supply executed gorm.DB object (already has Value) ex: PageResult(db.Find(&users), urlQuery)
func PageResult(db *gorm.DB, list interface{}, urlQuery URLQuery) *PaginationResult {
	var (
		count     int64
		limit     int = 10
		page      int = 0
		result    PaginationResult
		startDate time.Time
		endDate   time.Time
	)

	result.List = list
	db.Offset(-1).Count(&count)

	if len(urlQuery.ItemsPerPage) != 0 {
		limit, _ = strconv.Atoi(urlQuery.ItemsPerPage)
	}

	if len(urlQuery.Page) != 0 {
		page, _ = strconv.Atoi(urlQuery.Page)
		page = (page - 1) * limit
	}

	result.Keyword = urlQuery.Keyword
	result.ItemsPerPage = limit
	result.Page = page/limit + 1
	result.TotalItems = count
	result.TotalPage = int64(math.Ceil(float64(count) / float64(limit)))

	if urlQuery.StartDate != "" {
		startDate, _ = time.Parse(formatDate, urlQuery.StartDate)
		result.StartDate = &startDate
	}

	if urlQuery.EndDate != "" {
		endDate, _ = time.Parse(formatDate, urlQuery.EndDate)
		result.EndDate = &endDate
	}

	return &result
}
