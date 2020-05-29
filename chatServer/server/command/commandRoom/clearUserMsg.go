package commandRoom

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"server/common"
	"server/socket"
	"github.com/olivere/elastic/v7"
)

func Clearusermsg(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendClearUserMsg := socket.Cmd_r_clear_user_msg{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_CLEAR_USER_MSG,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetClearUserMsg socket.Cmd_c_clear_user_msg

	err := json.Unmarshal([]byte(msg), &packetClearUserMsg)
	if err != nil {
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_JSON_ERROR", userUuid, err)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}
	sendClearUserMsg.Base_R.Idem = packetClearUserMsg.Base_C.Idem

	if loginUuid == userUuid && false {
		//block處理
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_GUEST", userUuid, nil)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetClearUserMsg.Payload.Roomuuid)
	if !ok {
		//block處理
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_ROOMUUID_ERROR", userUuid, nil)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}
	// log.Printf("Blockroomchat roomInfo : %+v\n", roomInfo)

	if !common.Checkadmin(packetClearUserMsg.Payload.Roomuuid, userPlatform.Useruuid, "Quiet") {
		//block處理
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_NOT_ADMIN", userUuid, nil)
		// log.Printf("sendClearUserMsg : %+v\n", sendClearUserMsg)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	blockUser, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetClearUserMsg.Payload.Targetuuid)

	if !ok {
		//block處理
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = exception
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	if common.Checktargetadmin(packetClearUserMsg.Payload.Roomuuid, packetClearUserMsg.Payload.Targetuuid, blockUser.Globalrole, "Quiet") {
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_TARGET_IS_ADMIN", userPlatform.Useruuid, nil)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewMatchQuery("chatTarget", roomInfo.Roomcore.Roomuuid))
	boolQ.Filter(elastic.NewMatchQuery("myUuid", userPlatform.Useruuid))
	_, err = common.Elasticclient.DeleteByQuery(os.Getenv(roomInfo.Roomcore.Roomtype)).Query(boolQ).Do(context.Background())
	if err != nil {
		sendClearUserMsg.Base_R.Result = "err"
		sendClearUserMsg.Base_R.Exp = common.Exception("COMMAND_CLEARUSERMSG_ES_DELETE_ERROR", userUuid, err)
		sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}

	sendClearUserMsg.Base_R.Result = "ok"
	sendChatblockJson, _ := json.Marshal(sendClearUserMsg)
	common.Sendmessage(connCore, sendChatblockJson)

	clearUserMsgData := common.Redispubsubclearusermsgdata{
		Roomuuid:   packetClearUserMsg.Payload.Roomuuid,
		Targetuuid: packetClearUserMsg.Payload.Targetuuid,
	}
	clearUserMsgDataJson, _ := json.Marshal(clearUserMsgData)
	pubData := common.Syncdata{
		Synctype: "clearUserMsg",
		Data:     string(clearUserMsgDataJson),
	}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	return nil
}
