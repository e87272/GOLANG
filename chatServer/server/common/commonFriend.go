package common

import (
	"encoding/json"
	"strconv"
	"time"

	"../socket"
)

func Updatefriendstate(dataJson string) {

	var data Redispubsubinvitedata
	json.Unmarshal([]byte(dataJson), &data)

	sendInviteBroadcastJson := []byte(data.Datajson)
	userFriend := data.Userfriend
	userUuid := userFriend.Userplatform.Useruuid
	targetFriend := data.Targetfriend
	targetUuid := targetFriend.Userplatform.Useruuid

	connectMap, ok := Clientsconnectread(targetUuid)
	if ok {
		if userFriend.State == "delete" {
			Userfriendlistuserdelete(targetUuid, userUuid)
		} else {
			Userfriendlistuserinsert(targetUuid, userUuid, userFriend)
		}
		for _, connect := range connectMap {
			Sendmessage(connect, sendInviteBroadcastJson)
		}
	}

	connectMap, ok = Clientsconnectread(userUuid)
	if ok {
		var sendInviteBroadcast socket.Cmd_b_friend
		json.Unmarshal(sendInviteBroadcastJson, &sendInviteBroadcast)
		sendInviteBroadcast.Payload = targetFriend
		sendInviteBroadcastJson, _ = json.Marshal(sendInviteBroadcast)

		if targetFriend.State == "delete" {
			Userfriendlistuserdelete(userUuid, targetUuid)
		} else {
			Userfriendlistuserinsert(userUuid, targetUuid, targetFriend)
		}
		for _, connect := range connectMap {
			Sendmessage(connect, sendInviteBroadcastJson)
		}
	}
}

func Frienddelete(dataJson string) {

	var data Redispubsubfrienddeletedata
	json.Unmarshal([]byte(dataJson), &data)

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendFriendDeleteBroadcast := socket.Cmd_b_friend_delete{Base_B: socket.Base_B{
		Cmd:   socket.CMD_B_FRIEND_DELETE,
		Stamp: timeUnix,
	}}

	userUuid := data.Useruuid
	targetUuid := data.Targetuuid

	connectMap, ok := Clientsconnectread(targetUuid)
	if ok {
		Userfriendlistuserdelete(targetUuid, userUuid)
		sendFriendDeleteBroadcast.Payload = userUuid
		sendFriendDeleteBroadcastJson, _ := json.Marshal(sendFriendDeleteBroadcast)
		for _, connect := range connectMap {
			Sendmessage(connect, sendFriendDeleteBroadcastJson)
		}
	}

	connectMap, ok = Clientsconnectread(userUuid)
	if ok {
		Userfriendlistuserdelete(userUuid, targetUuid)
		sendFriendDeleteBroadcast.Payload = targetUuid
		sendFriendDeleteBroadcastJson, _ := json.Marshal(sendFriendDeleteBroadcast)
		for _, connect := range connectMap {
			Sendmessage(connect, sendFriendDeleteBroadcastJson)
		}
	}
}
