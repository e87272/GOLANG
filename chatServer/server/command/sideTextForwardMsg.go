package command

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

func Sidetextforwardmsg(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendForwardMsg := socket.Cmd_r_forward_msg{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_SIDETEXT_FORWARD_MSG,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	maxLength := 200

	var packetForwardMsg socket.Cmd_c_forward_msg

	// log.Printf("Playersendforwardmsg packetForwardMsg : %+v\n", packetForwardMsg)
	err := json.Unmarshal([]byte(msg), &packetForwardMsg)
	if err != nil {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_JSON_ERROR", userUuid, err)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return err
	}
	sendForwardMsg.Base_R.Idem = packetForwardMsg.Base_C.Idem

	//禁止訪客發話
	if loginUuid == userPlatform.Useruuid && false {
		//block處理
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_GUEST", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	var forwardChatMessage socket.Chatmessage
	err = json.Unmarshal([]byte(packetForwardMsg.Payload.Forwardchatmessage), &forwardChatMessage)
	if err != nil {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_FORWARD_JSON_ERROR", userUuid, err)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return err
	}

	// log.Printf("Playersendforwardmsg forwardChatMessage : %+v\n", forwardChatMessage)

	if forwardChatMessage.Historyuuid == "" {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_FORWARD_MSG_NULL", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_IP_READ_ERROR", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if common.Isnewuser(userPlatform.Useruuid) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_IS_NEW_USER", userUuid, nil)
		sendForwardMsg.Base_R.Exp.Code = ""
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if common.Isspeakcd(userPlatform.Useruuid, timeUnix) {
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_SPEAK_CD", userUuid, nil)
		sendForwardMsg.Base_R.Exp.Code = ""
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)
		return nil
	}

	if utf8.RuneCountInString(packetForwardMsg.Payload.Message) > maxLength {
		//block處理
		sendForwardMsg.Base_R.Result = "err"
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_MSG_TOO_LONG", userUuid, nil)
		sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
		common.Sendmessage(connCore, sendForwardMsgJson)

		_, err = database.Exec(
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
		sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_CHAT_BLOCK", userUuid, nil)
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
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert("sidetextdirtywordhistory", string(chatMessageJson))
		if err != nil {
			common.Essyserrorlog("COMMAND_SIDETEXTFORWARDMSG_ES_INSERT_ERROR", userPlatform.Useruuid, err)
		}
	}

	targetPlatform, ok := common.Clientssidetextuserread(loginUuid, packetForwardMsg.Payload.Chattarget)
	if !ok {

		targetUser, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userUuid, packetForwardMsg.Payload.Chattarget)
		if !ok {
			sendForwardMsg.Base_R.Result = "err"
			sendForwardMsg.Base_R.Exp = exception
			sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
			common.Sendmessage(connCore, sendForwardMsgJson)
			return nil
		}

		sideTextUuid := common.Getid().Hexstring()
		targetPlatform.Userplatform = targetUser.Userplatform
		targetPlatform.Sidetextuuid = sideTextUuid
		forwardUser := userPlatform
		backwardUser := targetPlatform.Userplatform
		if userPlatform.Useruuid > packetForwardMsg.Payload.Chattarget {
			forwardUser = targetPlatform.Userplatform
			backwardUser = userPlatform
		}

		_, err = database.Exec(
			"INSERT INTO sideText (sideTextUuid, forwardUuid, backwardUuid,forwardPlatformUuid,forwardPlatform,backwardPlatformUuid,backwardPlatform) VALUES (?, ?, ?, ?, ?, ? , ? )",
			sideTextUuid,
			forwardUser.Useruuid,
			backwardUser.Useruuid,
			forwardUser.Platformuuid,
			forwardUser.Platform,
			backwardUser.Platformuuid,
			backwardUser.Platform,
		)

		if err != nil {
			sendForwardMsg.Base_R.Result = "err"
			sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_INSERT_CHATTARGET_ERROR", userUuid, err)
			sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
			common.Sendmessage(connCore, sendForwardMsgJson)
			return nil
		}

		sideTextMap, err := common.Querysidetextmap(userPlatform.Useruuid)
		if err != nil {
			sendForwardMsg.Base_R.Result = "err"
			sendForwardMsg.Base_R.Exp = common.Exception("COMMAND_SIDETEXTFORWARDMSG_QUERY_SIDETEXT_ERROR", userUuid, err)
			sendForwardMsgJson, _ := json.Marshal(sendForwardMsg)
			common.Sendmessage(connCore, sendForwardMsgJson)
			return nil
		}
		common.Clientssidetextinsert(loginUuid, sideTextMap)
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
	sendMsgBroadcast := socket.Cmd_b_player_speak{Base_B: socket.Base_B{Cmd: socket.CMD_B_SIDETEXT, Stamp: timeUnix}}
	sendMsgBroadcast.Payload.Chatmessage = chatMessage
	sendMsgBroadcast.Payload.Chattarget = packetForwardMsg.Payload.Chattarget
	sendMsgBroadcastJson, _ := json.Marshal(sendMsgBroadcast)
	common.Redispubsidetextdata(userPlatform.Useruuid, targetPlatform.Userplatform.Useruuid, sendMsgBroadcastJson)
	common.Setredissidetextlastmessage(targetPlatform.Sidetextuuid, chatMessage)

	chatMessageHsitory := common.Chathistory{
		Historyuuid:        historyUuid,
		Chattarget:         targetPlatform.Sidetextuuid,
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

	err = common.Esinsert(os.Getenv("sideText"), string(chatMessageJson))
	if err != nil {
		common.Essyserrorlog("COMMAND_SIDETEXTFORWARDMSG_ES_CHAT_HISTORY_INSERT_ERROR", userUuid, err)
		return err
	}

	return nil
}
