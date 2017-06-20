package wrappers

import (
	"fmt"

	"gopkg.in/redis.v4"
)

type Redis struct {
	address string
}

// ############## Redis

// NewRedis creates a new redis object
func NewRedis(host string, port int) *redis.Client {

	if port <= 0 {
		port = 6379
	}

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", host, port),
	})
}
