package api

import (
	"net/http"
	"strings"

	"../common"
)

func serverAliveList(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/serverAliveList")
	if !ok {
		return
	}

	aliveIpList := r.FormValue("aliveIpList")
	aliveIpListArray := strings.Split(aliveIpList, ",")

	serverAliveMap := map[string]bool{}

	for _, ip := range aliveIpListArray {
		serverAliveMap[ip] = true
	}

	common.Mutexroomspopulation.Lock()
	defer common.Mutexroomspopulation.Unlock()
	for ip, _ := range common.Roomspopulation {
		_, ok := serverAliveMap[ip]
		if !ok {
			delete(common.Roomspopulation, ip)
		}
	}
	result := map[string]interface{}{}
	result["result"] = "ok"
	result["message"] = "ok"

	// log.Printf("Roomspopulation : %+v\n", common.Roomspopulation)

	common.ResponseWithJson(w, http.StatusOK, result)

}
