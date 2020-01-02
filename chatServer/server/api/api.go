package api

import (
	"net/http"

	"../common"
)

func Api() {

	common.Queryapikey()
	// var s Limit = make(chan int, 1)

	// http.Handle("/broadcast", s)

	http.HandleFunc("/server/serverRoomUserList", serverRoomUserList)
	http.HandleFunc("/server/serverAliveList", serverAliveList)

	http.HandleFunc("/broadcast/subscription", subscription)
	http.HandleFunc("/broadcast/sendGift", sendGift)

	http.HandleFunc("/emit/roomPopulation", roomPopulation)
	http.HandleFunc("/emit/roomPopulationList", roomPopulationList)

	http.HandleFunc("/emit/blockSearchRoom", blockSearchRoom)
	http.HandleFunc("/emit/blockSearchUser", blockSearchUser)
	http.HandleFunc("/emit/blockRoomchat", blockRoomChat)
	http.HandleFunc("/emit/createRoom", createRoom)
	http.HandleFunc("/emit/roomInfoEdit", roomInfoEdit)

	http.HandleFunc("/emit/tokenChange", tokenChange)
	http.HandleFunc("/emit/createApiClient", createApiClient)

	http.HandleFunc("/emit/dirtywordSync", dirtywordSync)
	http.HandleFunc("/emit/userInfoSync", userInfoSync)
	http.HandleFunc("/emit/roomInfoSync", roomInfoSync)
	http.HandleFunc("/emit/funcManagementSync", funcManagementSync)
	http.HandleFunc("/emit/globalMessageSync", globalMessageSync)
	http.HandleFunc("/emit/roomProclamationSync", roomProclamationSync)

	http.HandleFunc("/emit/getSnowFlake", getSnowFlake)

	http.HandleFunc("/emit/setLog", setLog)
}
