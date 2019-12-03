package commandRoom

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../socket"
)

func Playerenterroom(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendPlayerEnterRoom := socket.Cmd_r_player_enter_room{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_ENTER_ROOM,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetPlayerEnterRoom socket.Cmd_c_player_enter_room
	if err := json.Unmarshal([]byte(msg), &packetPlayerEnterRoom); err != nil {
		sendPlayerEnterRoom.Base_R.Result = "err"
		sendPlayerEnterRoom.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOM_JSON_ERROR", userUuid, err)
		sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
		common.Sendmessage(connect, sendPlayerEnterRoomJson)
		return err
	}
	sendPlayerEnterRoom.Base_R.Idem = packetPlayerEnterRoom.Base_C.Idem

	roomType := packetPlayerEnterRoom.Payload.Roomtype
	roomUuid := packetPlayerEnterRoom.Payload.Roomuuid

	if roomUuid == "" || len(roomUuid) != 16 {
		sendPlayerEnterRoom.Base_R.Result = "err"
		sendPlayerEnterRoom.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOM_ROOM_UUID_NULL", userUuid, nil)
		sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
		common.Sendmessage(connect, sendPlayerEnterRoomJson)
		return nil
	}

	switch roomType {
	case "liveGroup":
	default:
		//block處理
		sendPlayerEnterRoom.Base_R.Result = "err"
		sendPlayerEnterRoom.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOM_ROOM_UUID_NULL", userUuid, nil)
		sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
		common.Sendmessage(connect, sendPlayerEnterRoomJson)
		return nil
	}

	if common.Checkinroom(roomUuid, loginUuid) {
		//block處理
		sendPlayerEnterRoom.Base_R.Result = "err"
		sendPlayerEnterRoom.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOM_IN_ROOM", userUuid, nil)
		sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
		common.Sendmessage(connect, sendPlayerEnterRoomJson)
		return nil
	}

	roomInfo, ok, exception := common.Hierarchyroominfosearch(loginUuid, client, roomType, roomUuid)
	if !ok {
		//block處理
		sendPlayerEnterRoom.Base_R.Result = "err"
		sendPlayerEnterRoom.Base_R.Exp = exception
		sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
		common.Sendmessage(connect, sendPlayerEnterRoomJson)
		return nil
	}

	common.Clientsroominsert(loginUuid, roomUuid, socket.Roomcore{Roomuuid: roomInfo.Roomuuid, Roomtype: roomInfo.Roomtype})

	sendPlayerEnterRoom.Base_R.Result = "ok"
	sendPlayerEnterRoom.Payload = roomInfo
	sendPlayerEnterRoomJson, _ := json.Marshal(sendPlayerEnterRoom)
	common.Sendmessage(connect, sendPlayerEnterRoomJson)

	historyUuid := common.Getid().Hexstring()
	chatMessage := socket.Chatmessage{Historyuuid: historyUuid, From: userPlatform, Stamp: timeUnix, Message: "enter room", Style: "sys"}
	roomBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_ENTER_ROOM, Stamp: timeUnix}}
	roomBroadcast.Payload.Chatmessage = chatMessage
	roomBroadcast.Payload.Chattarget = roomInfo.Roomuuid
	roomBroadcastJson, _ := json.Marshal(roomBroadcast)
	common.Redispubroomsinfo(roomUuid, roomBroadcastJson)

	return nil
}
