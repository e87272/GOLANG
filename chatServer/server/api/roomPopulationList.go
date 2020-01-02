package api

import (
	"net/http"

	"../common"
)

func roomPopulationList(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/roomPopulationList")
	if !ok {
		return
	}

	roomPopulationMap := map[string]int{}
	common.Mutexroomspopulation.Lock()
	defer common.Mutexroomspopulation.Unlock()
	for _, rooms := range common.Roomspopulation {
		for roomUuid, population := range rooms {
			count, _ := roomPopulationMap[roomUuid]
			roomPopulationMap[roomUuid] = count + population
		}
	}
	result := map[string]interface{}{}
	result["result"] = "ok"
	result["message"] = roomPopulationMap

	// log.Printf("Roomspopulation : %+v\n", common.Roomspopulation)

	common.ResponseWithJson(w, http.StatusOK, result)

	// roomPopulationJson, _ := json.Marshal(common.Roomspopulation)
	// common.Essyslog(string(roomPopulationJson), "Roomspopulation", r.Header["Client-Name"][0])

	// roomPopulationMapJson, _ := json.Marshal(roomPopulationMap)
	// common.Essyslog(string(roomPopulationMapJson), "/emit/roomPopulationList", r.Header["Client-Name"][0])

}
