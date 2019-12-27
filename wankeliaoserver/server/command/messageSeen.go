package command

import (
	"encoding/json"
	"strconv"
	"time"

	"../common"
	"../socket"
)

func Messageseen(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendMessageSeen := socket.Cmd_r_message_seen{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_MESSAGE_SEEN,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetMessageSeen socket.Cmd_c_message_seen

	if err := json.Unmarshal([]byte(msg), &packetMessageSeen); err != nil {
		sendMessageSeen.Base_R.Result = "err"
		sendMessageSeen.Base_R.Exp = common.Exception("COMMAND_MESSAGESEEN_JSON_ERROR", userUuid, err)
		sendMessageSeenJson, _ := json.Marshal(sendMessageSeen)
		common.Sendmessage(connCore, sendMessageSeenJson)
		return err
	}
	sendMessageSeen.Base_R.Idem = packetMessageSeen.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendMessageSeen.Base_R.Result = "err"
		sendMessageSeen.Base_R.Exp = common.Exception("COMMAND_MESSAGESEEN_GUEST", userUuid, nil)
		sendMessageSeenJson, _ := json.Marshal(sendMessageSeen)
		common.Sendmessage(connCore, sendMessageSeenJson)
		return nil
	}

	historyUuid := common.Getid().Hexstring()
	switch packetMessageSeen.Payload.Chattype {
	case "room":

		roomInfo, ok := common.Roomsinforead(packetMessageSeen.Payload.Chattargetuuid)
		if !ok {
			//block處理
			messageSeen := socket.Cmd_r_message_seen{Base_R: socket.Base_R{Cmd: socket.CMD_R_MESSAGE_SEEN, Idem: packetMessageSeen.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("COMMAND_MESSAGESEEN_TARGET_ROOM_UUID_ERROR", userPlatform.Useruuid, nil)}}
			messageSeenJson, _ := json.Marshal(messageSeen)
			common.Sendmessage(connCore, messageSeenJson)
			return nil
		}

		if roomInfo.Roomcore.Roomtype == "liveGroup" {
			//block處理
			sendMessageSeen.Base_R.Result = "err"
			sendMessageSeen.Base_R.Exp = common.Exception("COMMAND_MESSAGESEEN_TARGET_ROOM_TYPE_ERROR", userUuid, nil)
			sendMessageSeenJson, _ := json.Marshal(sendMessageSeen)
			common.Sendmessage(connCore, sendMessageSeenJson)
			return nil
		}
		common.Setredisroomlastseen(roomInfo.Roomcore.Roomtype+"_"+roomInfo.Roomcore.Roomuuid+"_"+userPlatform.Useruuid, historyUuid)
		break
	case "sideText":
		targetPlatform, ok := common.Clientssidetextuserread(loginUuid, packetMessageSeen.Payload.Chattargetuuid)

		if !ok {
			//block處理
			sendMessageSeen.Base_R.Result = "err"
			sendMessageSeen.Base_R.Exp = common.Exception("COMMAND_MESSAGESEEN_TARGET_SIDE_TEXT_UUID_ERROR", userUuid, nil)
			sendMessageSeenJson, _ := json.Marshal(sendMessageSeen)
			common.Sendmessage(connCore, sendMessageSeenJson)
			return nil
		}
		common.Setredissidetextlastseen(targetPlatform.Sidetextuuid, historyUuid)

		beSeenMessage := socket.Cmd_b_message_be_seen{Base_B: socket.Base_B{Cmd: socket.CMD_B_MESSAGE_BE_SEEN, Stamp: timeUnix}}
		beSeenMessage.Payload.Chattarget = userPlatform.Useruuid
		beSeenMessage.Payload.Historyuuid = historyUuid
		beSeenMessageJson, _ := json.Marshal(beSeenMessage)

		userMessage := common.Redispubsubuserdata{Useruuid: targetPlatform.Userplatform.Useruuid, Datajson: string(beSeenMessageJson)}
		userMessageJson, _ := json.Marshal(userMessage)
		common.Redispubdata("user", string(userMessageJson))

		break
	}

	sendMessageSeen.Base_R.Result = "ok"
	sendMessageSeenJson, _ := json.Marshal(sendMessageSeen)
	common.Sendmessage(connCore, sendMessageSeenJson)

	return nil
}
