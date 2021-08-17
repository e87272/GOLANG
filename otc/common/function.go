package common

import (
	"crypto/rand"
	"math/big"
	"net"
	"strings"
)

func MyIp() string {

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	// log.Printf("localAddr : %+v\n", localAddr)

	idx := strings.LastIndex(localAddr, ":")

	myIp := localAddr[0:idx]

	return myIp
}

func RandomInt(n int) int {
	var result, _ = rand.Int(rand.Reader, big.NewInt(int64(n)))
	return int(result.Uint64())
}

func RandomString(n int) string {
	targetChar := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := ""
	for i := 0; i < n; i++ {
		j := RandomInt(len(targetChar))
		result += string(targetChar[j])
	}
	return result
}
