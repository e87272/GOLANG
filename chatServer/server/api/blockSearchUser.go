package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"../common"
	"../database"
)

func blockSearchUser(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/blockSearchUser")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/blockSearchUser", r.Header["Client-Name"][0])

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	// log.Printf("timeUnix : %s\n", timeUnix)

	rows, _ := database.Query("select blockUuid,blocktarget,platformUuid,platform,timeStamp from chatBlock where  timeStamp >= ? and blockType = 'user' ", timeUnix)
	// log.Printf("data : %+v\n", rows)

	var blockUuid string
	var blocktarget string
	var platformUuid string
	var platform string
	var timeStamp string

	var blockList = make(map[string]string)

	result := map[string]interface{}{}
	
	for rows.Next() {

		rows.Scan(&blockUuid, &blocktarget, &platformUuid, &platform, &timeStamp)
		// log.Printf("blockUuid : %+v\n", blockUuid)
		// log.Printf("blocktarget : %+v\n", blocktarget)
		// log.Printf("platformUuid : %+v\n", platformUuid)
		// log.Printf("platform : %+v\n", platform)
		// log.Printf("timeStamp : %+v\n", timeStamp)
		var blockData = make(map[string]string)
		blockData["blockUuid"] = blockUuid
		blockData["blocktarget"] = blocktarget
		blockData["platformUuid"] = platformUuid
		blockData["platform"] = platform
		blockData["timeStamp"] = timeStamp

		blockDataJson, err := json.Marshal(blockData)
		// log.Printf("blockDataJson : %+v\n", blockDataJson)

		// log.Printf("blockDataJsonstring : %s\n", string(blockDataJson))

		if err != nil {
			result["result"] = "err"
			result["message"] = "JSON_ERROR"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_BLOCKSEARCHUSER_JSON_ERROR", r.Header["Client-Name"][0], err)
			return
		}
		blockList[blockUuid] = string(blockDataJson)
	}
	rows.Close()

	// log.Printf("blockList : %+v\n", blockList)

	result["result"] = "ok"
	result["message"] = blockList

	common.ResponseWithJson(w, http.StatusOK, result)
}
