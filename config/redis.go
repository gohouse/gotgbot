package config

import (
	redis "github.com/go-redis/redis/v7"
)


var client redis.Client

func redisInit(conf *redis.Options) {
	if conf == nil {
		return
	}
	client = *(redis.NewClient(conf))
}

func Redis() *redis.Client {
	return &client
}