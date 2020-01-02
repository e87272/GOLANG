package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func userInfoSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/userInfoSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/userInfoSync", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	userUuid := r.FormValue("userUuid")

	pubData := common.Syncdata{Synctype: "userInfoSync", Data: userUuid}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
