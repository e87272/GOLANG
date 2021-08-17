package common

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

func SysLog(msg interface{}) {

	type sysLog struct {
		Message interface{} `json:"message"`
		Stamp   string      `json:"stamp"`
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sysLogJson, _ := json.Marshal(sysLog{Message: msg, Stamp: timeUnix})
	log.Printf("%+v ", string(sysLogJson))
	return
}

func SysErrorLog(msg interface{}, err error) {

	type sysErrorLog struct {
		Message interface{} `json:"message"`
		Error   string      `json:"error"`
		Stamp   string      `json:"stamp"`
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	sysErrorLogJson, _ := json.Marshal(sysErrorLog{Message: msg, Error: errMsg, Stamp: timeUnix})
	log.Printf("%+v ", string(sysErrorLogJson))
	return
}
