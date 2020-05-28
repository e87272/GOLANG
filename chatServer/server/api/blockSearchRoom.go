package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"server/common"
	"server/database"
)

func blockSearchRoom(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/blockSearchRoom")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/blockSearchRoom", r.Header["Client-Name"][0])

	roomUuid := r.FormValue("roomUuid")

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	// log.Printf("timeUnix : %s\n", timeUnix)

	rows, _ := database.Query(
		"select blockUuid,blocktarget,platformUuid,platform,timeStamp from chatBlock where  timeStamp >= ? and blocktarget = ? ",
		timeUnix,
		roomUuid,
	)
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
			result["message"] = "DB_SELECT_ERROR"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_BLOCKSEARCHROOM_DB_SELECT_ERROR", r.Header["Client-Name"][0], err)
			return
		}
		blockList[blockUuid] = string(blockDataJson)
	}
	rows.Close()

	// log.Printf("blockList : %+v\n", blockList)

	_, err := json.Marshal(blockList)

	if err != nil {
		result["result"] = "err"
		result["message"] = "JSON_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_BLOCKSEARCHROOM_JSON_ERROR", r.Header["Client-Name"][0], err)
		return
	}
	// log.Printf("Json : %+v\n", blockListJson)

	result["result"] = "ok"
	result["message"] = blockList

	common.ResponseWithJson(w, http.StatusOK, result)
}
