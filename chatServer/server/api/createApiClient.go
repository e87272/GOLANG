package api

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"

	"server/common"
	"server/database"
)

func createApiClient(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/createApiClient")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/createApiClient", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	clientName := r.FormValue("clientName")

	apiKeyByte := [96]byte{}
	for i := 0; i < len(apiKeyByte); i++ {
		apiKeyByte[i] = byte(rand.Intn(256))
	}
	apiKey := base64.URLEncoding.EncodeToString(apiKeyByte[:])

	_, err := database.Exec(
		"INSERT INTO apiKeyList (apiKey, clientName) VALUES (?, ?)",
		apiKey,
		clientName,
	)
	if err != nil {
		result["result"] = "err"
		result["message"] = "INSERT_KEY_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEAPICLIENT_DB_INSERT_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	result["result"] = "ok"
	result["message"] = apiKey
	common.ResponseWithJson(w, http.StatusOK, result)

	pubData := common.Syncdata{Synctype: "apiKeySync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))
	return
}
