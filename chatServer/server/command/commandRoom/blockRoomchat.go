package commandRoom

import (
	"encoding/json"

	"strconv"
	"time"

	"server/common"
	"server/database"
	"server/socket"
)

func Blockroomchat(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendChatblock := socket.Cmd_r_chatblock{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_CHATBLOCK,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	var packetChatBlock socket.Cmd_c_chatblock

	if err := json.Unmarshal([]byte(msg), &packetChatBlock); err != nil {
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_JSON_ERROR", userUuid, err)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}
	sendChatblock.Base_R.Idem = packetChatBlock.Base_C.Idem

	if loginUuid == userUuid && false {
		//block處理
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_GUEST", userUuid, nil)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	_, ok := common.Roomsinforead(packetChatBlock.Payload.Roomuuid)
	if !ok {
		//block處理
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_ROOMUUID_ERROR", userUuid, nil)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}
	// log.Printf("Blockroomchat Roomsinforead ok : %+v\n", ok)

	if !common.Checkadmin(packetChatBlock.Payload.Roomuuid, userPlatform.Useruuid, "Quiet") {
		//block處理
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_NOT_ADMIN", userUuid, nil)
		// log.Printf("sendChatblock : %+v\n", sendChatblock)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	// log.Printf("Blockroomchat Checkadmin ok : %+v\n", true)

	blockUser, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetChatBlock.Payload.Useruuid)

	if !ok {
		//block處理
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = exception
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	// log.Printf("Blockroomchat Hierarchytargetinfosearch blockUser : %+v\n", blockUser)

	if common.Checktargetadmin(packetChatBlock.Payload.Roomuuid, packetChatBlock.Payload.Useruuid, blockUser.Globalrole, "Quiet") {
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_TARGET_IS_ADMIN", userPlatform.Useruuid, nil)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return nil
	}

	// log.Printf("Blockroomchat Checktargetadmin ok : %+v\n", false)

	targetUuid := packetChatBlock.Payload.Roomuuid
	blockUserUuid := blockUser.Userplatform.Useruuid
	blockType := "room"
	blockTime, err := strconv.ParseInt(packetChatBlock.Payload.Blocktime, 10, 64)
	blockTime = blockTime*60*1000 + time.Now().UnixNano()/int64(time.Millisecond)
	if err != nil {
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_TIME_ERROR", userPlatform.Useruuid, err)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}

	var blockUuid = common.Getid().Hexstring()
	_, err = database.Exec(
		"DELETE FROM `chatBlock` WHERE blockUserUuid = ? and blocktarget = ? ",
		blockUserUuid,
		targetUuid,
	)
	if err != nil {
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_DB_DELETE_ERROR", userPlatform.Useruuid, err)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}

	// log.Printf("Blockroomchat DELETE : %+v\n", true)

	_, err = database.Exec(
		"INSERT INTO chatBlock (blockUuid, blockUserUuid, blocktarget , blockType, platformUuid, platform, timeStamp) VALUES (?, ? , ? , ? , ? , ? , ? )",
		blockUuid,
		blockUserUuid,
		targetUuid,
		blockType,
		blockUser.Userplatform.Platformuuid,
		blockUser.Userplatform.Platform,
		blockTime,
	)

	if err != nil {
		sendChatblock.Base_R.Result = "err"
		sendChatblock.Base_R.Exp = common.Exception("COMMAND_BLOCKROOMCHAT_DB_INSERT_ERROR", userPlatform.Useruuid, err)
		sendChatblockJson, _ := json.Marshal(sendChatblock)
		common.Sendmessage(connCore, sendChatblockJson)
		return err
	}

	// log.Printf("Blockroomchat INSERT : %+v\n", true)

	sendChatblock.Base_R.Result = "ok"
	sendChatblock.Payload.Useruuid = blockUserUuid
	sendChatblock.Payload.Roomuuid = targetUuid
	sendChatblockJson, _ := json.Marshal(sendChatblock)
	common.Sendmessage(connCore, sendChatblockJson)

	//更新列表
	pubData := common.Syncdata{Synctype: "blockSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	common.Pubsudoresult(packetChatBlock.Payload.Roomuuid, userPlatform, []string{"/su", "BU", blockUser.Userplatform.Useruuid, packetChatBlock.Payload.Blocktime}, blockUser.Userplatform, "Quiet")

	return nil
}
