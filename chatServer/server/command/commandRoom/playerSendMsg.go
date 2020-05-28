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

func Playersendmsg(connCore common.Conncore, msg []byte, loginUuid string) error {

	stamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeUnix := strconv.FormatInt(stamp, 10)

	sendPlayerMsg := socket.Cmd_r_player_send_msg{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_SEND_MSG,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid

	maxLength := 200

	var packetSendMsg socket.Cmd_c_player_send_msg

	if err := json.Unmarshal([]byte(msg), &packetSendMsg); err != nil {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_JSON_ERROR", userUuid, err)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return err
	}
	sendPlayerMsg.Base_R.Idem = packetSendMsg.Base_C.Idem

	//禁止訪客發話
	if loginUuid == userPlatform.Useruuid {
		//block處理
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_GUEST", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_IP_READ_ERROR", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	if common.Isnewuser(userPlatform.Useruuid) {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_IS_NEW_USER", userUuid, nil)
		sendPlayerMsg.Base_R.Exp.Code = ""
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	if common.Isspeakcd(userPlatform.Useruuid, timeUnix) {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_SPEAK_CD", userUuid, nil)
		sendPlayerMsg.Base_R.Exp.Code = ""
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	if utf8.RuneCountInString(packetSendMsg.Payload.Message) > maxLength {
		//block處理
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_MSG_TOO_LONG", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)

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

	if common.Checkblock(loginUuid, packetSendMsg.Payload.Chattarget, userPlatform.Useruuid) {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_CHAT_BLOCK", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	switch packetSendMsg.Payload.Style {
	case "url":
		if !common.Checkadmin(packetSendMsg.Payload.Chattarget, userPlatform.Useruuid, "Hyperlink") {
			packetSendMsg.Payload.Style = "string"
		}
	case "string":
	default:
		packetSendMsg.Payload.Style = "string"
	}

	isDirtyWord, clearMessage := common.Matchdirtyword(packetSendMsg.Payload.Message, maxLength)

	if isDirtyWord {

		chatMessageHsitory := common.Chathistory{
			Historyuuid:        common.Getid().Hexstring(),
			Chattarget:         packetSendMsg.Payload.Chattarget,
			Myuuid:             userPlatform.Useruuid,
			Myplatformuuid:     userPlatform.Platformuuid,
			Myplatform:         userPlatform.Platform,
			Stamp:              timeUnix,
			Message:            packetSendMsg.Payload.Message,
			Style:              packetSendMsg.Payload.Style,
			Ip:                 clientIp,
			Forwardchatmessage: packetSendMsg.Payload.Forwardchatmessage,
		}
		// Index a second tweet (by string)
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert("roomdirtywordhistory", string(chatMessageJson))

		if err != nil {
			common.Essyserrorlog("COMMAND_PLAYERSENDMSG_ES_INSERT_ERROR", userPlatform.Useruuid, err)
		}
	}

	if !common.Checkinroom(packetSendMsg.Payload.Chattarget, loginUuid) {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_NOT_IN_ROOM", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetSendMsg.Payload.Chattarget)
	if !ok {
		sendPlayerMsg.Base_R.Result = "err"
		sendPlayerMsg.Base_R.Exp = common.Exception("COMMAND_PLAYERSENDMSG_NOT_IN_ROOM", userUuid, nil)
		sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
		common.Sendmessage(connCore, sendPlayerMsgJson)
		return nil
	}

	sendPlayerMsg.Base_R.Result = "ok"
	sendPlayerMsgJson, _ := json.Marshal(sendPlayerMsg)
	common.Sendmessage(connCore, sendPlayerMsgJson)

	historyUuid := common.Getid().Hexstring()
	chatMessage := socket.Chatmessage{
		Historyuuid: historyUuid,
		From:        userPlatform,
		Stamp:       timeUnix,
		Message:     clearMessage,
		Style:       packetSendMsg.Payload.Style,
		Ip:          clientIp,
	}

	if roomInfo.Roomcore.Roomtype != "liveGroup" {
		common.Setredisroomlastmessage(roomInfo.Roomcore.Roomuuid, chatMessage)
	}

	chatMessageHsitory := common.Chathistory{
		Historyuuid:    historyUuid,
		Chattarget:     packetSendMsg.Payload.Chattarget,
		Myuuid:         userUuid,
		Myplatformuuid: userPlatform.Platformuuid,
		Myplatform:     userPlatform.Platform,
		Stamp:          timeUnix,
		Message:        chatMessage.Message,
		Style:          chatMessage.Style,
		Ip:             clientIp,
	}
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	err := common.Esinsert(os.Getenv(roomInfo.Roomcore.Roomtype), string(chatMessageJson))
	if err != nil {
		common.Essyslog("Esinsert "+os.Getenv(roomInfo.Roomcore.Roomtype)+" err: "+err.Error(), loginUuid, userUuid)
	}

	//訊息過久直接略過
	if time.Now().UnixNano()/int64(time.Millisecond)-stamp > common.Packetdroptime {
		timeUnix := time.Now().UnixNano() / int64(time.Millisecond)
		common.Essyslog(string(chatMessageJson), "timeUnix : "+strconv.FormatInt(timeUnix, 10)+" packetStamp : "+strconv.FormatInt(timeUnix, 10), "Playersendmsg")
		return nil
	}

	sendMsgBroadcast := socket.Cmd_b_player_speak{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_SPEAK, Stamp: timeUnix}}
	sendMsgBroadcast.Payload.Chatmessage = chatMessage
	sendMsgBroadcast.Payload.Chattarget = packetSendMsg.Payload.Chattarget
	sendMsgBroadcastJson, _ := json.Marshal(sendMsgBroadcast)

	common.Redispubroomdata(packetSendMsg.Payload.Chattarget, sendMsgBroadcastJson)

	return nil
}
