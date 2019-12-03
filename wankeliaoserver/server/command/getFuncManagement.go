package command

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Getfuncmanagement(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	Sendfuncmanagement := socket.Cmd_r_get_func_management{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_FUNC_MANAGEMENT,
		Stamp: timeUnix,
	}}

	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetFuncManagement socket.Cmd_c_get_func_management

	if err := json.Unmarshal([]byte(msg), &packetFuncManagement); err != nil {
		Sendfuncmanagement.Base_R.Result = "err"
		Sendfuncmanagement.Base_R.Exp = common.Exception("COMMAND_GETFUNCMANAGEMENT_JSON_ERROR", userUuid, err)
		SendfuncmanagementJson, _ := json.Marshal(Sendfuncmanagement)
		common.Sendmessage(connect, SendfuncmanagementJson)
		return err
	}
	Sendfuncmanagement.Base_R.Idem = packetFuncManagement.Base_C.Idem

	common.Mutexfunctionmanagement.Lock()
	defer common.Mutexfunctionmanagement.Unlock()
	functionManagement := common.Functionmanagement

	Sendfuncmanagement.Result = "ok"
	Sendfuncmanagement.Payload = functionManagement
	SendfuncmanagementJson, _ := json.Marshal(Sendfuncmanagement)
	common.Sendmessage(connect, SendfuncmanagementJson)

	return nil
}
