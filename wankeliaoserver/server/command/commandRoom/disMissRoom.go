package commandRoom

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../database"
	"../../socket"
)

func Dismissroom(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendDisMissRoom := socket.Cmd_r_dis_miss_room{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_DIS_MISS_ROOM,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetDisMissRoom socket.Cmd_c_dis_miss_room

	if err := json.Unmarshal([]byte(msg), &packetDisMissRoom); err != nil {
		sendDisMissRoom.Base_R.Result = "err"
		sendDisMissRoom.Base_R.Exp = common.Exception("COMMAND_DISMISSROOM_JSON_ERROR", userUuid, err)
		sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
		common.Sendmessage(connect, sendDisMissRoomJson)
		return err
	}
	sendDisMissRoom.Base_R.Idem = packetDisMissRoom.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendDisMissRoom.Base_R.Result = "err"
		sendDisMissRoom.Base_R.Exp = common.Exception("COMMAND_DISMISSROOM_GUEST", userUuid, nil)
		sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
		common.Sendmessage(connect, sendDisMissRoomJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetDisMissRoom.Payload)
	if !ok {
		//block處理
		sendDisMissRoom.Base_R.Result = "err"
		sendDisMissRoom.Base_R.Exp = common.Exception("COMMAND_DISMISSROOM_ROOM_UUID_ERROR", userUuid, nil)
		sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
		common.Sendmessage(connect, sendDisMissRoomJson)
		return nil
	}
	if roomInfo.Roomtype != "privateGroup" {
		//block處理
		sendDisMissRoom.Base_R.Result = "err"
		sendDisMissRoom.Base_R.Exp = common.Exception("COMMAND_DISMISSROOM_ROOM_UUID_ERROR", userUuid, nil)
		sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
		common.Sendmessage(connect, sendDisMissRoomJson)
		return nil
	}
	if !common.Checkadmin(packetDisMissRoom.Payload, userPlatform.Useruuid, "DisMissGroup") {
		//block處理
		sendDisMissRoom.Base_R.Result = "err"
		sendDisMissRoom.Base_R.Exp = common.Exception("COMMAND_DISMISSROOM_NOT_ADMIN", userUuid, nil)
		sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
		common.Essyslog(string(sendDisMissRoomJson), loginUuid, userUuid)
		common.Sendmessage(connect, sendDisMissRoomJson)
		return nil
	}

	userListName := roomInfo.Roomtype + "UserList"

	rows, err := database.Query("SELECT userUuid FROM "+userListName+" WHERE roomUuid = ?",
		roomInfo.Roomuuid,
	)
	if err != nil {
		return err
	}
	//群留下不刪除
	// _, _ = database.Exec(
	// 	"DELETE FROM `"+roomInfo.Roomtype+"` WHERE roomUuid = ? ",
	// 	roomInfo.Roomuuid,
	// )

	_, _ = database.Exec(
		"DELETE FROM `"+userListName+"` WHERE roomUuid = ? ",
		roomInfo.Roomuuid,
	)

	//解散逐一踢人通知
	var targetUuid string
	for rows.Next() {
		rows.Scan(&targetUuid)

		targetKickRoomUser := socket.Cmd_b_kick_room_user{Base_B: socket.Base_B{Cmd: socket.CMD_B_KICK_ROOM_USER, Stamp: timeUnix}}
		targetKickRoomUser.Payload = socket.Roomcore{Roomuuid: roomInfo.Roomuuid, Roomtype: roomInfo.Roomtype}
		targetKickRoomUserJson, _ := json.Marshal(targetKickRoomUser)

		userMessage := common.Redispubsubuserdata{Useruuid: targetUuid, Datajson: string(targetKickRoomUserJson)}
		userMessageJson, _ := json.Marshal(userMessage)
		common.Redispubdata("user", string(userMessageJson))
	}
	rows.Close()

	sendDisMissRoom.Base_R.Result = "ok"
	sendDisMissRoomJson, _ := json.Marshal(sendDisMissRoom)
	common.Sendmessage(connect, sendDisMissRoomJson)

	return nil
}
