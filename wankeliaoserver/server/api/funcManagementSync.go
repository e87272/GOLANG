package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func funcManagementSync(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/funcManagementSync")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/funcManagementSync", r.Header["Client-Name"][0])

	result := map[string]interface{}{}

	pubData := common.Syncdata{Synctype: "funcManagementSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
