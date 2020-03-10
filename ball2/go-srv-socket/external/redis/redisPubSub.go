package redis

import (
	"log"
)

func Redispubdata(redisIndex string, pubDataJson string) {

	//加鎖 加鎖 加鎖
	mutexRedis.Lock()
	defer mutexRedis.Unlock()

	log.Printf("redisPubSub Redispubdata redisIndex : %+v\n", redisIndex)
	log.Printf("redisPubSub Redispubdata pubDataJson : %+v\n", pubDataJson)

	err := redisClient.Publish(redisIndex, pubDataJson).Err()

	if err != nil {
		log.Printf("Redispubdata err : %+v\n", err)
		return
	}

}

func Subscribe(channelName string, callBack func(string)) {

	// log.Printf("Subscriberoom \n")

	//参数1 频道名 字符串类型
	Redispubsub := redisClient.Subscribe(channelName)

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive()
	if err != nil {
		log.Printf("Subscribe %+v err : %+v\n", channelName, err)
		return
	}
	ch := Redispubsub.Channel()

	for msg := range ch {
		callBack(msg.Payload)
	}

}
