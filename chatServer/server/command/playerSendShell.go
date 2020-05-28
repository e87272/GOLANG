package command

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"server/common"
	"server/shell"
	"server/socket"
)

func Playersendshell(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	// log.Printf("timeUnix : %s\n", timeUnix)

	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	

	var packetSendShell socket.Cmd_c_player_send_shell

	if err := json.Unmarshal([]byte(msg), &packetSendShell); err != nil {
		// log.Printf("Playersendmsg err : %+v\n", err)
		return err
	}

	if strings.Split(packetSendShell.Payload.Message, " ")[0] == "/su" {
		return shell.Shell(connCore, loginUuid, userPlatform, packetSendShell, timeUnix)
	}

	return nil
}
