package api

import (
	"encoding/json"
	"net/http"

	"server/common"
)

func roomProclamationSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/roomProclamationSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/roomProclamationSync", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	roomUuid := r.FormValue("roomUuid")

	result["result"] = "ok"

	pubData := common.Syncdata{Synctype: "proclamationSync", Data: roomUuid}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
