package command

import (
	"encoding/json"
	"strconv"
	"time"

	"server/common"
	"server/database"
	"server/socket"
)

func Frienddelete(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendDelete := socket.Cmd_r_friend_delete{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_FRIEND_DELETE,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetDelete socket.Cmd_c_friend_delete
	err := json.Unmarshal([]byte(msg), &packetDelete)
	if err != nil {
		sendDelete.Base_R.Result = "err"
		sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_JSON_ERROR", userUuid, err)
		sendDeleteJson, _ := json.Marshal(sendDelete)
		common.Sendmessage(connCore, sendDeleteJson)
		return err
	}
	sendDelete.Base_R.Idem = packetDelete.Base_C.Idem

	if loginUuid == userPlatform.Useruuid && false {
		sendDelete.Base_R.Result = "err"
		sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_GUEST", userUuid, nil)
		sendDeleteJson, _ := json.Marshal(sendDelete)
		common.Sendmessage(connCore, sendDeleteJson)
		return nil
	}

	targetUserUuid := packetDelete.Payload
	targetFriend, ok := common.Userfriendlistuserread(userUuid, targetUserUuid)
	if ok {
		switch targetFriend.State {
		case "friend":
			// 已經是好友
			_, err = database.Exec(
				"DELETE FROM `friendList` WHERE fromUuid = ? and toUuid = ? ",
				userUuid,
				targetUserUuid,
			)
			if err != nil {
				sendDelete.Base_R.Result = "err"
				sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_FRIEND_DELETELIST_ERROR", userUuid, err)
				sendDeleteJson, _ := json.Marshal(sendDelete)
				common.Sendmessage(connCore, sendDeleteJson)
				return err
			}
			_, err = database.Exec(
				"DELETE FROM `friendList` WHERE fromUuid = ? and toUuid = ? ",
				targetUserUuid,
				userUuid,
			)
			if err != nil {
				sendDelete.Base_R.Result = "err"
				sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_FRIEND_DELETELIST_ERROR", userUuid, err)
				sendDeleteJson, _ := json.Marshal(sendDelete)
				common.Sendmessage(connCore, sendDeleteJson)
				return err
			}
		default:
			// 未知狀態
			sendDelete.Base_R.Result = "err"
			sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_FRIEND_STATE_ERROR", userUuid, nil)
			sendDeleteJson, _ := json.Marshal(sendDelete)
			common.Sendmessage(connCore, sendDeleteJson)
			return nil
		}
	} else {
		// 未知狀態
		sendDelete.Base_R.Result = "err"
		sendDelete.Base_R.Exp = common.Exception("COMMAND_FRIENDDELETE_FRIEND_STATE_ERROR", userUuid, nil)
		sendDeleteJson, _ := json.Marshal(sendDelete)
		common.Sendmessage(connCore, sendDeleteJson)

	}

	sendDelete.Base_R.Result = "ok"
	sendDeleteJson, _ := json.Marshal(sendDelete)
	common.Sendmessage(connCore, sendDeleteJson)

	userFriend := socket.Friendplatform{}
	userFriend.Userplatform = userPlatform
	userFriend.State = "delete"
	targetFriend.State = "delete"

	sendDeleteBroadcast := socket.Cmd_b_friend{}
	sendDeleteBroadcast.Base_B.Stamp = timeUnix
	sendDeleteBroadcast.Base_B.Cmd = socket.CMD_B_FRIEND
	sendDeleteBroadcast.Payload = userFriend
	sendDeleteBroadcastJson, _ := json.Marshal(sendDeleteBroadcast)
	common.Redispubfriendupdatestate(userFriend, targetFriend, sendDeleteBroadcastJson)

	return nil
}
