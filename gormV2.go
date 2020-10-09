package handler

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	gormDBs           map[string]*gorm.DB = make(map[string]*gorm.DB)
	defaultGormConfig *gorm.Config        = &gorm.Config{
		PrepareStmt: true,
	}
)

func getGormConfig(config *gormConfig) *gorm.Config {
	if config == nil {
		return defaultGormConfig
	}
	var gormConfig gorm.Config
	copier.Copy(&gormConfig, config)

	return &gormConfig
}

func gormDebug(db *gorm.DB) *gorm.DB {
	var debug string
	debug = strings.ToLower(os.Getenv("DEBUG_MODE"))

	if debug == "true" || debug == "1" {
		return db.Debug()
	}

	return db
}

// NewGormConfig return initial gormConfig
func NewGormConfig() *gormConfig {
	var config gormConfig

	config.SkipDefaultTransaction = false
	config.PrepareStmt = false

	return &config
}

// GetGormDB returns database instance
func GetGormDB(alias string) *gorm.DB {
	return gormDebug(gormDBs[alias])
}

// ConnectMysql open connection to mysql server.
// Example datasource "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func ConnectMysql(alias, dataSource string, config *gormConfig) {
	var err error

	if gormDBs[alias] == nil {
		gormDBs[alias], err = gorm.Open(mysql.Open(dataSource), getGormConfig(config))

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}
}

// ConnectPostgres open connection to postgres server
// "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func ConnectPostgres(alias, dataSource string, config *gormConfig) {
	var err error

	if gormDBs[alias] == nil {
		gormDBs[alias], err = gorm.Open(postgres.Open(dataSource), getGormConfig(config))

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}
}

// ConnectMssql open connection to postgres server
// "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
func ConnectMssql(alias, dataSource string, config *gormConfig) {
	var err error

	if gormDBs[alias] == nil {
		gormDBs[alias], err = gorm.Open(sqlserver.Open(dataSource), getGormConfig(config))

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}
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
// Warning: must supply executed gorm.DB object () ex: PageResult(db.Find(&users), urlQuery)
func PageResult(dbResult *gorm.DB, resultSet interface{}, urlQuery URLQuery) *PaginationResult {
	var (
		count     int64
		limit     int = 10
		page      int = 0
		result    PaginationResult
		startDate time.Time
		endDate   time.Time
	)

	result.List = resultSet
	dbResult.Offset(-1).Count(&count)

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
