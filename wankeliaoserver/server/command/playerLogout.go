package command

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Playerlogout(connect *websocket.Conn, msg []byte, loginUuid string) error {

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
		common.Sendmessage(connect, sendLogoutJson)
		return err
	}
	sendLogout.Base_R.Idem = packetLogout.Base_C.Idem

	sendLogout.Base_R.Result = "ok"
	sendLogoutJson, _ := json.Marshal(sendLogout)
	common.Sendmessage(connect, sendLogoutJson)

	for _, roomInfo := range userRoom {
		common.Roomsclientdelete(roomInfo.Roomuuid, loginUuid)
		if len(common.Roomsread(roomInfo.Roomuuid)) == 0 {
			common.Roomsdelete(roomInfo.Roomuuid)
			common.Roomsinfodelete(roomInfo.Roomuuid)
		}

		// 離開為單一不用通知
		// chatMessage := socket.Chatmessage{Historyuuid: "sys", From: userPlatform, Stamp: timeUnix, Message: "exit room", Style: "exit room"}
		// logoutBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}}
		// logoutBroadcast.Payload.Chatmessage = chatMessage
		// logoutBroadcast.Payload.Chattarget = roomInfo.Roomuuid
		// logoutBroadcastJson, _ := json.Marshal(logoutBroadcast)
		// common.Redispubroomdata(roomInfo.Roomuuid, logoutBroadcastJson)

	}
	common.Clientsdelete(loginUuid)
	common.Usersinfodelete(userPlatform.Useruuid)

	connect.Close()

	return nil
}
