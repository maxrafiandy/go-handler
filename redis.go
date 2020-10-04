package handler

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	redisDBs map[string]*redis.Client = make(map[string]*redis.Client)
)

// NewRedisProp returns new database property
func NewRedisProp(host, port, pass string, db int) *RedisProp {
	return &RedisProp{
		host:     host,
		port:     port,
		pass:     pass,
		database: db,
	}
}

// AddRedis returns new client of redis host
func AddRedis(alias string, prop *RedisProp) *redis.Client {
	if redisDBs[alias] == nil {
		redisDBs[alias] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", prop.host, prop.port),
			Password: prop.pass,
			DB:       prop.database,
		})
	}

	return redisDBs[alias]
}

// GetRedis returns new client of redis host
func GetRedis(alias string) *redis.Client {
	return redisDBs[alias]
}
