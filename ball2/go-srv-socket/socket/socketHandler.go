package socket

import (
	"encoding/json"
	"net/http"

	"../commonData"
	"../commonFunc"
)

func SocketHandler(w http.ResponseWriter, r *http.Request, connCore commonData.ConnCore) {

	// log.Printf("connect RemoteAddr: %+v\n", connect.RemoteAddr().String())
	// log.Printf("connect LocalAddr: %+v\n", connect.LocalAddr().String())
	userInfo, ok := commonFunc.ClientsUserInfoRead(connCore.LoginUuid)
	if !ok {
		userInfo.UserUuid = connCore.LoginUuid
	}

	//ReadMessage只能讀一次 猜測是因為讀取指標的問題

	_, msg, err := connCore.Conn.ReadMessage()
	if err != nil {
		commonFunc.EsSysErrorLog(string(msg), userInfo.UserUuid, err)
		return
	}

	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	// log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
		commonFunc.EsSysErrorLog(string(msg), userInfo.UserUuid, err)
		return
	}

	commonFunc.EsSysLog(string(msg), connCore.LoginUuid, userInfo.UserUuid)

	switch mapResult["cmd"] {
	case CMD_C_TOKEN_CHANGE:

		tokenChange(connCore, msg, userInfo)

		break
	case CMD_C_PLAYER_ENTER_ROOM:

		playerEnterRoom(connCore, msg, userInfo)

		break
	case CMD_C_PLAYER_SEND_MSG:

		playerSendMsg(connCore, msg, userInfo)

		break
	case CMD_C_LIVE_GAME_INFO:

		liveGameInfo(connCore, msg, userInfo)

		break

	default:

	}

	return
}
