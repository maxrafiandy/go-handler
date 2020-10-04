package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"    // driver for mssql
	_ "github.com/jinzhu/gorm/dialects/mysql"    // driver for mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" // driver for postgres
)

var (
	gormV1DBs map[string]*gorm.DB = make(map[string]*gorm.DB)
)

// NewGormProp returns new database property
func NewGormProp(host, port, user, pass, db, driver string) *Gormv1Prop {
	prop := &Gormv1Prop{
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

// AddGormDBv1 connects and returns specific gorm connection
func AddGormDBv1(alias string, prop *Gormv1Prop) *gorm.DB {
	GormProps[alias] = prop
	if gormV1DBs[alias] == nil {
		var err error
		gormV1DBs[alias], err = gorm.Open(prop.driver, prop.connectionString)

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
			gormV1DBs[alias] = gormV1DBs[alias].Debug()
		}
	}

	return gormV1DBs[alias].Set("gorm:auto_preload", true)
}

// GetGormDBv1 returns the "name" gorm connection
func GetGormDBv1(alias string) *gorm.DB {
	return gormV1DBs[alias].Set("gorm:auto_preload", true)
}
