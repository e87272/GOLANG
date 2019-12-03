package commandRoom

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../socket"
)

func Playerenterroombatch(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendEnterRoomBatch := socket.Cmd_r_player_enter_room_batch{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_ENTER_ROOM_BATCH,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userInfo, _ := common.Usersinforead(userPlatform.Useruuid)
	userUuid := userPlatform.Useruuid

	var packetEnterRoomBatch socket.Cmd_c_player_enter_room_batch
	if err := json.Unmarshal([]byte(msg), &packetEnterRoomBatch); err != nil {
		sendEnterRoomBatch.Base_R.Result = "err"
		sendEnterRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOMBATCH_JSON_ERROR", userUuid, err)
		sendEnterRoomBatchJson, _ := json.Marshal(sendEnterRoomBatch)
		common.Sendmessage(connect, sendEnterRoomBatchJson)
		return err
	}
	sendEnterRoomBatch.Base_R.Idem = packetEnterRoomBatch.Base_C.Idem

	roomBatchrResult := []struct {
		Result      string             `json:"result"`
		Roominfo    socket.Roominfo    `json:"roomInfo"`
		Lastmessage socket.Chatmessage `json:"lastMessage"`
		Membercount int                `json:"memberCount"`
	}{}

	for key, roomCore := range packetEnterRoomBatch.Payload {
		roomBatchrResult = append(roomBatchrResult, struct {
			Result      string             `json:"result"`
			Roominfo    socket.Roominfo    `json:"roomInfo"`
			Lastmessage socket.Chatmessage `json:"lastMessage"`
			Membercount int                `json:"memberCount"`
		}{})

		roomType := roomCore.Roomtype
		roomUuid := roomCore.Roomuuid

		if roomUuid == "" || len(roomUuid) != 16 {
			code := common.Essyserrorlog("COMMAND_PLAYERENTERROOMBATCH_ROOM_UUID_NULL", userPlatform.Useruuid, nil)
			roomBatchrResult[key].Result = code
			continue
		}

		switch roomType {
		case "privateGroup":
			if strings.Index(userInfo.Privategroup, roomUuid) == -1 {
				//block處理
				code := common.Essyserrorlog("COMMAND_PLAYERENTERROOMBATCH_ROOM_UUID_NULL", userPlatform.Useruuid, nil)
				roomBatchrResult[key].Result = code
			}
		case "vipGroup":
			if strings.Index(userInfo.Vipgroup, roomUuid) == -1 {
				//block處理
				code := common.Essyserrorlog("COMMAND_PLAYERENTERROOMBATCH_ROOM_UUID_NULL", userPlatform.Useruuid, nil)
				roomBatchrResult[key].Result = code
				continue
			}
		default:
			code := common.Essyserrorlog("COMMAND_PLAYERENTERROOMBATCH_ROOM_TYPE_NULL", userPlatform.Useruuid, nil)
			roomBatchrResult[key].Result = code
			continue
		}

		if common.Checkinroom(roomUuid, loginUuid) {
			code := common.Essyserrorlog("COMMAND_PLAYERENTERROOMBATCH_IN_ROOM", userPlatform.Useruuid, nil)
			roomBatchrResult[key].Result = code
			continue
		}

		roomInfo, ok, exception := common.Hierarchyroominfosearch(loginUuid, client, roomType, roomUuid)
		if !ok {
			roomBatchrResult[key].Result = exception.Code
			continue
		}

		memberCount, ok, exception := common.Hierarchymembercount(loginUuid, client, roomType, roomUuid)
		if !ok {
			roomBatchrResult[key].Result = exception.Code
			continue
		}

		common.Clientsroominsert(loginUuid, roomUuid, socket.Roomcore{Roomuuid: roomInfo.Roomuuid, Roomtype: roomInfo.Roomtype})

		lastMessage := socket.Chatmessage{}
		lastMessage = common.Getredisroomlastmessage(roomUuid)

		roomBatchrResult[key].Lastmessage = lastMessage
		roomBatchrResult[key].Roominfo = roomInfo
		roomBatchrResult[key].Membercount = memberCount
	}

	sendEnterRoomBatch.Base_R.Result = "ok"
	sendEnterRoomBatch.Payload = roomBatchrResult
	sendEnterRoomBatchJson, _ := json.Marshal(sendEnterRoomBatch)
	common.Sendmessage(connect, sendEnterRoomBatchJson)

	for _, result := range roomBatchrResult {

		if result.Result != "" {
			continue
		}

		firstFromUuid := common.Getredisfirstenterroom(result.Roominfo.Roomuuid + userPlatform.Useruuid)
		if firstFromUuid == "" {
			continue
		}
		common.Deleteredisfirstenterroom(result.Roominfo.Roomuuid + userPlatform.Useruuid)

		// common.Essyslog("firstenterroom Roomname : "+result.Roominfo.Roomname, loginUuid, client.Userplatform.Useruuid)
		historyUuid := common.Getid().Hexstring()
		chatMessage := socket.Chatmessage{Historyuuid: historyUuid, From: userPlatform, Stamp: timeUnix, Message: "enter room", Style: "sys"}
		roomBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_ENTER_ROOM, Stamp: timeUnix}}
		roomBroadcast.Payload.Chatmessage = chatMessage
		roomBroadcast.Payload.Chattarget = result.Roominfo.Roomuuid
		roomBroadcastJson, _ := json.Marshal(roomBroadcast)
		common.Redispubroomsinfo(result.Roominfo.Roomuuid, roomBroadcastJson)

		common.Setredisroomlastmessage(result.Roominfo.Roomuuid, chatMessage)

		chatMessageHsitory := common.Chathistory{
			Historyuuid:    historyUuid,
			Chattarget:     result.Roominfo.Roomuuid,
			Myuuid:         userUuid,
			Myplatformuuid: userPlatform.Platformuuid,
			Myplatform:     userPlatform.Platform,
			Stamp:          timeUnix,
			Message:        chatMessage.Message,
			Style:          chatMessage.Style,
		}
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert(os.Getenv(result.Roominfo.Roomtype), string(chatMessageJson))
		if err != nil {
			common.Essyslog("Esinsert "+os.Getenv(result.Roominfo.Roomtype)+" err: "+err.Error(), loginUuid, userUuid)
		}

	}
	return nil
}
