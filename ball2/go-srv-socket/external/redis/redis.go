package redis

import (
	"os"
	"sync"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var mutexRedis = new(sync.Mutex)

func RedisInit() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redisHost"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
