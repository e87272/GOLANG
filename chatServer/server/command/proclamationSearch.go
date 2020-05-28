package command

import (
	"encoding/json"
	"strconv"
	"time"

	"server/common"
	"server/socket"
)

func Proclamationsearch(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendProclamationSearch := socket.Cmd_r_proclamation{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PROCLAMATION,
		Stamp: timeUnix,
	}}

	var packetProclamation socket.Cmd_c_proclamation

	if err := json.Unmarshal([]byte(msg), &packetProclamation); err != nil {
		sendProclamationSearch.Base_R.Result = "err"
		sendProclamationSearch.Base_R.Exp = common.Exception("COMMAND_PROCLAMATIONSEARCH_JSON_ERROR", "", err)
		sendProclamationSearchJson, _ := json.Marshal(sendProclamationSearch)
		common.Sendmessage(connCore, sendProclamationSearchJson)
		return err
	}
	sendProclamationSearch.Base_R.Idem = packetProclamation.Base_C.Idem

	proclamationlist := common.Proclamationsearch(packetProclamation.Payload.Roomuuid)

	sendProclamationSearch.Result = "ok"
	sendProclamationSearch.Payload = proclamationlist
	sendProclamationSearchJson, _ := json.Marshal(sendProclamationSearch)
	common.Sendmessage(connCore, sendProclamationSearchJson)

	return nil
}
