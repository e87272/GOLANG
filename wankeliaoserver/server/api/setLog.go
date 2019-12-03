package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func setLog(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/setLog")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/setLog", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)

	return
}
