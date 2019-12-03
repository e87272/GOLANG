package shell

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"../common"
	"../socket"
	"github.com/gorilla/websocket"
)

func queryBlockList(connect *websocket.Conn, userPlatform socket.Userplatform, packetSendMsg socket.Cmd_c_player_send_shell, timeUnix string) error {

	argument := regexp.MustCompile(" +-").Split(strings.Trim(packetSendMsg.Payload.Message, " "), -1)
	shellCmd := regexp.MustCompile(" +").Split(argument[0], -1)

	if len(shellCmd) != 2 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_QUERYBLOCKLIST_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	var blockChatList = make(map[string]string)
	for userUuid, roomMap := range common.BlockchatList {
		for roomUuid, time := range roomMap {
			if roomUuid == packetSendMsg.Payload.Chattarget {
				blockChatList[userUuid] = strconv.FormatInt(time, 10)
			}
		}
	}

	// log.Printf("common.Blocklist : %+v\n", common.BlockchatList)
	// log.Printf("blockList : %+v\n", blockChatList)

	blockChatListJson, err := json.Marshal(blockChatList)

	if err != nil {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_QUERYBLOCKLIST_JSON_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		// log.Printf("err : %+v\n", err)
		return nil
	}
	// log.Printf("Json : %+v\n", blockListJson)

	SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}, Payload: string(blockChatListJson)}
	SendMsgJson, _ := json.Marshal(SendMsg)
	common.Sendmessage(connect, SendMsgJson)

	return nil
}
