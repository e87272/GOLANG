package api

import (
	"encoding/json"
	"net/http"

	"server/common"
)

func roomPopulation(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/roomPopulation")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/roomPopulation", r.Header["Client-Name"][0])

	roomUuid := r.FormValue("roomUuid")
	result := map[string]interface{}{}
	// log.Printf("roomUuid = %s\n", roomUuid)
	if roomUuid == "" {
		result["result"] = "err"
		result["message"] = "ROOMUUID_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_ROOMPOPULATION_ROOMUUID_ERROR", r.Header["Client-Name"][0], nil)
		return
	}

	result["result"] = "ok"
	result["message"] = common.Roomspopulationcountread(roomUuid)

	// log.Printf("Roomspopulation : %+v\n", common.Roomspopulation)

	common.ResponseWithJson(w, http.StatusOK, result)
}
