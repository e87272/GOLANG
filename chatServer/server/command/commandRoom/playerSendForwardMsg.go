package commandRoom

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"server/common"
	"server/database"
	"server/socket"
)

func Playersendforwardmsg(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendForwardMsg := socket.Cmd_r_forward_msg{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_FORWARD_MSG,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	maxLength := 200

	var packetForwardMsg socket.Cmd_c_forward_msg

	// log.Printf("Playersendforwardmsg packetForwardMsg : %+v\n", packetForwardMsg)

	if err := json.Unmarshal([]byte(msg), &packetForwardMsg); err != nil {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_JSON_ERROR", userUuid, err)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return err
	}
	sendForwardMsg.Base_R.Idem = packetForwardMsg.Base_C.Idem

	//禁止訪客發話
	if loginUuid == userPlatform.Useruuid && false {
		//block處理
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_GUEST", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	var forwardChatMessage socket.Chatmessage
	if err := json.Unmarshal([]byte(packetForwardMsg.Payload.Forwardchatmessage), &forwardChatMessage); err != nil {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_FORWARD_JSON_ERROR", userUuid, err)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return err
	}

	// log.Printf("Playersendforwardmsg forwardChatMessage : %+v\n", forwardChatMessage)

	if forwardChatMessage.Historyuuid == "" {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_FORWARD_MSG_NULL", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_IP_READ_ERROR", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if common.Isnewuser(userPlatform.Useruuid) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_IS_NEW_USER", userUuid, nil)
		sendForwardMsg.Base_R.Exp.Code = ""
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if common.Isspeakcd(userPlatform.Useruuid, timeUnix) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_SPEAK_CD", userUuid, nil)
		sendForwardMsg.Base_R.Exp.Code = ""
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if utf8.RuneCountInString(packetForwardMsg.Payload.Message) > maxLength {
		//block處理
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_MSG_TOO_LONG", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)

		_, err := database.Exec(
			"DELETE FROM `chatBlock` WHERE blockUserUuid = ? and blocktarget = ? ",
			userPlatform.Useruuid,
			userPlatform.Useruuid,
		)
		if err != nil {
			common.Essyslog("chatBlock DELETE err: "+err.Error(), loginUuid, userUuid)
			return nil
		}
		_, err = database.Exec(
			"INSERT INTO chatBlock (blockUuid, blockUserUuid, blocktarget , blockType, platformUuid, platform, timeStamp) VALUES (?, ? , ? , ? , ? , ? , ? )",
			common.Getid().Hexstring(),
			userPlatform.Useruuid,
			userPlatform.Useruuid,
			"user",
			userPlatform.Platformuuid,
			userPlatform.Platform,
			15*60*1000+time.Now().UnixNano()/int64(time.Millisecond),
		)

		if err != nil {
			common.Essyslog("chatBlock INSERT err: "+err.Error(), loginUuid, userUuid)
			return nil
		}

		_, err = database.Exec(
			"INSERT INTO chatIpBlock (blockUuid,blockip,timeStamp) VALUES (?, ? , ? )",
			common.Getid().Hexstring(),
			clientIp,
			15*60*1000+time.Now().UnixNano()/int64(time.Millisecond),
		)

		if err != nil {
			common.Essyslog("chatIpBlock INSERT err: "+err.Error(), loginUuid, userUuid)
			// return nil
		}

		//更新列表
		pubData := common.Syncdata{Synctype: "blockSync", Data: ""}
		pubDataJson, _ := json.Marshal(pubData)
		common.Redispubdata("sync", string(pubDataJson))

		return nil
	}

	if common.Checkblock(loginUuid, packetForwardMsg.Payload.Chattarget, userPlatform.Useruuid) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_CHAT_BLOCK", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	switch packetForwardMsg.Payload.Style {
	case "url":
		if !common.Checkadmin(packetForwardMsg.Payload.Chattarget, userPlatform.Useruuid, "Hyperlink") {
			packetForwardMsg.Payload.Style = "string"
		}
	case "string":
	default:
		packetForwardMsg.Payload.Style = "string"
	}

	isDirtyWord, clearMessage := common.Matchdirtyword(packetForwardMsg.Payload.Message, maxLength)

	if isDirtyWord {

		chatMessageHsitory := common.Chathistory{
			Historyuuid:        common.Getid().Hexstring(),
			Chattarget:         packetForwardMsg.Payload.Chattarget,
			Myuuid:             userPlatform.Useruuid,
			Myplatformuuid:     userPlatform.Platformuuid,
			Myplatform:         userPlatform.Platform,
			Stamp:              timeUnix,
			Message:            packetForwardMsg.Payload.Message,
			Style:              packetForwardMsg.Payload.Style,
			Ip:                 clientIp,
			Forwardchatmessage: packetForwardMsg.Payload.Forwardchatmessage,
		}
		// Index a second tweet (by string)
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert("roomdirtywordhistory", string(chatMessageJson))

		if err != nil {
			common.Essyserrorlog("COMMAND_PLAYERSENDFORWARDMSG_ES_INSERT_ERROR", userPlatform.Useruuid, err)
		}
	}

	if !common.Checkinroom(packetForwardMsg.Payload.Chattarget, loginUuid) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_NOT_IN_ROOM", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetForwardMsg.Payload.Chattarget)
	if !ok {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDFORWARDMSG_NOT_IN_ROOM", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	sendForwardMsg.Base_R.Result = "ok"
	sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
	common.Sendmessage(connCore, sendForwardMsgJson)

	historyUuid := common.Getid().Hexstring()
	chatMessage := socket.Chatmessage{
		Historyuuid:        historyUuid,
		From:               userPlatform,
		Stamp:              timeUnix,
		Message:            clearMessage,
		Style:              packetForwardMsg.Payload.Style,
		Ip:                 clientIp,
		Forwardchatmessage: packetForwardMsg.Payload.Forwardchatmessage,
	}
	sendMsgBroadcast := socket.Cmd_b_player_speak{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_SPEAK, Stamp: timeUnix}}
	sendMsgBroadcast.Payload.Chatmessage = chatMessage
	sendMsgBroadcast.Payload.Chattarget = packetForwardMsg.Payload.Chattarget
	sendMsgBroadcastJson, _ := json.Marshal(sendMsgBroadcast)

	common.Redispubroomdata(packetForwardMsg.Payload.Chattarget, sendMsgBroadcastJson)

	if roomInfo.Roomcore.Roomtype != "liveGroup" {
		common.Setredisroomlastmessage(roomInfo.Roomcore.Roomuuid, chatMessage)
	}

	chatMessageHsitory := common.Chathistory{
		Historyuuid:        historyUuid,
		Chattarget:         packetForwardMsg.Payload.Chattarget,
		Myuuid:             userUuid,
		Myplatformuuid:     userPlatform.Platformuuid,
		Myplatform:         userPlatform.Platform,
		Stamp:              timeUnix,
		Message:            chatMessage.Message,
		Style:              chatMessage.Style,
		Ip:                 clientIp,
		Forwardchatmessage: packetForwardMsg.Payload.Forwardchatmessage,
	}
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	err := common.Esinsert(os.Getenv(roomInfo.Roomcore.Roomtype), string(chatMessageJson))
	if err != nil {
		common.Essyslog("Esinsert "+os.Getenv(roomInfo.Roomcore.Roomtype)+" err: "+err.Error(), loginUuid, userUuid)
	}

	return nil
}
