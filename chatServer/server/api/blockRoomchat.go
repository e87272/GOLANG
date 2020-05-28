package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"server/common"
	"server/database"
	"server/socket"
)

func blockRoomChat(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/blockRoomchat")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/blockRoomchat", r.Header["Client-Name"][0])

	result := map[string]interface{}{}

	platformUuid := r.FormValue("platformUuid")
	platform := r.FormValue("platform")
	targetUuid := r.FormValue("targetUuid")
	blockType := r.FormValue("blockType")
	blockTime, err := strconv.ParseInt(r.FormValue("blockTime"), 10, 64)
	blockTime = blockTime*60*1000 + time.Now().UnixNano()/int64(time.Millisecond)
	if err != nil {
		result["result"] = "err"
		result["message"] = "TIME_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_BLOCKROOMCHAT_TIME_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	// log.Printf("platformUuid = %s\n", platformUuid)
	// log.Printf("platform = %s\n", platform)
	// log.Printf("targetUuid = %s\n", targetUuid)
	// log.Printf("blockType = %s\n", blockType)
	// log.Printf("blockTime = %+v\n", blockTime)

	var userPlatform socket.Userplatform
	row := database.QueryRow(
		"select uuid,platformUuid,platform from users where platformUuid = ? and platform = ?",
		platformUuid,
		platform,
	)
	err = row.Scan(&userPlatform.Useruuid, &userPlatform.Platformuuid, &userPlatform.Platform)
	// log.Printf("user : %+v\n", user)
	if err != nil {
		result["result"] = "err"
		result["message"] = "DB_NO_DATA"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_BLOCKROOMCHAT_DB_NO_DATA", r.Header["Client-Name"][0], err)
		return
	}

	if blockType == "user" {
		targetUuid = userPlatform.Useruuid
	}

	var blockUuid = common.Getid().Hexstring()
	_, err = database.Exec(
		"DELETE FROM `chatBlock` WHERE blockUserUuid = ? and blocktarget = ? ",
		userPlatform.Useruuid,
		targetUuid,
	)
	if err != nil {
		result["result"] = "err"
		result["message"] = "DB_DELETE_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_BLOCKROOMCHAT_DB_DELETE_ERROR", r.Header["Client-Name"][0], err)
		return
	}
	_, err = database.Exec(
		"INSERT INTO chatBlock (blockUuid, blockUserUuid, blocktarget , blockType, platformUuid, platform, timeStamp) VALUES (?, ? , ? , ? , ? , ? , ? )",
		blockUuid,
		userPlatform.Useruuid,
		targetUuid,
		blockType,
		platformUuid,
		platform,
		blockTime,
	)

	if err != nil {
		result["result"] = "err"
		result["message"] = "DB_INSERT_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_BLOCKROOMCHAT_DB_INSERT_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	result["result"] = "ok"
	result["message"] = blockUuid

	common.ResponseWithJson(w, http.StatusOK, result)

	//更新列表
	pubData := common.Syncdata{Synctype: "blockSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))
	// log.Printf("BlockchatList : %+v\n", common.BlockchatList)

}
