package command

import (
	"encoding/json"
	"strconv"
	"time"

	"../common"
	"../socket"
)

func Playerlogout(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendLogout := socket.Cmd_r_player_logout{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_LOGOUT,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userRoom := client.Room
	userUuid := userPlatform.Useruuid
	

	var packetLogout socket.Cmd_c_player_logout

	if err := json.Unmarshal([]byte(msg), &packetLogout); err != nil {
		sendLogout.Base_R.Result = "err"
		sendLogout.Base_R.Exp = common.Exception("COMMAND_PLAYERLOGOUT_JSON_ERROR", userUuid, err)
		sendLogoutJson, _ := json.Marshal(sendLogout)
		common.Sendmessage(connCore, sendLogoutJson)
		return err
	}
	sendLogout.Base_R.Idem = packetLogout.Base_C.Idem

	sendLogout.Base_R.Result = "ok"
	sendLogoutJson, _ := json.Marshal(sendLogout)
	common.Sendmessage(connCore, sendLogoutJson)

	for _, roomCore := range userRoom {
		common.Roomsclientdelete(roomCore.Roomuuid, loginUuid)
		if len(common.Roomsread(roomCore.Roomuuid)) == 0 {
			common.Roomsdelete(roomCore.Roomuuid)
			common.Roomsinfodelete(roomCore.Roomuuid)
		}

		// 離開為單一不用通知

	}
	common.Clientsdelete(loginUuid)
	common.Usersinfodelete(userPlatform.Useruuid)

	connCore.Conn.Close()

	return nil
}
