package socket

import (
	"encoding/json"

	"strconv"
	"time"

	"../commonData"
	"../commonFunc"
)

func playerEnterRoom(connCore commonData.ConnCore, msg []byte, userInfo commonData.UserInfo) {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendPlayerEnterRoom := cmd_r_player_enter_room{base_R: base_R{
		Cmd:    CMD_R_PLAYER_ENTER_ROOM,
		Stamp:  timeUnix,
		Result: "ok",
	}}

	var packetPlayerEnterRoom cmd_c_player_enter_room

	if err := json.Unmarshal([]byte(msg), &packetPlayerEnterRoom); err != nil {
		sendPlayerEnterRoom.Result = "err"
		sendPlayerEnterRoom.Exp = commonFunc.Exception("COMMAND_PLAYERENTERROOM_JSON_ERROR", "userUuid : "+userInfo.UserUuid, err)

		socketResultJson, _ := json.Marshal(sendPlayerEnterRoom)
		commonFunc.SendMessage(connCore, socketResultJson)

		return
	}
	sendPlayerEnterRoom.base_R.Idem = packetPlayerEnterRoom.base_C.Idem

	roomUuid := packetPlayerEnterRoom.Payload

	if roomUuid == "" {
		sendPlayerEnterRoom.Result = "err"
		sendPlayerEnterRoom.Exp = commonFunc.Exception("COMMAND_PLAYERENTERROOM_ROOM_UUID_NULL", "userUuid : "+userInfo.UserUuid, nil)
		socketResultJson, _ := json.Marshal(sendPlayerEnterRoom)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}

	commonFunc.ClientsRoomInsert(connCore.LoginUuid, roomUuid)
	commonFunc.RoomsClientInsert(roomUuid, connCore, userInfo.UserUuid)

	socketResultJson, _ := json.Marshal(sendPlayerEnterRoom)
	commonFunc.SendMessage(connCore, socketResultJson)

	historyUuid := commonFunc.GetUuid()
	chatMessage := commonData.ChatMessage{
		HistoryUuid: historyUuid,
		From:        userInfo,
		Stamp:       timeUnix,
		Message:     "enter room",
		Style:       "sys",
	}
	sendBroadcastMsg := cmd_b_player_room_msg{base_B: base_B{Cmd: CMD_B_PLAYER_ROOM_MSG, Stamp: timeUnix}}
	sendBroadcastMsg.Payload.ChatMessage = chatMessage
	sendBroadcastMsg.Payload.ChatTarget = roomUuid
	sendBroadcastMsgJson, _ := json.Marshal(sendBroadcastMsg)

	var roomMessage struct {
		RoomUuid string
		Message  []byte
	}
	roomMessage.RoomUuid = roomUuid
	roomMessage.Message = sendBroadcastMsgJson

	commonFunc.PubRedisJson("roomMessage", roomMessage)

	return
}
