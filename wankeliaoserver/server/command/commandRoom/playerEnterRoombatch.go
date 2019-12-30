package commandRoom

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic"

	"../../common"
	"../../socket"
)

func Playerenterroombatch(connCore common.Conncore, msg []byte, loginUuid string) error {

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
		common.Sendmessage(connCore, sendEnterRoomBatchJson)
		return err
	}
	sendEnterRoomBatch.Base_R.Idem = packetEnterRoomBatch.Base_C.Idem

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendEnterRoomBatch.Base_R.Result = "err"
		sendEnterRoomBatch.Base_R.Exp = common.Exception("COMMAND_PLAYERENTERROOMBATCH_IP_READ_ERROR", userUuid, nil)
		sendEnterRoomBatchJson, _ := json.Marshal(sendEnterRoomBatch)
		common.Sendmessage(connCore, sendEnterRoomBatchJson)
		return nil
	}

	roomBatchrResult := []struct {
		Result      string               `json:"result"`
		Roominfo    socket.Roominfo      `json:"roomInfo"`
		Newmessage  []socket.Chatmessage `json:"newMessage"`
		Lastmessage socket.Chatmessage   `json:"lastMessage"`
		Membercount int                  `json:"memberCount"`
	}{}

	for key, roomCore := range packetEnterRoomBatch.Payload {
		roomBatchrResult = append(roomBatchrResult, struct {
			Result      string               `json:"result"`
			Roominfo    socket.Roominfo      `json:"roomInfo"`
			Newmessage  []socket.Chatmessage `json:"newMessage"`
			Lastmessage socket.Chatmessage   `json:"lastMessage"`
			Membercount int                  `json:"memberCount"`
		}{})

		// log.Printf("roomBatchrResult : %+v\n", roomBatchrResult)

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

		common.Clientsroominsert(loginUuid, roomUuid, socket.Roomcore{Roomuuid: roomUuid, Roomtype: roomType})

		lastMessage := common.Hierarchyroomlastmessage(loginUuid, userUuid, roomCore)

		if roomInfo.Roomicon != "" {
			roomInfo.Roomicon = os.Getenv("linkPath") + roomInfo.Roomicon
		}

		oldLastMessageUuid := common.Getredisroomlastseen(roomType + "_" + roomUuid + "_" + userUuid)
		if oldLastMessageUuid == "" {
			oldLastMessageUuid = common.Getid().Hexstring()
		}

		boolQ := elastic.NewBoolQuery()
		boolQ.Filter(elastic.NewMatchQuery("chatTarget", roomUuid))
		boolQ.Filter(elastic.NewRangeQuery("historyUuid").Gt(oldLastMessageUuid))
		searchResult, err := common.Elasticclient.Search(strings.ToLower(os.Getenv(roomType))).Query(boolQ).Sort("historyUuid", false).Do(context.Background())

		log.Printf("chatTarget : %+v\n", roomUuid)
		log.Printf("oldLastMessageUuid : %+v\n", oldLastMessageUuid)
		log.Printf("strings.ToLower(os.Getenv(roomType)) : %+v\n", strings.ToLower(os.Getenv(roomType)))
		log.Printf("searchResult.Hits.Hits : %+v\n", searchResult.Hits.Hits)

		if err != nil {
			roomBatchrResult[key].Result = common.Exception("COMMAND_PLAYERENTERROOMBATCH_SEARCH_ERROR", userUuid, err).Message
		}

		roomBatchrResult[key].Newmessage = []socket.Chatmessage{}

		for _, hit := range searchResult.Hits.Hits {
			var chatHistory common.Chathistory
			_ = json.Unmarshal(hit.Source, &chatHistory)

			chatMessage := socket.Chatmessage{
				Historyuuid: chatHistory.Historyuuid,
				From: socket.Userplatform{
					Useruuid:     chatHistory.Myuuid,
					Platformuuid: chatHistory.Myplatformuuid,
					Platform:     chatHistory.Myplatform,
				},
				Stamp:   chatHistory.Stamp,
				Message: chatHistory.Message,
				Style:   chatHistory.Style,
				Ip:      chatHistory.Ip,
			}
			roomBatchrResult[key].Newmessage = append(roomBatchrResult[key].Newmessage, chatMessage)
		}

		roomBatchrResult[key].Lastmessage = lastMessage
		roomBatchrResult[key].Roominfo = roomInfo
		roomBatchrResult[key].Membercount = memberCount

		log.Printf("roomBatchrResult[key] : %+v\n", roomBatchrResult[key])
	}

	sendEnterRoomBatch.Base_R.Result = "ok"
	sendEnterRoomBatch.Payload = roomBatchrResult
	sendEnterRoomBatchJson, _ := json.Marshal(sendEnterRoomBatch)
	common.Sendmessage(connCore, sendEnterRoomBatchJson)

	for _, result := range roomBatchrResult {

		if result.Result != "" {
			continue
		}

		firstFromUuid := common.Getredisfirstenterroom(result.Roominfo.Roomcore.Roomuuid + userPlatform.Useruuid)
		if firstFromUuid == "" {
			continue
		}
		common.Deleteredisfirstenterroom(result.Roominfo.Roomcore.Roomuuid + userPlatform.Useruuid)

		// common.Essyslog("firstenterroom Roomname : "+result.Roominfo.Roomname, loginUuid, client.Userplatform.Useruuid)
		historyUuid := common.Getid().Hexstring()
		chatMessage := socket.Chatmessage{
			Historyuuid: historyUuid,
			From:        userPlatform,
			Stamp:       timeUnix,
			Message:     "enter room",
			Style:       "sys",
			Ip:          clientIp,
		}
		roomBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_ENTER_ROOM, Stamp: timeUnix}}
		roomBroadcast.Payload.Chatmessage = chatMessage
		roomBroadcast.Payload.Chattarget = result.Roominfo.Roomcore.Roomuuid
		roomBroadcastJson, _ := json.Marshal(roomBroadcast)
		common.Redispubroomsinfo(result.Roominfo.Roomcore.Roomuuid, roomBroadcastJson)

		common.Setredisroomlastmessage(result.Roominfo.Roomcore.Roomuuid, chatMessage)

		chatMessageHsitory := common.Chathistory{
			Historyuuid:    historyUuid,
			Chattarget:     result.Roominfo.Roomcore.Roomuuid,
			Myuuid:         userUuid,
			Myplatformuuid: userPlatform.Platformuuid,
			Myplatform:     userPlatform.Platform,
			Stamp:          timeUnix,
			Message:        chatMessage.Message,
			Style:          chatMessage.Style,
			Ip:             clientIp,
		}
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert(os.Getenv(result.Roominfo.Roomcore.Roomtype), string(chatMessageJson))
		if err != nil {
			common.Essyslog("Esinsert "+os.Getenv(result.Roominfo.Roomcore.Roomtype)+" err: "+err.Error(), loginUuid, userUuid)
		}

	}
	return nil
}
