package common

import (
	"encoding/json"
	"context"

	"server/socket"
)

func Redispubdata(redisIndex string, pubDataJson string) {

	//加鎖 加鎖 加鎖
	Mutexredis.Lock()
	defer Mutexredis.Unlock()
	
	ctx := context.Background()
	err := Redisclient.Publish(ctx,redisIndex, pubDataJson).Err()

	if err != nil {
		Essyserrorlog("COMMON_REDISPUBDATA_REDIS_ERROR", "redisIndex : "+redisIndex+"   pubDataJson : "+pubDataJson, err)
	}

}

func Redispubroomdata(roomUuid string, broadcastData []byte) {

	pubData := Redispubsubroomdata{RoomUuid: roomUuid, Datajson: string(broadcastData)}
	pubDataJson, _ := json.Marshal(pubData)
	Redispubdata("room", string(pubDataJson))

}

func Redispubsidetextdata(userUuid string, targetUuid string, sideTextData []byte) {

	pubData := Redispubsubsidetextdata{Useruuid: userUuid, Targetuuid: targetUuid, Datajson: string(sideTextData)}
	// log.Printf("Redispubsidetextdata pubData : %+v\n", pubData)
	pubDataJson, _ := json.Marshal(pubData)
	Redispubdata("sideText", string(pubDataJson))

}

func Redispubroomsinfo(roomUuid string, roomBroadcastJson []byte) {

	pubData := Redispubsubroomsinfo{Ip: Myiplastdigit(), Roomuuid: roomUuid, Usercount: Roomsmembercount(roomUuid), Datajson: string(roomBroadcastJson)}
	pubDataJson, _ := json.Marshal(pubData)
	Redispubdata("roomsinfo", string(pubDataJson))

}

func Redispubfriendupdatestate(userFriend socket.Friendplatform, targetFriend socket.Friendplatform, stateData []byte) {

	pubData := Redispubsubinvitedata{
		Userfriend:   userFriend,
		Targetfriend: targetFriend,
		Datajson:     string(stateData),
	}
	pubDataJson, _ := json.Marshal(pubData)

	syncData := Syncdata{
		Synctype: "friendUpdateState",
		Data:     string(pubDataJson),
	}
	syncDataJson, _ := json.Marshal(syncData)

	Redispubdata("sync", string(syncDataJson))

}
