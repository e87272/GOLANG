package commandRoom

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic"

	"server/common"
	"server/socket"
)

func Getroomhistory(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendRoomHistory := socket.Cmd_r_get_chat_history{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_CHAT_HISTORY,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetRoomHistory socket.Cmd_c_get_room_chat_history

	if err := json.Unmarshal([]byte(msg), &packetRoomHistory); err != nil {
		sendRoomHistory.Base_R.Result = "err"
		sendRoomHistory.Base_R.Exp = common.Exception("COMMAND_GETROOMHISTORY_JSON_ERROR", userUuid, err)
		sendRoomHistoryJson, _ := json.Marshal(sendRoomHistory)
		common.Sendmessage(connCore, sendRoomHistoryJson)
		return err
	}

	sendRoomHistory.Base_R.Idem = packetRoomHistory.Base_C.Idem

	roomCore, ok := common.Clientsroomread(loginUuid, packetRoomHistory.Payload.Roomcore.Roomuuid)

	if !ok {
		//block處理
		sendRoomHistory.Base_R.Result = "err"
		sendRoomHistory.Base_R.Exp = common.Exception("COMMAND_GETROOMHISTORY_ROOM_UUID_ERROR", userUuid, nil)
		sendRoomHistoryJson, _ := json.Marshal(sendRoomHistory)
		common.Sendmessage(connCore, sendRoomHistoryJson)
		return nil
	}

	if packetRoomHistory.Payload.Historyuuid == "" {
		packetRoomHistory.Payload.Historyuuid = common.Getid().Hexstring()
	}

	//走SEEN後移除
	common.Setredisroomlastseen(roomCore.Roomtype+"_"+roomCore.Roomuuid+"_"+userPlatform.Useruuid, common.Getid().Hexstring())

	endStampHex := packetRoomHistory.Payload.Historyuuid
	// log.Printf("endStampHex : %+v\n", endStampHex)
	endStamp, _ := strconv.ParseInt(endStampHex, 16, 64)
	// log.Printf("endStamp : %+v\n", endStamp)
	rangeStamp := common.Timetoshift(common.Oncechathistorylong)
	// log.Printf("rangeStamp : %+v\n", rangeStamp)
	startStamp := endStamp - rangeStamp
	if startStamp < 0 {
		startStamp = 0
	}
	// log.Printf("startStamp : %+v\n", startStamp)
	startStampHex := common.Addzerohexstring(startStamp)
	// log.Printf("startStampHex : %+v\n", startStampHex)
	historyUuid := startStampHex

	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewMatchQuery("chatTarget", roomCore.Roomuuid))
	boolQ.Filter(elastic.NewRangeQuery("historyUuid").Gte(startStampHex).Lt(endStampHex))

	// Search with a term query
	searchResult, err := common.Elasticclient.Search(os.Getenv(roomCore.Roomtype)).Query(boolQ).Sort("historyUuid", false).Size(common.Maxchathistory).Do(context.Background()) // execute

	if err != nil {
		sendRoomHistory.Base_R.Result = "err"
		sendRoomHistory.Base_R.Exp = common.Exception("COMMAND_GETROOMHISTORY_ES_SEARCH_ERROR", userUuid, err)
		sendRoomHistoryJson, _ := json.Marshal(sendRoomHistory)
		common.Sendmessage(connCore, sendRoomHistoryJson)
		return err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		chatHistoryList := make([]socket.Chatmessage, 0, searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index
			// log.Printf("hit : %+v\n", hit)
			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var chatHistory common.Chathistory
			err := json.Unmarshal(*hit.Source, &chatHistory)
			if err != nil {
				// Deserialization failed
			}
			// Work with tweet
			// log.Printf("ChatMessage : %+v\n", chatHistory)

			fromUserPlatform := socket.Userplatform{
				Useruuid:     chatHistory.Myuuid,
				Platformuuid: chatHistory.Myplatformuuid,
				Platform:     chatHistory.Myplatform,
			}
			chatMessage := socket.Chatmessage{
				Historyuuid:        chatHistory.Historyuuid,
				From:               fromUserPlatform,
				Stamp:              chatHistory.Stamp,
				Message:            chatHistory.Message,
				Style:              chatHistory.Style,
				Ip:                 chatHistory.Ip,
				Forwardchatmessage: chatHistory.Forwardchatmessage,
			}

			chatHistoryList = append(chatHistoryList, chatMessage)

			historyUuid = chatHistory.Historyuuid
		}
		if searchResult.Hits.TotalHits < int64(common.Maxchathistory) {
			historyUuid = startStampHex
		}
		sendRoomHistory.Payload.Message = chatHistoryList

	} else {
		sendRoomHistory.Payload.Message = make([]socket.Chatmessage, 0)
	}

	sendRoomHistory.Result = "ok"
	sendRoomHistory.Payload.Historyuuid = historyUuid
	sendRoomHistory.Payload.Chattarget = roomCore.Roomuuid

	// log.Printf("sendRoomHistory : %+v\n", sendRoomHistory)

	sendRoomHistoryJson, _ := json.Marshal(sendRoomHistory)
	common.Sendmessage(connCore, sendRoomHistoryJson)

	return nil
}
