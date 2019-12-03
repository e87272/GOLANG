package command

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../shell"
	"../socket"
)

func Playersendshell(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	// log.Printf("timeUnix : %s\n", timeUnix)

	var packetSendShell socket.Cmd_c_player_send_shell

	if err := json.Unmarshal([]byte(msg), &packetSendShell); err != nil {
		// log.Printf("Playersendmsg err : %+v\n", err)
		return err
	}

	userPlatform, _ := common.Clientsuserplatformread(loginUuid)

	if strings.Split(packetSendShell.Payload.Message, " ")[0] == "/su" {
		return shell.Shell(connect, loginUuid, userPlatform, packetSendShell, timeUnix)
	}

	return nil
}
