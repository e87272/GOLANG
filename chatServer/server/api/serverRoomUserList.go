package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func serverRoomUserList(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/server/serverRoomUserList")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/serverRoomUserList", r.Header["Client-Name"][0])

	roomUuid := r.FormValue("roomUuid")
	// log.Printf("roomUuid = %s\n", roomUuid)

	result := map[string]interface{}{}
	result["result"] = "ok"

	jsonListByte, err := json.Marshal(common.Roomsreadmemberlist(roomUuid))
	if err != nil {
		result["result"] = "err"
		result["message"] = "JSON_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_SERVERROOMUSERLIST_JSON_ERROR", r.Header["Client-Name"][0], nil)
		return
	}

	result["message"] = string(jsonListByte)

	// log.Printf("Roomspopulation : %+v\n", common.Roomspopulation)

	common.ResponseWithJson(w, http.StatusOK, result)
}
