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

func Getnewsidetext(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendNewSideText := socket.Cmd_r_get_new_side_text{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_NEW_SIDETEXT,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetNewSidetext socket.Cmd_c_get_new_side_text

	if err := json.Unmarshal([]byte(msg), &packetNewSidetext); err != nil {
		sendNewSideText.Base_R.Result = "err"
		sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_JSON_ERROR", userUuid, err)
		sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
		common.Sendmessage(connCore, sendNewSideTextJson)
		return err
	}
	sendNewSideText.Base_R.Idem = packetNewSidetext.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendNewSideText.Base_R.Result = "err"
		sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_GUEST", userUuid, nil)
		sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
		common.Sendmessage(connCore, sendNewSideTextJson)
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
		boolQ.Filter(elastic.NewMatchQuery("chatTarget", sideText.Sidetextuuid))
		boolQ.Filter(elastic.NewRangeQuery("historyUuid").Gt(oldLastMessageUuid))
		searchResult, err := common.Elasticclient.Search(os.Getenv("sideText")).Query(boolQ).Sort("historyUuid", false).Do(context.Background())
		if err != nil {
			sendNewSideText.Base_R.Result = "err"
			sendNewSideText.Base_R.Exp = common.Exception("COMMAND_GETNEWSIDETEXT_SEARCH_ERROR", userUuid, err)
			sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
			common.Sendmessage(connCore, sendNewSideTextJson)
			return nil
		}

		var newSidetext socket.Newsidetext
		newSidetext.Targetuserplatform = sideText.Userplatform
		newSidetext.Newmessagecount = strconv.FormatInt(searchResult.Hits.TotalHits.Value, 10)
		newSidetext.Lastmessage = common.Hierarchysidetextlastmessage(loginUuid, userUuid, sideText.Sidetextuuid)

		newSidetextList = append(newSidetextList, newSidetext)
	}

	sendNewSideText.Base_R.Result = "ok"
	sendNewSideText.Payload.Newsidetextlist = newSidetextList
	sendNewSideTextJson, _ := json.Marshal(sendNewSideText)
	common.Sendmessage(connCore, sendNewSideTextJson)

	return nil
}
