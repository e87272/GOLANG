package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func roomInfoSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/roomInfoSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/roomInfoSync", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	roomUuid := r.FormValue("roomUuid")
	roomType := r.FormValue("roomType")

	result["message"] = common.Queryroominfo(r.Header["Client-Name"][0], roomType, roomUuid)

	if result["message"] != "ok" {
		result["result"] = "err"
		// result["message"] = "API_ROOMINFOSYNC_QUERYROOMINFO_ERROR"
	} else {

		result["result"] = "ok"

		//更新列表
		data := map[string]string{}
		data["roomType"] = roomType
		data["roomUuid"] = roomUuid
		dataJson, _ := json.Marshal(data)

		pubData := common.Syncdata{Synctype: "roomsInfoSync", Data: string(dataJson)}
		pubDataJson, _ := json.Marshal(pubData)
		common.Redispubdata("sync", string(pubDataJson))
	}
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
