package services

import "github.com/go-redis/redis"

// var redisClient *redis.ClusterClient
// var redisClient *redis.SentinelClient
var redisClient *redis.Client

func RedisClient() *redis.Client {
	return redisClient
}

func SetRedisClient(c *redis.Client) {
	redisClient = c
}
