package goredis

import (
	"github.com/go-redis/redis/v7"
)

func New(address, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
