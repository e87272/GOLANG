package socket

import (
	"encoding/json"

	"strconv"
	"time"

	"../commonData"
	"../commonFunc"
)

func tokenChange(connCore commonData.ConnCore, msg []byte, userInfo commonData.UserInfo) {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendTokenChange := cmd_r_token_change{base_R: base_R{
		Cmd:    CMD_R_TOKEN_CHANGE,
		Stamp:  timeUnix,
		Result: "ok",
	}}

	var packetToken cmd_c_token_change

	if err := json.Unmarshal([]byte(msg), &packetToken); err != nil {
		sendTokenChange.Result = "err"
		sendTokenChange.Exp = commonFunc.Exception("COMMAND_TOKENCHANGE_JSON_ERROR", "userUuid : "+userInfo.UserUuid, err)
		socketResultJson, _ := json.Marshal(sendTokenChange)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}
	sendTokenChange.base_R.Idem = packetToken.base_C.Idem

	if packetToken.Payload == "" {
		sendTokenChange.Result = "err"
		sendTokenChange.Exp = commonFunc.Exception("COMMAND_TOKENCHANGE_NULL", "userUuid : "+userInfo.UserUuid, nil)
		socketResultJson, _ := json.Marshal(sendTokenChange)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}

	result, ok := commonFunc.AesTokenDecrypt(packetToken.Payload)
	if !ok {
		sendTokenChange.Result = "err"
		sendTokenChange.Exp = commonFunc.Exception("COMMAND_TOKENCHANGE_TOKEN_ERROR", "userUuid : "+userInfo.UserUuid, nil)
		socketResultJson, _ := json.Marshal(sendTokenChange)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}
	if time.Now().UnixNano()/int64(time.Millisecond)-result.Timestamp > 60*1000 {
		sendTokenChange.Result = "err"
		sendTokenChange.Exp = commonFunc.Exception("COMMAND_TOKENCHANGE_TOKEN_EXPIRED", "userUuid : "+userInfo.UserUuid, nil)
		socketResultJson, _ := json.Marshal(sendTokenChange)
		commonFunc.SendMessage(connCore, socketResultJson)
		return
	}
	if result.Content == "" {
		userInfo.UserUuid = connCore.LoginUuid
	} else {
		userInfo.UserUuid = result.Content
	}
	sendTokenChange.Payload = userInfo

	socketResultJson, _ := json.Marshal(sendTokenChange)
	commonFunc.SendMessage(connCore, socketResultJson)

	commonFunc.UsersInfoInsert(userInfo.UserUuid, userInfo)

	var client = commonData.Client{Room: make(map[string]string), ConnCore: connCore, UserInfo: userInfo}
	commonFunc.ClientsInsert(connCore.LoginUuid, client)

	if connCore.LoginUuid == userInfo.UserUuid {
		var connection = make(map[string]commonData.ConnCore)
		connection[connCore.LoginUuid] = connCore
		commonFunc.ClientsConnectionsInsert(userInfo.UserUuid, connection)
		return
	}

	// log.Printf("Tokenchange Clients : %+v\n", common.Clients)
	commonFunc.SetRedisUserInfo(userInfo.UserUuid, userInfo)

	_, ok = commonFunc.ClientsConnectionsRead(userInfo.UserUuid)
	if !ok {
		var userConnect = make(map[string]commonData.ConnCore)
		userConnect[connCore.LoginUuid] = connCore
		commonFunc.ClientsConnectionsInsert(userInfo.UserUuid, userConnect)
	} else {
		commonFunc.ClientsConnectionsLoginUuidInsert(userInfo.UserUuid, connCore.LoginUuid, connCore)
	}
	return
}
