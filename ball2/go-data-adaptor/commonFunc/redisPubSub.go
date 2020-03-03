package commonFunc

import (
	"encoding/json"
	"log"

	"../external/redis"
)

func PubRedisJson(channel string, pubData interface{}) {

	pubDataJson, _ := json.Marshal(pubData)

	redis.Redispubdata(channel, string(pubDataJson))
}

func SubRedisRoomMessage(jsonData string) {

	var roomMessage struct {
		RoomUuid string
		Message  []byte
	}

	if err := json.Unmarshal([]byte(jsonData), &roomMessage); err != nil {
		log.Printf("SubRedisRoomMessage json err : %+v\n", err)
		return
	}

	// go Broadcast(roomMessage.RoomUuid, roomMessage.Message)
}
