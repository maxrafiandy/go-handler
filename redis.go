package handler

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/jinzhu/copier"
)

var (
	redisDBs map[string]*redis.Client = make(map[string]*redis.Client)
)

// NewRedisOptions returns new database property
func NewRedisOptions(host, port, pass string, db int) *redisOptions {
	var prop redisOptions

	prop.Addr = fmt.Sprintf("%s:%s", host, port)
	prop.Password = pass
	prop.DB = db

	return &prop
}

// AddRedis returns new client of redis host
func AddRedis(alias string, prop *redisOptions) *redis.Client {
	var (
		err error
		opt redis.Options
	)

	if redisDBs[alias] == nil {
		copier.Copy(&opt, prop)
		redisDBs[alias] = redis.NewClient(&opt)

		if _, err = redisDBs[alias].Ping().Result(); err != nil {
			log.Fatal(err)
		}
	}

	return redisDBs[alias]
}

// GetRedis returns new client of redis host
func GetRedis(alias string) *redis.Client {
	return redisDBs[alias]
}
