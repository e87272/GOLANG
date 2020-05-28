package common

import (
	"encoding/json"
	"context"
	"time"

	"server/socket"
)

var sideTextLastMessagePrefix string = "sidetext_lastMessage_"
var sideTextLastSeenPrefix string = "sidetext_lastSeen_"

var roomLastMessagePrefix string = "room_lastMessage_"
var roomLastSeenPrefix string = "room_lastSeen_"

var firstenterroomPrefix string = "enter_room_first_"

var roomInfoPrefix string = "roominfo_"

var userInfoPrefix string = "userinfo_"

var memberCountPrefix string = "member_count_"

var roomStationPrefix string = "roomstation_"

var userTokenPrefix string = "userToken_"

func Getredissidetextlastmessage(key string) (socket.Chatmessage, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),sideTextLastMessagePrefix + key).Result()
	// log.Printf("Getredissidetextlastmessage key : %+v\n", key)
	// log.Printf("Getredissidetextlastmessage result : %+v\n", result)
	if err != nil {
		Essyserrorlog("COMMON_GETREDISSIDETEXTLASTMESSAGE_ERROR", "sideTextUuid:"+sideTextLastMessagePrefix+key, err)
		return socket.Chatmessage{}, false
	}
	var chatMessage socket.Chatmessage
	json.Unmarshal([]byte(result), &chatMessage)
	// log.Printf("Getredissidetextlastmessage chatMessage : %+v\n", chatMessage)
	return chatMessage, true
}

func Setredissidetextlastmessage(key string, chatMessage socket.Chatmessage) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredissidetextlastmessage chatMessage : %+v\n", chatMessage)
	chatMessageJson, _ := json.Marshal(chatMessage)
	//設0存永久
	// log.Printf("Setredissidetextlastmessage string(chatMessageJson) : %+v\n", string(chatMessageJson))
	// log.Printf("Setredissidetextlastmessage key : %+v\n", sideTextLastMessagePrefix+key)
	
	Redisclient.Set(context.Background(),sideTextLastMessagePrefix+key, string(chatMessageJson), 30*24*time.Hour)
}

func Getredissidetextlastseen(key string) string {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),sideTextLastSeenPrefix + key).Result()
	// log.Printf("Getredissidetextlastseen key : %+v\n", key)
	// log.Printf("Getredissidetextlastseen result : %+v\n", result)
	if err != nil {
		// Essyserrorlog("COMMON_GETREDISSIDETEXTLASTSEEN_ERROR", "sideTextUuid:"+sideTextLastSeenPrefix+key, err)
		return ""
	}
	return result
}

func Setredissidetextlastseen(key string, hsitoryUuid string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredissidetextlastseen key : %+v\n", sideTextLastSeenPrefix+key)
	// log.Printf("Setredissidetextlastseen hsitoryUuid : %+v\n", sideTextLastSeenPrefix+key)
	
	Redisclient.Set(context.Background(),sideTextLastSeenPrefix+key, hsitoryUuid, 30*24*time.Hour)
}

func Getredisroomlastmessage(key string) (socket.Chatmessage, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),roomLastMessagePrefix + key).Result()
	// log.Printf("Getredisroomlastmessage key : %+v\n", key)
	// log.Printf("Getredisroomlastmessage result : %+v\n", result)
	if err != nil {
		Essyserrorlog("COMMON_GETREDISROOMLASTMESSAGE_ERROR", "roomUuid:"+roomLastMessagePrefix+key, err)
		return socket.Chatmessage{}, false
	}
	var chatMessage socket.Chatmessage
	json.Unmarshal([]byte(result), &chatMessage)
	// log.Printf("Getredisroomlastmessage chatMessage : %+v\n", chatMessage)
	return chatMessage, true
}

func Setredisroomlastmessage(key string, chatMessage socket.Chatmessage) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredisroomlastmessage chatMessage : %+v\n", chatMessage)
	chatMessageJson, _ := json.Marshal(chatMessage)
	//設0存永久
	// log.Printf("Setredisroomlastmessage string(chatMessageJson) : %+v\n", string(chatMessageJson))
	// log.Printf("Setredisroomlastmessage key : %+v\n", roomLastMessagePrefix+key)
	
	Redisclient.Set(context.Background(),roomLastMessagePrefix+key, string(chatMessageJson), 30*24*time.Hour)
}

