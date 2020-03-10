package commonFunc

import (
	"encoding/json"
	"log"
	"time"

	"../commonData"
	"../external/redis"
)

var redisTestPrefix string = "redisTest_"

var userInfoPrefix string = "userInfo_"

var liveGameInfoPrefix string = "liveGameInfo_"

func GetRedisUserInfo(key string) (commonData.UserInfo, bool) {

	result, err := redis.GetRedis(userInfoPrefix + key)
	if err != nil {
		log.Printf("GetRedisUserInfo get err : %+v\n", err)
		return commonData.UserInfo{}, false
	}

	var userInfo commonData.UserInfo

	if err := json.Unmarshal([]byte(result), &userInfo); err != nil {
		log.Printf("GetRedisUserInfo json err : %+v\n", err)
		return commonData.UserInfo{}, false
	}
	return userInfo, true
}

func SetRedisUserInfo(key string, userInfo commonData.UserInfo) {

	userInfosJson, _ := json.Marshal(userInfo)

	redis.SetRedis(userInfoPrefix+key, string(userInfosJson), 30*24*time.Hour)

}

func GetRedisLiveGameInfo(key string) []commonData.ChatMessage {

	result, err := redis.GetRedis(liveGameInfoPrefix + key)
	if err != nil {
		log.Printf("GetRedisLiveGameInfo get err : %+v\n", err)
		return []commonData.ChatMessage{}
	}

	var liveGameInfo []commonData.ChatMessage

	if err := json.Unmarshal([]byte(result), &liveGameInfo); err != nil {
		log.Printf("GetRedisUserInfo json err : %+v\n", err)
		return []commonData.ChatMessage{}
	}
	return liveGameInfo
}
