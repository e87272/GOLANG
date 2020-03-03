package commonFunc

import (
	"../external/redis"
)

func FuncInit() {

	go redis.Subscribe("test", SubRedisTest)

}
