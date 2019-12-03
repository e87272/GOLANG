package command

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Healthcheck(connect *websocket.Conn, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	// log.Printf("timeUnix : %s\n", timeUnix)

	var packetHealthCheck socket.Cmd_c_healthcheck

	if err := json.Unmarshal([]byte(msg), &packetHealthCheck); err != nil {
		return err
	}

	SendHealthCheck := socket.Cmd_r_healthcheck{Base_R: socket.Base_R{Cmd: socket.CMD_R_PING, Idem: packetHealthCheck.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("","",nil)}, Payload: "PONG"}
	SendHealthCheckJson, _ := json.Marshal(SendHealthCheck)
	common.Sendmessage(connect, SendHealthCheckJson)

	return nil
}
