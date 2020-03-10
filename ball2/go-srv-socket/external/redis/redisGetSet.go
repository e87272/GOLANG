package redis

import (
	"time"
)

func GetRedis(key string) (string, error) {
	mutexRedis.Lock()
	defer mutexRedis.Unlock()
	result, err := redisClient.Get(key).Result()
	return result, err
}

func SetRedis(key string, str string, expiration time.Duration) {
	mutexRedis.Lock()
	defer mutexRedis.Unlock()

	redisClient.Set(key, str, expiration)
}
