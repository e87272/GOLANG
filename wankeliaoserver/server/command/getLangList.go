package command

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Getlanglist(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	SendGetLangList := socket.Cmd_r_get_lang_list{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_LANG_LIST,
		Stamp: timeUnix,
	}}

	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetGetLangList socket.Cmd_c_get_lang_list

	if err := json.Unmarshal([]byte(msg), &packetGetLangList); err != nil {
		SendGetLangList.Base_R.Result = "err"
		SendGetLangList.Base_R.Exp = common.Exception("COMMAND_GETLANGLIST_JSON_ERROR", userUuid, err)
		SendGetLangListJson, _ := json.Marshal(SendGetLangList)
		common.Sendmessage(connect, SendGetLangListJson)
		return err
	}
	SendGetLangList.Base_R.Idem = packetGetLangList.Base_C.Idem

	common.Mutexmutilangerrormsg.Lock()
	defer common.Mutexmutilangerrormsg.Unlock()
	mutiLangErrorMsg, ok := common.Mutilangerrormsg[packetGetLangList.Payload]
	if !ok {
		mutiLangErrorMsg, ok = common.Mutilangerrormsg["zh-CN"]
		if !ok {
			SendGetLangList.Base_R.Result = "err"
			SendGetLangList.Base_R.Exp = common.Exception("COMMAND_GETLANGLIST_SERVER_LANG_FILE_ERROR", userUuid, nil)
			SendGetLangListJson, _ := json.Marshal(SendGetLangList)
			common.Sendmessage(connect, SendGetLangListJson)
			return nil
		}
	}

	SendGetLangList.Result = "ok"
	SendGetLangList.Payload = mutiLangErrorMsg
	SendGetLangListJson, _ := json.Marshal(SendGetLangList)
	common.Sendmessage(connect, SendGetLangListJson)

	return nil
}
