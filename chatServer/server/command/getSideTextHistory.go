package command

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic"

	"../common"
	"../socket"
)

func Getsidetexthistory(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendSideTextHistory := socket.Cmd_r_get_side_text_history{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_SIDETEXT_HISTORY,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetSideTextHistory socket.Cmd_c_get_side_text_history

	if err := json.Unmarshal([]byte(msg), &packetSideTextHistory); err != nil {
		sendSideTextHistory.Base_R.Result = "err"
		sendSideTextHistory.Base_R.Exp = common.Exception("COMMAND_GETSIDETEXTHISTORY_JSON_ERROR", userUuid, err)
		sendSideTextHistoryJson, _ := json.Marshal(sendSideTextHistory)
		common.Sendmessage(connCore, sendSideTextHistoryJson)
		return err
	}
	sendSideTextHistory.Base_R.Idem = packetSideTextHistory.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendSideTextHistory.Base_R.Result = "err"
		sendSideTextHistory.Base_R.Exp = common.Exception("COMMAND_GETSIDETEXTHISTORY_GUEST", userUuid, nil)
		sendSideTextHistoryJson, _ := json.Marshal(sendSideTextHistory)
		common.Sendmessage(connCore, sendSideTextHistoryJson)
		return nil
	}

	if packetSideTextHistory.Payload.Historyuuid == "" {
		packetSideTextHistory.Payload.Historyuuid = common.Getid().Hexstring()
	}

	targetSideTextUser, ok := common.Clientssidetextuserread(loginUuid, packetSideTextHistory.Payload.Chattarget)
	// log.Printf("targetSideTextUser : %+v\n", targetSideTextUser)

	if !ok {
		//沒聊過歷史訊息為空
		sendSideTextHistory.Base_R.Result = "ok"
		sendSideTextHistory.Payload.Message = make([]socket.Chatmessage, 0)
		sendSideTextHistory.Payload.Historyuuid = common.Getid().Hexstring()
		sendSideTextHistory.Payload.Chattarget = targetSideTextUser.Userplatform.Useruuid
		sendSideTextHistoryJson, _ := json.Marshal(sendSideTextHistory)
		common.Sendmessage(connCore, sendSideTextHistoryJson)
		return nil
	}

	//走SEEN後移除
	common.Setredissidetextlastseen(targetSideTextUser.Sidetextuuid, common.Getid().Hexstring())

	endStampHex := packetSideTextHistory.Payload.Historyuuid
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
	boolQ.Filter(elastic.NewMatchQuery("chatTarget", targetSideTextUser.Sidetextuuid))
	boolQ.Filter(elastic.NewRangeQuery("historyUuid").Gte(startStampHex).Lt(endStampHex))

	// Search with a term query
	searchResult, err := common.Elasticclient.Search(os.Getenv("sideText")).Query(boolQ).Sort("historyUuid", false).Size(common.Maxchathistory).Do(context.Background()) // execute

	if err != nil {
		sendSideTextHistory.Base_R.Result = "err"
		sendSideTextHistory.Base_R.Exp = common.Exception("COMMAND_GETSIDETEXTHISTORY_ES_SEARCH_ERROR", userUuid, err)
		sendSideTextHistoryJson, _ := json.Marshal(sendSideTextHistory)
		common.Sendmessage(connCore, sendSideTextHistoryJson)
		return err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits.Value > 0 {
		chatHistoryList := make([]socket.Chatmessage, 0, searchResult.Hits.TotalHits.Value)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index
			// log.Printf("hit : %+v\n", hit)
			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var chatHistory common.Chathistory
			err := json.Unmarshal(hit.Source, &chatHistory)
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
		if searchResult.Hits.TotalHits.Value < int64(common.Maxchathistory) {
			historyUuid = startStampHex
		}
		sendSideTextHistory.Payload.Message = chatHistoryList

	} else {
		sendSideTextHistory.Payload.Message = make([]socket.Chatmessage, 0)
	}
	sendSideTextHistory.Result = "ok"
	sendSideTextHistory.Payload.Historyuuid = historyUuid
	sendSideTextHistory.Payload.Chattarget = targetSideTextUser.Userplatform.Useruuid

	sendSideTextHistoryJson, _ := json.Marshal(sendSideTextHistory)
	common.Sendmessage(connCore, sendSideTextHistoryJson)

	return nil
}
