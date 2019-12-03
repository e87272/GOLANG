package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func dirtywordSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/dirtywordSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/dirtywordSync", r.Header["Client-Name"][0])

	result := map[string]interface{}{}

	pubData := common.Syncdata{Synctype: "dirtywordSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
