package api

import (
	"encoding/json"
	"net/http"

	"../common"
	"../database"
)

func tokenChange(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/tokenChange")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/tokenChange", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	platformUuid := r.FormValue("platformUuid")
	platform := r.FormValue("platform")

	userUuid := ""
	row := database.QueryRow("select uuid from users where platformUuid = ? and platform = ?", platformUuid, platform)
	err := row.Scan(&userUuid)
	if err == database.ErrNoRows {
		userUuid = common.Getid().Hexstring()
		_, err := database.Exec(
			"INSERT INTO users (uuid, platformUuid, platform, globalRole) VALUES (?, ?, ?, ?)",
			userUuid,
			platformUuid,
			platform,
			"",
		)
		if err != nil {

			result["result"] = "err"
			result["message"] = "API_TOKENCHANGE_INSERT_USER_ERROR"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_TOKENCHANGE_INSERT_USER_ERROR", r.Header["Client-Name"][0], err)
			return
		}
	} else if err != nil {
		result["result"] = "err"
		result["message"] = "API_TOKENCHANGE_SELECT_USER_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_TOKENCHANGE_SELECT_USER_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	result["result"] = "ok"
	result["message"] = ""
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
