package commonFunc

import (
	"log"

	"../external/redis"
)

func PubRedisTest(broadcastJson []byte) {

	redis.Redispubdata("test", string(broadcastJson))
}

func SubRedisTest(jsonData string) {

	log.Printf("SubRedisTest : %v\n", jsonData)

}
