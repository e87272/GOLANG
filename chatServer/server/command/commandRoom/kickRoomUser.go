package commandRoom

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"../../common"
	"../../database"
	"../../socket"
)

func Kickroomuser(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendKickRoomUser := socket.Cmd_r_kick_room_user{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_KICK_ROOM_USER,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetKickRoomUser socket.Cmd_c_kick_room_user

	if err := json.Unmarshal([]byte(msg), &packetKickRoomUser); err != nil {
		sendKickRoomUser.Base_R.Result = "err"
		sendKickRoomUser.Base_R.Exp = common.Exception("COMMAND_KICKROOMUSER_JSON_ERROR", userUuid, err)
		sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
		common.Sendmessage(connCore, sendKickRoomUserJson)
		return err
	}
	sendKickRoomUser.Base_R.Idem = packetKickRoomUser.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendKickRoomUser.Base_R.Result = "err"
		sendKickRoomUser.Base_R.Exp = common.Exception("COMMAND_KICKROOMUSER_GUEST", userUuid, nil)
		sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
		common.Sendmessage(connCore, sendKickRoomUserJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetKickRoomUser.Payload.Roomcore.Roomuuid)
	if !ok {
		//block處理
		sendKickRoomUser.Base_R.Result = "err"
		sendKickRoomUser.Base_R.Exp = common.Exception("COMMAND_KICKROOMUSER_ROOM_UUID_ERROR", userUuid, nil)
		sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
		common.Sendmessage(connCore, sendKickRoomUserJson)
		return nil
	}
	if packetKickRoomUser.Payload.Targetuuid != userPlatform.Useruuid && !common.Checkadmin(packetKickRoomUser.Payload.Roomcore.Roomuuid, userPlatform.Useruuid, "KickPartner") {
		//block處理
		sendKickRoomUser.Base_R.Result = "err"
		sendKickRoomUser.Base_R.Exp = common.Exception("COMMAND_KICKROOMUSER_NOT_ADMIN", userUuid, nil)
		sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
		common.Sendmessage(connCore, sendKickRoomUserJson)
		return nil
	}
	if roomInfo.Roomcore.Roomtype == "liveGroup" {
		//block處理
		sendKickRoomUser.Base_R.Result = "err"
		sendKickRoomUser.Base_R.Exp = common.Exception("COMMAND_KICKROOMUSER_ROOM_TYPE_ERROR", userUuid, nil)
		sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
		common.Sendmessage(connCore, sendKickRoomUserJson)
		return nil
	}

	targetUserAry := strings.Split(packetKickRoomUser.Payload.Targetuuid, ",")

	for _, targetUuid := range targetUserAry {

		userListName := roomInfo.Roomcore.Roomtype + "UserList"

		_, _ = database.Exec(
			"DELETE FROM `"+userListName+"` WHERE roomUuid = ? and userUuid = ? ",
			roomInfo.Roomcore.Roomuuid,
			targetUuid,
		)

		//刪除redis資料，下次更新時撈DB
		common.Deleteredisuserinfo(targetUuid)

		common.Membercount(userListName, roomInfo.Roomcore.Roomuuid, userUuid)

		targetKickRoomUser := socket.Cmd_b_kick_room_user{Base_B: socket.Base_B{Cmd: socket.CMD_B_KICK_ROOM_USER, Stamp: timeUnix}}
		targetKickRoomUser.Payload = socket.Roomcore{Roomuuid: roomInfo.Roomcore.Roomuuid, Roomtype: roomInfo.Roomcore.Roomtype}
		targetKickRoomUserJson, _ := json.Marshal(targetKickRoomUser)

		userMessage := common.Redispubsubuserdata{Useruuid: targetUuid, Datajson: string(targetKickRoomUserJson)}
		userMessageJson, _ := json.Marshal(userMessage)
		//更新列表
		pubData := common.Syncdata{
			Synctype: "userInfoSyncAndEmit",
			Data:     string(userMessageJson),
		}
		pubDataJson, _ := json.Marshal(pubData)
		// common.Essyslog(string(pubDataJson), loginUuid, userPlatform.Useruuid)
		common.Redispubdata("sync", string(pubDataJson))

	}

	sendKickRoomUser.Base_R.Result = "ok"
	sendKickRoomUserJson, _ := json.Marshal(sendKickRoomUser)
	common.Sendmessage(connCore, sendKickRoomUserJson)

	return nil
}
