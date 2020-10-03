package handler

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
		log.Fatalf("%s is not a database dialect.", driver)
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
			errs := &Error{
				Description: err.Error(),
				Errors:      err,
			}
			Logger(errs)
			log.Fatal(errs)
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
func Pagination(db *gorm.DB, u URLQuery) *gorm.DB {
	// set limit
	nlimit := 10
	if len(u.ItemsPerPage) != 0 {
		n, _ := strconv.Atoi(u.ItemsPerPage)
		nlimit = n
	}
	db = db.Limit(nlimit)

	// set page offset
	npage := 0
	if len(u.Page) != 0 {
		n, _ := strconv.Atoi(u.Page)
		npage = (n - 1) * nlimit
	}
	db = db.Offset(npage)

	return db
}

// DateRange set single date or range date.
// if there are 2 dates passed from url - start_date and end_date -
// it will return a between start_date and end_date. if it only pass
// start_date then it return where date(column) = start_date,
// otherwise date(column) = NOW()
func DateRange(db *gorm.DB, vars URLQuery, column string) (*gorm.DB, error) {
	var (
		between    string = fmt.Sprintf("date(%s) = ?", column)
		expression string = `^([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))$`
		regexps    *regexp.Regexp
		err        error = &Error{Description: invalidDateFormat}
	)

	checkValid := func(value string) bool {
		regexps, _ = regexp.Compile(expression)
		return regexps.MatchString(value)
	}

	switch {
	// if start and end sets
	case vars.StartDate != "" && vars.EndDate != "":
		if !checkValid(vars.StartDate) || !checkValid(vars.EndDate) {
			return db, err
		}
		between = fmt.Sprintf("date(%s) between ? and ?", column)
		return db.Where(between, vars.StartDate, vars.EndDate), nil
	case vars.StartDate != "":
		if !checkValid(vars.StartDate) {
			return db, err
		}
		return db.Where(between, vars.StartDate), nil
	case vars.EndDate != "":
		if !checkValid(vars.EndDate) {
			return db, err
		}
		between = fmt.Sprintf("date(%s) between ? and ?", column)
		return db.Where(between, time.Now().Format(formatDateYMD), vars.EndDate), nil
	default:
		return db.Where(between, time.Now().Format(formatDateYMD)), nil
	}
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
