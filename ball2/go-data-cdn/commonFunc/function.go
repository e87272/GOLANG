package commonFunc

import (
	"net"
	"strings"

	"../data"

	_ "net/http/pprof"
)

func InitFunc() {

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

func Exception(msg string, userUuid string, err error) data.Exception {

	if msg == "" {
		return data.Exception{}
	}

	var code = EsSysErrorLog(msg, userUuid, err)

	return data.Exception{Code: code, Message: msg}
}
