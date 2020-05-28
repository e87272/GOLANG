package api

import (
	"encoding/json"
	"net/http"

	"server/common"
)

func globalMessageSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/globalMessageSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/globalMessageSync", r.Header["Client-Name"][0])

	pubData := common.Syncdata{Synctype: "globalMessageSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	result := map[string]interface{}{}
	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
