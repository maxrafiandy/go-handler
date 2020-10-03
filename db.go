package handler

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"    // driver for mssql
	_ "github.com/jinzhu/gorm/dialects/mysql"    // driver for mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" // driver for postgres
)

// NewGormProp returns new database property
func NewGormProp(host, port, user, pass, db, driver string) *GormProp {
	prop := &GormProp{
		host:     host,
		port:     port,
		user:     user,
		pass:     pass,
		database: db,
		driver:   driver,
	}

	switch prop.driver {
	case "mysql":
		prop.connectionString = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			prop.user, prop.pass, prop.host, prop.port, prop.database)
		// host=myhost port=myport user=gorm dbname=gorm password=mypassword
	case "postgres":
		prop.connectionString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			prop.host, prop.port, prop.user, prop.database, prop.pass)
	// sqlserver://username:password@localhost:1433?database=dbname
	case "mssql":
		prop.connectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			prop.user, prop.pass, prop.host, prop.port, prop.database)
	case "sqlite3":
	default:
		log.Fatalf("[go-handler] %s is not a database dialect", driver)
	}

	return prop
}

// NewRedisProp returns new database property
func NewRedisProp(host, port, user, pass string, db int) *RedisProp {
	return &RedisProp{
		host:     host,
		port:     port,
		pass:     pass,
		database: db,
	}
}

// GormAdd connects and returns specific gorm connection
func GormAdd(alias string, prop *GormProp) *gorm.DB {
	GormProps[alias] = prop
	if gormDBs[alias] == nil {
		var err error
		gormDBs[alias], err = gorm.Open(prop.driver, prop.connectionString)

		if err != nil {
			err = &Error{
				Description: err.Error(),
				Errors:      err,
			}
			log.Fatalf("[go-handler] fatal: %v", err.Error())
		}

		var debug string = os.Getenv("DEBUG_MODE")
		debug = strings.ToLower(debug)

		if debug == "true" || debug == "1" {
			gormDBs[alias] = gormDBs[alias].Debug()
		}
	}

	return gormDBs[alias].Set("gorm:auto_preload", true)
}

// GormGet returns the "name" gorm connection
func GormGet(alias string) *gorm.DB {
	return gormDBs[alias].Set("gorm:auto_preload", true)
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
func PageResult(db *gorm.DB, urlQuery URLQuery) *PaginationResult {
	var (
		count     int
		limit     int = 10
		page      int = 0
		result    PaginationResult
		startDate time.Time
		endDate   time.Time
	)

	result.List = db.Value
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
	result.TotalPage = int(math.Ceil(float64(count) / float64(limit)))

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

// RedisAdd returns new client of redis host
func RedisAdd(alias string, prop *RedisProp) *redis.Client {
	if redisDBs[alias] == nil {
		redisDBs[alias] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", prop.host, prop.port),
			Password: prop.pass,
			DB:       prop.database,
		})
	}

	return redisDBs[alias]
}

// RedisGet returns new client of redis host
func RedisGet(alias string) *redis.Client {
	return redisDBs[alias]
}
