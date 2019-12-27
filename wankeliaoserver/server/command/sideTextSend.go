package command

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"../common"
	"../database"
	"../socket"
)

func Sidetextsend(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendSidetext := socket.Cmd_r_player_side_text{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_PLAYER_SIDETEXT,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	maxLength := 200

	var packetSideText socket.Cmd_c_player_side_text

	if err := json.Unmarshal([]byte(msg), &packetSideText); err != nil {
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_JSON_ERROR", userUuid, err)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)
		return err
	}
	sendSidetext.Base_R.Idem = packetSideText.Base_C.Idem

	//禁止訪客發話
	if loginUuid == userPlatform.Useruuid {
		//block處理
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_GUEST", userUuid, nil)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)
		return nil
	}

	//自我私聊
	if packetSideText.Payload.Chattarget == userPlatform.Useruuid {
		//block處理
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_SIDE_TEXT_YOURSELF", userUuid, nil)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)
		return nil
	}

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_IP_READ_ERROR", userUuid, nil)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)
		return nil
	}

	if common.Isspeakcd(userPlatform.Useruuid, timeUnix) {
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_SPEAK_CD", userUuid, nil)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)
		return nil
	}

	if utf8.RuneCountInString(packetSideText.Payload.Message) > maxLength {
		sendSidetext.Base_R.Result = "err"
		sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_MSG_TOO_LONG", userUuid, nil)
		sendSidetextJson, _ := json.Marshal(sendSidetext)
		common.Sendmessage(connCore, sendSidetextJson)

		_, err := database.Exec(
			"DELETE FROM `chatBlock` WHERE blockUserUuid = ? and blocktarget = ? ",
			userPlatform.Useruuid,
			userPlatform.Useruuid,
		)
		if err != nil {
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_DELETE_CHATBLOCK_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
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
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_CHATBLOCK_INSERT_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
			return nil
		}

		//更新列表
		pubData := common.Syncdata{Synctype: "blockSync", Data: ""}
		pubDataJson, _ := json.Marshal(pubData)
		common.Redispubdata("sync", string(pubDataJson))

		return nil
	}

	targetPlatform, ok := common.Clientssidetextuserread(loginUuid, packetSideText.Payload.Chattarget)
	// log.Printf("targetPlatform : %+v\n", targetPlatform)
	if !ok {
		sideTextUuid := common.Getid().Hexstring()
		var platformUuid string
		var platform string

		row := database.QueryRow("select platformUuid, platform from users where uuid = ?", packetSideText.Payload.Chattarget)
		err := row.Scan(&platformUuid, &platform)

		if err == database.ErrNoRows {
			//block處理
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_SELECT_UUID_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
			return nil
		} else if err != nil {
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_SELECT_UUID_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
			return nil
		}

		targetPlatform.Userplatform = socket.Userplatform{Useruuid: packetSideText.Payload.Chattarget, Platformuuid: platformUuid, Platform: platform}
		targetPlatform.Sidetextuuid = sideTextUuid
		// log.Printf("!ok targetPlatform : %+v\n", targetPlatform)
		forwardUser := userPlatform
		backwardUser := targetPlatform.Userplatform
		if userPlatform.Useruuid > packetSideText.Payload.Chattarget {
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
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_INSERT_CHATTARGET_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
			return nil
		}

		sideTextMap, err := common.Querysidetextmap(userPlatform.Useruuid)
		if err != nil {
			sendSidetext.Base_R.Result = "err"
			sendSidetext.Base_R.Exp = common.Exception("COMMAND_SIDETEXTSEND_QUERY_SIDETEXT_ERROR", userUuid, err)
			sendSidetextJson, _ := json.Marshal(sendSidetext)
			common.Sendmessage(connCore, sendSidetextJson)
			return nil
		}
		common.Clientssidetextinsert(loginUuid, sideTextMap)
	}

	sideTextUuid := targetPlatform.Sidetextuuid
	historyUuid := common.Getid().Hexstring()

	isDirtyWord, clearMessage := common.Matchdirtyword(packetSideText.Payload.Message, maxLength)

	if isDirtyWord {

		chatMessageHsitory := common.Chathistory{
			Historyuuid:        historyUuid,
			Chattarget:         sideTextUuid,
			Myuuid:             userPlatform.Useruuid,
			Myplatformuuid:     userPlatform.Platformuuid,
			Myplatform:         userPlatform.Platform,
			Stamp:              timeUnix,
			Message:            packetSideText.Payload.Message,
			Style:              packetSideText.Payload.Style,
			Ip:                 clientIp,
			Forwardchatmessage: packetSideText.Payload.Forwardchatmessage,
		}
		// Index a second tweet (by string)
		chatMessageJson, _ := json.Marshal(chatMessageHsitory)

		err := common.Esinsert("sidetextdirtywordhistory", string(chatMessageJson))

		if err != nil {
			common.Essyserrorlog("COMMAND_SIDETEXTSEND_ES_DIRTYWORD_HISTORY_INSERT_ERROR", userUuid, err)
		}
	}

	sendSidetext.Base_R.Result = "ok"
	sendSidetextJson, _ := json.Marshal(sendSidetext)
	common.Sendmessage(connCore, sendSidetextJson)

	switch packetSideText.Payload.Style {
	case "url", "string":
	default:
		packetSideText.Payload.Style = "string"
	}

	Sidetextmessage := socket.Chatmessage{
		Historyuuid: historyUuid,
		From:        userPlatform,
		Stamp:       timeUnix,
		Message:     clearMessage,
		Style:       packetSideText.Payload.Style,
		Ip:          clientIp,
	}
	sendSideTextMessage := socket.Cmd_b_side_text{Base_B: socket.Base_B{Cmd: socket.CMD_B_SIDETEXT, Stamp: timeUnix}}
	sendSideTextMessage.Payload.Chatmessage = Sidetextmessage
	sendSideTextMessage.Payload.Chattarget = packetSideText.Payload.Chattarget
	sendSideTextMessageJson, _ := json.Marshal(sendSideTextMessage)

	// log.Printf("Redispubsidetextdata sendSideTextMessage : %+v\n", sendSideTextMessage)
	// log.Printf("Redispubsidetextdata userPlatform.Useruuid : %+v\n", userPlatform.Useruuid)
	// log.Printf("Redispubsidetextdata targetPlatform.Userplatform.Useruuid : %+v\n", targetPlatform.Userplatform.Useruuid)

	common.Redispubsidetextdata(userPlatform.Useruuid, targetPlatform.Userplatform.Useruuid, sendSideTextMessageJson)

	common.Setredissidetextlastmessage(sideTextUuid, Sidetextmessage)

	chatMessageHsitory := common.Chathistory{
		Historyuuid:    historyUuid,
		Chattarget:     sideTextUuid,
		Myuuid:         userPlatform.Useruuid,
		Myplatformuuid: userPlatform.Platformuuid,
		Myplatform:     userPlatform.Platform,
		Stamp:          timeUnix,
		Message:        clearMessage,
		Style:          packetSideText.Payload.Style,
		Ip:             clientIp,
	}
	// Index a second tweet (by string)
	chatMessageJson, _ := json.Marshal(chatMessageHsitory)

	err := common.Esinsert(os.Getenv("sideText"), string(chatMessageJson))

	if err != nil {
		common.Essyserrorlog("COMMAND_SIDETEXTSEND_ES_CHAT_HISTORY_INSERT_ERROR", userUuid, err)
		return err
	}
	return nil
}
