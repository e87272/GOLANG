package shell

import (
	"encoding/json"
	"regexp"
	"strings"

	"../common"
	"../socket"
)

func Shell(connCore common.Conncore, loginUuid string, userPlatform socket.Userplatform, packetSendShell socket.Cmd_c_player_send_shell, timeUnix string) error {

	shellCmd := regexp.MustCompile(" +").Split(strings.Trim(packetSendShell.Payload.Message, " "), -1)
	// shellCmd[0] == "/su"
	if len(shellCmd) <= 1 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELL_SHELL_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	switch shellCmd[1] {
	case "LP": //proclamation
		return linkProclamation(connCore, userPlatform, packetSendShell, timeUnix)
	case "NP": //proclamation
		return normalProclamation(connCore, userPlatform, packetSendShell, timeUnix)
	case "BU": //blockuser
		return blockUser(connCore, userPlatform, packetSendShell, timeUnix)
	case "BL": //blockList
		return queryBlockList(connCore, userPlatform, packetSendShell, timeUnix)
	default:
		return shelldemo(shellCmd, connCore, userPlatform, packetSendShell, timeUnix)

		// SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SHELL_ERROR", userPlatform.Useruuid, nil)}}
		// SendMsgJson, _ := json.Marshal(SendMsg)
		// common.Sendmessage(connCore, SendMsgJson)
		// return nil
	}
}
