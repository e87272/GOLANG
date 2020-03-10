package socket

import (
	"encoding/json"
	"strconv"
	"time"

	"../commonData"
	"../commonFunc"
)

func playerSendMsg(connCore commonData.ConnCore, msg []byte, userInfo commonData.UserInfo) {

	stamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeUnix := strconv.FormatInt(stamp, 10)

	sendPlayerMsg := cmd_r_player_send_msg{base_R: base_R{
		Cmd:    CMD_R_PLAYER_SEND_MSG,
		Stamp:  timeUnix,
		Result: "ok",
	}}
	var packetSendMsg cmd_c_player_send_msg

	if err := json.Unmarshal([]byte(msg), &packetSendMsg); err != nil {
		sendPlayerMsg.Result = "err"
		sendPlayerMsg.Exp = commonFunc.Exception("COMMAND_PLAYERSENDMSG_JSON_ERROR", "userUuid : "+userInfo.UserUuid, err)

		socketResultJson, _ := json.Marshal(sendPlayerMsg)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}
	sendPlayerMsg.base_R.Idem = packetSendMsg.base_C.Idem

	//禁止訪客發話
	// if connCore.LoginUuid == userInfo.UserUuid {
	// 	sendPlayerMsg.Result = "err"
	// 	sendPlayerMsg.Exp = commonFunc.Exception("COMMAND_PLAYERSENDMSG_GUEST", "userUuid : "+userInfo.UserUuid, nil)

	// 	socketResultJson, _ := json.Marshal(sendPlayerMsg)
	// 	commonFunc.SendMessage(connCore, socketResultJson)
	// 	return
	// }

	if !commonFunc.CheckInRoom(packetSendMsg.Payload.ChatTarget, connCore.LoginUuid) {
		sendPlayerMsg.Result = "err"
		sendPlayerMsg.Exp = commonFunc.Exception("COMMAND_PLAYERSENDMSG_NOT_IN_ROOM", "userUuid : "+userInfo.UserUuid, nil)
		socketResultJson, _ := json.Marshal(sendPlayerMsg)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}

	socketResultJson, _ := json.Marshal(sendPlayerMsg)
	commonFunc.SendMessage(connCore, socketResultJson)

	historyUuid := commonFunc.GetUuid()
	chatMessage := commonData.ChatMessage{
		HistoryUuid: historyUuid,
		From:        userInfo,
		Stamp:       timeUnix,
		Message:     packetSendMsg.Payload.Message,
		Style:       packetSendMsg.Payload.Style,
	}

	chatMessageHsitory := commonData.ChatHistory{
		HistoryUuid: historyUuid,
		Chattarget:  packetSendMsg.Payload.ChatTarget,
		MyUuid:      userInfo.UserUuid,
		Stamp:       timeUnix,
		Message:     chatMessage.Message,
		Style:       chatMessage.Style,
	}
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	commonFunc.EsSysLog(string(chatMessageJson), connCore.LoginUuid, userInfo.UserUuid)

	sendBroadcastMsg := cmd_b_player_room_msg{base_B: base_B{Cmd: CMD_B_PLAYER_ROOM_MSG, Stamp: timeUnix}}
	sendBroadcastMsg.Payload.ChatMessage = chatMessage
	sendBroadcastMsg.Payload.ChatTarget = packetSendMsg.Payload.ChatTarget
	sendBroadcastMsgJson, _ := json.Marshal(sendBroadcastMsg)

	var roomMessage struct {
		RoomUuid string
		Message  []byte
	}
	roomMessage.RoomUuid = packetSendMsg.Payload.ChatTarget
	roomMessage.Message = sendBroadcastMsgJson

	commonFunc.PubRedisJson("roomMessage", roomMessage)

	return
}
