package commonFunc

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	_ "net/http/pprof"

	"../data"
	"../external/elasticSearch"
)

func EsSysLog(msg string, loginUuid string, userUuid string) {

	if msg == "" {
		return
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sysErrorLog := data.SysLog{Code: loginUuid, UserUuid: userUuid, Message: msg, Stamp: timeUnix}
	sysErrorLogJson, _ := json.Marshal(sysErrorLog)
	elasticSearch.EsInsert(os.Getenv("sysLog"), string(sysErrorLogJson))

	return
}

func EsSysErrorLog(msg string, userUuid string, err error) string {

	if msg == "" {
		return ""
	}

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	var code = "(" + bkdrHash(msg, 36, 5) + ")"

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sysErrorLog := data.SysErrorLog{UserUuid: userUuid, Code: code, Message: msg, Error: errMsg, Stamp: timeUnix}
	sysErrorLogJson, _ := json.Marshal(sysErrorLog)
	elasticSearch.EsInsert(os.Getenv("sysErrorLog"), string(sysErrorLogJson))

	return code
}

// BKDR-Hash
func bkdrHash(text string, base int64, length int) string {
	const seed = int64(131)

	var divisor int64 = 1
	for i := 0; i < length; i++ {
		divisor *= base
	}

	var hash = int64(0)
	var textByte = []byte(text)
	var textLength = len(textByte)
	for i := 0; i < textLength; i++ {
		hash = (hash*seed + int64(textByte[i])) % divisor
	}

	var code = strconv.FormatInt(hash, int(base))
	var codeLength = len(code)
	if codeLength < length {
		code = strings.Repeat("0", length-codeLength) + code
	}

	return code
}
