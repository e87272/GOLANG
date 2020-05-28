package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"server/common"
	"server/database"
)

func createRoom(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/emit/createRoom")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/emit/createRoom", r.Header["Client-Name"][0])

	result := map[string]interface{}{}
	roomType := r.FormValue("roomType")
	roomName := r.FormValue("roomName")
	platformUuid := r.FormValue("platformUuid")
	platform := r.FormValue("platform")
	file, _, err := r.FormFile("roomIcon")
	roomIconLink := ""
	if err == nil {
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			result["result"] = "err"
			result["message"] = "ROOM_ICON_ERROR"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_CREATEROOM_ROOM_ICON_ERROR", r.Header["Client-Name"][0], err)
			return
		}
		if len(bytes) > 3*1024*104 {
			result["result"] = "err"
			result["message"] = "ROOM_ICON_TOO_LARGE"
			common.ResponseWithJson(w, http.StatusOK, result)
			common.Essyserrorlog("API_CREATEROOM_ROOM_ICON_TOO_LARGE", r.Header["Client-Name"][0], nil)
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
		common.Essyserrorlog("API_CREATEROOM_ROOM_TYPE_ERROR", r.Header["Client-Name"][0], nil)
		return
	}
	roomUuid, ok, exception := common.Hierarchytokensearch("api", r.Header["Client-Name"][0], platformUuid, platform)

	if !ok {
		result["result"] = "err"
		result["message"] = exception.Message
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEROOM_ROOM_SEARCH_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	if roomUuid != "" {
		result["result"] = "err"
		result["message"] = "ROOM_EXIST"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEROOM_ROOM_EXIST", r.Header["Client-Name"][0], err)
		return
	}

	var userUuid string
	row := database.QueryRow(
		"select uuid from users where platformUuid = ? and platform = ?",
		platformUuid,
		platform,
	)
	err = row.Scan(&userUuid)
	if err != nil {
		result["result"] = "err"
		result["message"] = "SELECT_USER_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEROOM_SELECT_USER_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	roomUuid = common.Getid().Hexstring()
	_, err = database.Exec(
		"INSERT INTO `"+roomType+"` (roomUuid, roomName, roomIcon, owner, station, ownerPlatformUuid) VALUES (?, ?, ?, ?, ?, ?)",
		roomUuid,
		roomName,
		roomIconLink,
		userUuid,
		r.Header["Client-Name"][0],
		platformUuid,
	)
	if err != nil {
		result["result"] = "err"
		result["message"] = "INSERT_ROOM_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEROOM_INSERT_ROOM_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	chatMessageHsitory := common.Chathistory{
		Historyuuid:    roomUuid,
		Chattarget:     roomUuid,
		Myuuid:         userUuid,
		Myplatformuuid: platformUuid,
		Myplatform:     platform,
		Stamp:          timeUnix,
		Message:        "create room",
		Style:          "sys",
	}
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	err = common.Esinsert(os.Getenv("privateGroup"), string(chatMessageJson))
	if err != nil {
		common.Essyslog("Esinsert "+os.Getenv("privateGroup")+" err: "+err.Error(), "", userUuid)
	}

	userListName := roomType + "UserList"
	uuid := common.Getid().Hexstring()
	_, err = database.Exec(
		"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?, ?, ?)",
		uuid,
		roomUuid,
		userUuid,
		os.Getenv("roleUuidStreamer"),
	)
	if err != nil {
		result["result"] = "err"
		result["message"] = "INSERT_ROOM_USER_ERROR"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_CREATEROOM_INSERT_ROOM_USER_ERROR", r.Header["Client-Name"][0], err)
		return
	}

	result["result"] = "ok"
	result["message"] = roomUuid
	common.ResponseWithJson(w, http.StatusOK, result)
	return
}
