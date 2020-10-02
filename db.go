package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

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
