package commonFunc

import (
	"net"
	"strings"

	"../commonData"
	"../external/redis"

	_ "net/http/pprof"
)

func InitFunc() {

	go redis.Subscribe("roomMessage", SubRedisRoomMessage)

}

func MyIp() string {

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	// log.Printf("localAddr : %+v\n", localAddr)

	idx := strings.LastIndex(localAddr, ":")

	myIp := localAddr[0:idx]

	return myIp
}

func Exception(msg string, userUuid string, err error) commonData.Exception {

	if msg == "" {
		return commonData.Exception{}
	}

	var code = EsSysErrorLog(msg, userUuid, err)

	return commonData.Exception{Code: code, Message: msg}
}