func Getredisroomlastseen(key string) string {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),roomLastSeenPrefix + key).Result()
	// log.Printf("Getredisroomlastseen key : %+v\n", key)
	// log.Printf("Getredisroomlastseen result : %+v\n", result)
	if err != nil {
		Essyserrorlog("COMMON_GETREDISROOMLASTSEEN_ERROR", "roomLastSeenKey:"+roomLastSeenPrefix+key, err)
		return ""
	}
	return result
}

func Setredisroomlastseen(key string, hsitoryUuid string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredisroomlastseen key : %+v\n", roomLastSeenPrefix+key)
	// log.Printf("Setredisroomlastseen hsitoryUuid : %+v\n", roomLastSeenPrefix+key)
	
	Redisclient.Set(context.Background(),roomLastSeenPrefix+key, hsitoryUuid, 30*24*time.Hour)
}

func Getredisfirstenterroom(key string) string {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),firstenterroomPrefix + key).Result()
	// log.Printf("Getredisroomlastseen key : %+v\n", key)
	// log.Printf("Getredisroomlastseen result : %+v\n", result)
	if err != nil {
		// Essyserrorlog("COMMON_GETREDISFIRSTENTERROOM_ERROR", "firstenterroomPrefixKey:"+firstenterroomPrefix+key, err)
		return ""
	}
	return result
}

func Setredisfirstenterroom(key string, fromUuid string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredisroomlastseen key : %+v\n", roomLastSeenPrefix+key)
	// log.Printf("Setredisroomlastseen hsitoryUuid : %+v\n", roomLastSeenPrefix+key)
	
	Redisclient.Set(context.Background(),firstenterroomPrefix+key, fromUuid, 30*24*time.Hour)
}

func Deleteredisfirstenterroom(key string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	Redisclient.Del(context.Background(),firstenterroomPrefix + key)
}

func Getredisroominfo(key string) (socket.Roominfo, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),roomInfoPrefix + key).Result()
	if err != nil {
		return socket.Roominfo{}, false
	}
	var roomInfo socket.Roominfo
	json.Unmarshal([]byte(result), &roomInfo)
	return roomInfo, true
}

func Setredisroominfo(key string, roomInfo socket.Roominfo) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	roomInfoJson, _ := json.Marshal(roomInfo)
	
	Redisclient.Set(context.Background(),roomInfoPrefix+key, roomInfoJson, 7*24*time.Hour)
}

func Getredisuserinfo(key string) (socket.User, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	result, err := Redisclient.Get(context.Background(),userInfoPrefix + key).Result()
	if err != nil {
		return socket.User{}, false
	}
	var userInfo socket.User
	json.Unmarshal([]byte(result), &userInfo)
	return userInfo, true
}

func Setredisuserinfo(key string, userInfo socket.User) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	userInfoJson, _ := json.Marshal(userInfo)
	
	Redisclient.Set(context.Background(),userInfoPrefix+key, userInfoJson, 7*24*time.Hour)
}

func Deleteredisuserinfo(key string) {

	Redisclient.Del(context.Background(),userInfoPrefix+key)
}

func Getredismembercount(key string) (int, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	result, err := Redisclient.Get(context.Background(),memberCountPrefix + key).Result()
	if err != nil {
		return 0, false
	}
	var count int
	json.Unmarshal([]byte(result), &count)
	return count, true
}

func Setredismembercount(key string, count int) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	Redisclient.Set(context.Background(),memberCountPrefix+key, count, 7*24*time.Hour)
}

func Getredisroomstation(key string) (string, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	roomUuid, err := Redisclient.Get(context.Background(),roomStationPrefix + key).Result()
	if err != nil {
		return "", false
	}
	// log.Printf("Getredisroomstation key : %+v\n", key)
	// log.Printf("Getredisroomstation roomUuid : %+v\n", roomUuid)
	return roomUuid, true
}

func Setredisroomstation(key string, roomUuid string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	// log.Printf("Setredisroomstation key : %+v\n", key)
	// log.Printf("Setredisroomstation roomUuid : %+v\n", roomUuid)
	
	Redisclient.Set(context.Background(),roomStationPrefix+key, roomUuid, 7*24*time.Hour)
}

func Getredisusertoken(key string) (string, bool) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	token, err := Redisclient.Get(context.Background(),userTokenPrefix + key).Result()
	if err != nil {
		return "", false
	}
	return token, true
}

func Setredisusertoken(key string, token string) {
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	Redisclient.Set(context.Background(),userTokenPrefix+key, token, 30*24*time.Hour)
}