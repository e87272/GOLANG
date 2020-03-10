package socket

import (
	"encoding/json"
	"strconv"
	"time"

	"../commonData"
	"../commonFunc"
)

func liveGameInfo(connCore commonData.ConnCore, msg []byte, userInfo commonData.UserInfo) {

	stamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeUnix := strconv.FormatInt(stamp, 10)

	sendLiveGameInfo := cmd_r_live_game_info{base_R: base_R{
		Cmd:    CMD_R_LIVE_GAME_INFO,
		Stamp:  timeUnix,
		Result: "ok",
	}}
	var packetLiveGameInfo cmd_c_live_game_info

	if err := json.Unmarshal([]byte(msg), &packetLiveGameInfo); err != nil {

		sendLiveGameInfo.Result = "err"
		sendLiveGameInfo.Exp = commonFunc.Exception("COMMAND_PLAYERSENDMSG_JSON_ERROR", "userUuid : "+userInfo.UserUuid, err)

		socketResultJson, _ := json.Marshal(sendLiveGameInfo)
		commonFunc.SendMessage(connCore, socketResultJson)

		return
	}
	sendLiveGameInfo.base_R.Idem = packetLiveGameInfo.base_C.Idem

	sendLiveGameInfo.Payload = commonFunc.GetRedisLiveGameInfo(packetLiveGameInfo.Payload)

	socketResultJson, _ := json.Marshal(sendLiveGameInfo)
	commonFunc.SendMessage(connCore, socketResultJson)

	return
}
