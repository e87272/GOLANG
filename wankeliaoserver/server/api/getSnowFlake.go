package api

import (
	"encoding/json"
	"net/http"

	"../common"
)

func getSnowFlake(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/getSnowFlake")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/getSnowFlake", r.Header["Client-Name"][0])

	result := map[string]interface{}{}

	result["result"] = "ok"
	result["message"] = common.Getid().Hexstring()

	common.ResponseWithJson(w, http.StatusOK, result)

	// log.Printf("BlockchatList : %+v\n", common.BlockchatList)

}
