package commonFunc

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	_ "net/http/pprof"

	"../commonData"
	"../external/elasticSearch"
	"../external/stamp"
)

func EsSysLog(apiName string, msg string) {
	log.Printf("Elasticclient EsSysLog apiName : %+v\n", apiName)

	timeUnix := strconv.FormatInt(stamp.Now(), 10)
	sysLog := commonData.SysLog{ApiName: apiName, Message: msg, Stamp: timeUnix}
	sysLogJson, _ := json.Marshal(sysLog)
	elasticSearch.EsInsert(os.Getenv("sysLog"), string(sysLogJson))

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

	timeUnix := strconv.FormatInt(stamp.Now(), 10)
	sysErrorLog := commonData.SysErrorLog{UserUuid: userUuid, Code: code, Message: msg, Error: errMsg, Stamp: timeUnix}
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
