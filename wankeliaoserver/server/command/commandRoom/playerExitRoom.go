package commandRoom

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../socket"
)

func Playerexitroom(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendExitRoomBatch := socket.Cmd_r_player_exit_room{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_EXIT_ROOM,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	// log.Printf("Playerexitroom : %s\n", loginUuid)
	var packetExitRoom socket.Cmd_c_player_exit_room
	if err := json.Unmarshal([]byte(msg), &packetExitRoom); err != nil {
		sendExitRoomBatch.Base_R.Result = "err"
		sendExitRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYEREXITROOM_JSON_ERROR", userUuid, err)
		sendExitRoomBatchJson, _ := json.Marshal(sendExitRoomBatch)
		common.Sendmessage(connect, sendExitRoomBatchJson)
		return err
	}
	sendExitRoomBatch.Base_R.Idem = packetExitRoom.Base_C.Idem

	roomType := packetExitRoom.Payload.Roomtype
	roomUuid := packetExitRoom.Payload.Roomuuid

	if roomUuid == "" || len(roomUuid) != 16 {
		sendExitRoomBatch.Base_R.Result = "err"
		sendExitRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYEREXITROOM_ROOM_UUID_NULL", userUuid, nil)
		sendExitRoomBatchJson, _ := json.Marshal(sendExitRoomBatch)
		common.Sendmessage(connect, sendExitRoomBatchJson)
		return nil
	}

	// userInfoJson, _ := json.Marshal(userInfo)
	// common.Essyslog("Playerexitroom userInfo : "+string(userInfoJson), loginUuid, client.Userplatform.Useruuid)

	switch roomType {
	case "privateGroup":
	case "vipGroup":
	case "liveGroup":
	default:
		sendExitRoomBatch.Base_R.Result = "err"
		sendExitRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYEREXITROOM_ROOM_TYPE_ERROR", userUuid, nil)
		sendExitRoomBatchJson, _ := json.Marshal(sendExitRoomBatch)
		common.Sendmessage(connect, sendExitRoomBatchJson)
		return nil
	}

	roomCore, ok := common.Clientsroomread(loginUuid, packetExitRoom.Payload.Roomuuid)

	if !ok {
		sendExitRoomBatch.Base_R.Result = "err"
		sendExitRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYEREXITROOM_ROOM_UUID_ERROR", userUuid, nil)
		sendExitRoomBatchJson, _ := json.Marshal(sendExitRoomBatch)
		common.Sendmessage(connect, sendExitRoomBatchJson)
		return nil
	}

	common.Clientsroomdelete(loginUuid, roomCore.Roomuuid)
	common.Roomsclientdelete(roomCore.Roomuuid, loginUuid)

	sendExitRoomBatch.Base_R.Result = "ok"
	sendExitRoomBatchJson, _ := json.Marshal(sendExitRoomBatch)
	common.Sendmessage(connect, sendExitRoomBatchJson)

	// 離開為單一不用通知
	// historyUuid := common.Getid().Hexstring()
	// chatMessage := socket.Chatmessage{Historyuuid: historyUuid, From: userPlatform, Stamp: timeUnix, Message: "exit room", Style: "sys"}
	// roomBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}}
	// roomBroadcast.Payload.Chatmessage = chatMessage
	// roomBroadcast.Payload.Chattarget = roomCore.Roomuuid
	// roomBroadcastJson, _ := json.Marshal(roomBroadcast)
	// common.Redispubroomsinfo(roomCore.Roomuuid, roomBroadcastJson)

	if len(common.Roomsread(roomCore.Roomuuid)) == 0 {
		common.Roomsdelete(roomCore.Roomuuid)
		common.Roomsinfodelete(roomCore.Roomuuid)
	}

	return nil
}
