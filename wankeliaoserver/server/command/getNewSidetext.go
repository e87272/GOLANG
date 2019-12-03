package command

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"../common"
	"../socket"
)

func Getnewsidetext(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendNewSideText := socket.Cmd_r_get_new_side_text{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_NEW_SIDETEXT,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetNewSidetext socket.Cmd_c_get_new_side_text

	if err := json.Unmarshal([]byte(msg), &packetNewSidetext); err != nil {
		sendNewSideText.Base_R.Result = "err"
		sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_JSON_ERROR", userUuid, err)
		sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
		common.Sendmessage(connect, sendNewSideTextJson)
		return err
	}
	sendNewSideText.Base_R.Idem = packetNewSidetext.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendNewSideText.Base_R.Result = "err"
		sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_GUEST", userUuid, nil)
		sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
		common.Sendmessage(connect, sendNewSideTextJson)
		return nil
	}

	var newSidetextList []socket.Newsidetext

	sideTextMap, ok := common.Clientssidetextmapread(loginUuid)
	if !ok {
		sideTextMap, _ := common.Querysidetextmap(userPlatform.Useruuid)
		common.Clientssidetextinsert(loginUuid, sideTextMap)
	}
	for _, sideText := range sideTextMap {
		oldLastMessageUuid := common.Getredissidetextlastseen(sideText.Sidetextuuid)
		if oldLastMessageUuid == "" {
			oldLastMessageUuid = common.Getid().Hexstring()
		}
		boolQ := elastic.NewBoolQuery()
		boolQ.Must(elastic.NewMatchQuery("chatTarget", sideText.Sidetextuuid))
		boolQ.Filter(elastic.NewRangeQuery("historyUuid").Gt(oldLastMessageUuid))
		searchResult, err := common.Elasticclient.Search(os.Getenv("sideText")).Query(boolQ).Sort("historyUuid", false).Do(context.Background())
		if err != nil {
			sendNewSideText.Base_R.Result = "err"
			sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_SEARCH_ERROR", userUuid, err)
			sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
			common.Sendmessage(connect, sendNewSideTextJson)
			return nil
		}

		var newSidetext socket.Newsidetext
		newSidetext.Targetuserplatform = sideText.Userplatform
		newSidetext.Message = []socket.Chatmessage{}
		newSidetext.Lastmessage = common.Getredissidetextlastmessage(sideText.Sidetextuuid)

		for _, hit := range searchResult.Hits.Hits {
			var chatHistory common.Chathistory
			_ = json.Unmarshal(hit.Source, &chatHistory)

			var chatMessage socket.Chatmessage
			chatMessage.Historyuuid = chatHistory.Historyuuid
			chatMessage.From.Useruuid = chatHistory.Myuuid
			chatMessage.From.Platformuuid = chatHistory.Myplatformuuid
			chatMessage.From.Platform = chatHistory.Myplatform
			chatMessage.Stamp = chatHistory.Stamp
			chatMessage.Message = chatHistory.Message
			chatMessage.Style = chatHistory.Style
			newSidetext.Message = append(newSidetext.Message, chatMessage)
		}

		newSidetextList = append(newSidetextList, newSidetext)
	}

	sendNewSideText.Base_R.Result = "ok"
	sendNewSideText.Payload.Newsidetextlist = newSidetextList
	sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
	common.Sendmessage(connect, sendNewSideTextJson)

	return nil
}
