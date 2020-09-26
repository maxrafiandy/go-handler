package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

// NewGormProp returns new database property
func NewGormProp(host, port, user, pass, db, driver string) *GormProp {
	prop := &GormProp{
		Host:     host,
		Port:     port,
		User:     user,
		Pass:     pass,
		Database: db,
		Driver:   driver,
	}

	switch prop.Driver {
	case "mysql":
		prop.ConnectionString = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			prop.User, prop.Pass, prop.Host, prop.Port, prop.Database)
		// host=myhost port=myport user=gorm dbname=gorm password=mypassword
	case "postgres":
		prop.ConnectionString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			prop.Host, prop.Port, prop.User, prop.Database, prop.Pass)
	// sqlserver://username:password@localhost:1433?database=dbname
	case "mssql":
		prop.ConnectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			prop.User, prop.Pass, prop.Host, prop.Port, prop.Database)
	case "sqlite3":
	default:
		log.Fatalf("%s is not a database dialect.", driver)
	}

	return prop
}

// GormAdd connects and returns specific gorm connection
func GormAdd(alias string, prop *GormProp) *gorm.DB {
	GormProps[alias] = prop
	if gormDBs[alias] == nil {
		var err error
		gormDBs[alias], err = gorm.Open(prop.Driver, prop.ConnectionString)

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
			Addr:     fmt.Sprintf("%s:%s", prop.Host, prop.Port),
			Password: prop.Pass,
			DB:       prop.Database,
		})
	}

	return redisDBs[alias]
}

// RedisGet returns new client of redis host
func RedisGet(alias string) *redis.Client {
	return redisDBs[alias]
}
