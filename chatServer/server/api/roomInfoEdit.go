package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"../common"
	"../database"
)

func roomInfoEdit(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/roomInfoEdit")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/roomInfoEdit", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	platformUuid := r.FormValue("platformUuid")
	platform := r.FormValue("platform")
	roomUuid := r.FormValue("roomUuid")
	roomType := r.FormValue("roomType")
	roomName := r.FormValue("roomName")
	file, _, err := r.FormFile("roomIcon")
	roomIconLink := ""
	if err == nil {
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			result["result"] = "err"
			result["message"] = "ROOM_ICON_ERROR"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_ROOMINFOEDIT_ROOM_ICON_ERROR", r.Header["Client-Name"][0], err)
			return
		}
		if len(bytes) > 3*1024*104 {
			result["result"] = "err"
			result["message"] = "ROOM_ICON_TOO_LARGE"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_ROOMINFOEDIT_ROOM_ICON_TOO_LARGE", r.Header["Client-Name"][0], nil)
			return
		}
		thisTime := time.Now()
		fileName := "/" + common.Getid().Hexstring() + ".png"
		uploadPath := os.Getenv("uploadPath") + strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month()))
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		roomIconLink = "/" + strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month())) + fileName
		ioutil.WriteFile(uploadPath+fileName, bytes, os.ModePerm)
	}

	switch roomType {
	case "privateGroup", "vipGroup", "liveGroup":
	default:
		result["result"] = "err"
		result["message"] = "ROOM_TYPE_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_ROOMINFOEDIT_ROOM_TYPE_ERROR_LARGE", r.Header["Client-Name"][0], nil)
		return
	}

	var ownerUuid string
	row := database.QueryRow(
		"select uuid from users where platformUuid = ? and platform = ?",
		platformUuid,
		platform,
	)
	err = row.Scan(&ownerUuid)
	if err != nil {
		result["result"] = "err"
		result["message"] = "SELECT_USER_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_ROOMINFOEDIT_SELECT_USER_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	if roomIconLink != "" {
		_, err = database.Exec("UPDATE "+roomType+" SET roomName = ? , roomIcon = ? , owner = ? where roomUuid = ?", roomName, roomIconLink, ownerUuid, roomUuid)
	} else {
		_, err = database.Exec("UPDATE "+roomType+" SET roomName = ? , owner = ? where roomUuid = ?", roomName, ownerUuid, roomUuid)
	}
	if err != nil {
		result["result"] = "err"
		result["message"] = "UPDATE_ROOM_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_ROOMINFOEDIT_UPDATE_ROOM_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	result["message"] = common.Queryroominfo(r.Header["Client-Name"][0], roomType, roomUuid)

	if result["message"] != "ok" {
		result["result"] = "err"
		result["message"] = "API_ROOMINFOEDIT_QUERYROOMINFO_ERROR"
	} else {

		result["result"] = "ok"

		//更新列表
		data := map[string]string{}
		data["roomType"] = roomType
		data["roomUuid"] = roomUuid
		dataJson, _ := json.Marshal(data)

		pubData := common.Syncdata{Synctype: "roomsInfoSync", Data: string(dataJson)}
		pubDataJson, _ := json.Marshal(pubData)
		common.Redispubdata("sync", string(pubDataJson))
	}

	result["result"] = "ok"
	result["message"] = roomUuid
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
