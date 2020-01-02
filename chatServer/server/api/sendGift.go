package api

import (
	"log"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"../common"
	"../database"
	"../socket"
)

func sendGift(w http.ResponseWriter, r *http.Request) {

	ok := common.Apicheck(w, r, "/broadcast/sendGift")
	if !ok {
		return
	}

	rStr, _ := json.Marshal(r.Form)
	common.Essyslog(string(rStr), "/broadcast/sendGift", r.Header["Client-Name"][0])

	result := map[string]interface{}{}

	platformUuid := r.FormValue("platformUuid")
	platform := r.FormValue("platform")
	ownerPlatformUuid := r.FormValue("ownerPlatformUuid")
	imgUrl := r.FormValue("imgUrl")
	// log.Printf("platformUuid = %s\n", platformUuid)
	// log.Printf("platform = %s\n", platform)
	// log.Printf("roomUuid = %s\n", roomUuid)
	// log.Printf("imgUrl = %s\n", imgUrl)

	var userPlatform socket.Userplatform
	row := database.QueryRow("select uuid,platformUuid,platform from users where platformUuid = ? and platform = ?", platformUuid, platform)
	err := row.Scan(&userPlatform.Useruuid, &userPlatform.Platformuuid, &userPlatform.Platform)
	// log.Printf("user : %+v\n", user)
	if err != nil {

		result["result"] = "err"
		result["message"] = "DB_NO_DATA"
		common.ResponseWithJson(w, http.StatusOK, result)
		common.Essyserrorlog("API_SENDGIFT_DB_NO_DATA", r.Header["Client-Name"][0], err)
		return
	}

	roomUuid, ok, exception := common.Hierarchytokensearch(r.Header["Client-Name"][0], r.Header["Client-Name"][0], ownerPlatformUuid, "MM")
	
	if !ok {

		result["result"] = "err"
		result["message"] = exception.Message
		common.ResponseWithJson(w, http.StatusOK, result)
		return
	}

	historyUuid := common.Getid().Hexstring()
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	chatMessage := socket.Chatmessage{Historyuuid: historyUuid, From: userPlatform, Stamp: timeUnix, Message: imgUrl, Style: "gift"}
	sendGiftBroadcast := socket.Cmd_b_send_gift{Base_B: socket.Base_B{Cmd: socket.CMD_B_SEND_GIFT, Stamp: timeUnix}}
	sendGiftBroadcast.Payload.Chatmessage = chatMessage
	sendGiftBroadcast.Payload.Chattarget = roomUuid
	sendGiftBroadcastJson, _ := json.Marshal(sendGiftBroadcast)

	log.Printf("sendGiftBroadcastJson : %+v\n", sendGiftBroadcast)

	common.Redispubroomdata(roomUuid, sendGiftBroadcastJson)

	result["result"] = "ok"
	result["message"] = "ok"
	common.ResponseWithJson(w, http.StatusOK, result)

	chatMessageHsitory := common.Chathistory{Historyuuid: historyUuid, Chattarget: roomUuid, Myuuid: userPlatform.Useruuid, Myplatformuuid: userPlatform.Platformuuid, Myplatform: userPlatform.Platform, Stamp: timeUnix, Message: imgUrl, Style: "send gift"}
	// Index a second tweet (by string)
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	// log.Printf("chatMessageHsitory : %+v\n", chatMessageHsitory)

	err = common.Esinsert(os.Getenv("sysRoomLog"), string(chatMessageJson))

	// log.Printf("Playersendmsg err : %+v\n", err)
	return
}
