package shell

import (
	"encoding/json"

	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"server/common"
	"server/database"
	"server/socket"
)

func shelldemo(shellCmd []string, connCore common.Conncore, userPlatform socket.Userplatform, packetSendShell socket.Cmd_c_player_send_shell, timeUnix string) error {

	switch shellCmd[1] {

	case "DSG": //demoSendGift

		if len(shellCmd) < 2 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		historyUuid := common.Getid().Hexstring()
		timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		// roominfo := common.Roomsinforead(packetSendShell.Payload.Roominfo.Roomcore.Roomuuid)
		chatMessage := socket.Chatmessage{
			Historyuuid: historyUuid,
			From:        userPlatform,
			Stamp:       timeUnix,
			Message:     "https://smzb2.zanxingtv.com//img//external//beststartv//gift//gif//appreciate.gif",
			Style:       "send gift",
		}
		sendGiftBroadcast := socket.Cmd_b_send_gift{Base_B: socket.Base_B{Cmd: socket.CMD_B_SEND_GIFT, Stamp: timeUnix}}
		sendGiftBroadcast.Payload.Chatmessage = chatMessage
		sendGiftBroadcast.Payload.Chattarget = packetSendShell.Payload.Chattarget
		sendGiftBroadcastJson, _ := json.Marshal(sendGiftBroadcast)

		// log.Printf("sendGiftBroadcastJson : %+v\n", sendGiftBroadcastJson)

		common.Redispubroomdata(packetSendShell.Payload.Chattarget, sendGiftBroadcastJson)

		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return nil

	case "DSB": //demoSubscription

		if len(shellCmd) < 2 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		historyUuid := common.Getid().Hexstring()
		timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		chatMessage := socket.Chatmessage{
			Historyuuid: historyUuid,
			From:        userPlatform,
			Stamp:       timeUnix,
			Message:     "subscription",
			Style:       "subscription",
		}
		subscriptionBroadcast := socket.Cmd_b_subscription{Base_B: socket.Base_B{Cmd: socket.CMD_B_SUBSCRIPTION, Stamp: timeUnix}}
		subscriptionBroadcast.Payload.Chatmessage = chatMessage
		subscriptionBroadcast.Payload.Chattarget = packetSendShell.Payload.Chattarget
		subscriptionBroadcastJson, _ := json.Marshal(subscriptionBroadcast)

		// log.Printf("subscriptionBroadcastJson : %+v\n", subscriptionBroadcastJson)

		common.Redispubroomdata(packetSendShell.Payload.Chattarget, subscriptionBroadcastJson)

		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return nil

	case "DAD": //demoAddadmin

		if len(shellCmd) < 4 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		roomInfo, _ := common.Roomsinforead(packetSendShell.Payload.Chattarget)
		roomType := roomInfo.Roomcore.Roomtype
		roomUuid := roomInfo.Roomcore.Roomuuid
		role := "admin"
		addUuid := shellCmd[3]

		row := database.QueryRow("select uuid from users where uuid = ?", addUuid)
		err := row.Scan(&addUuid)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SELECT_USER_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		var adminSetJson string
		row = database.QueryRow("select adminSet from "+roomType+" where roomUuid = ?", roomUuid)
		err = row.Scan(&adminSetJson)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SELECT_ADMINSET_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		var adminSet map[string]string
		err = json.Unmarshal([]byte(adminSetJson), &adminSet)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_JSON_ADMINSET_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}
		userRole, ok := adminSet[addUuid]
		adminSetRole := ""
		if ok {
			userRoleAry := strings.Split(userRole, ",")
			roleAry := strings.Split(role, ",")
			roleMap := make(map[string]bool)
			for _, val := range userRoleAry {
				roleMap[val] = true
			}
			for _, val := range roleAry {
				roleMap[val] = true
			}
			for key := range roleMap {
				adminSetRole = adminSetRole + key + ","
			}
			if adminSetRole != "" {
				adminSetRole = adminSetRole[0 : len(adminSetRole)-1]
			}
			adminSet[addUuid] = adminSetRole
		} else {
			adminSet[addUuid] = role
		}

		//D則刪除
		if shellCmd[2] == "d" {
			delete(adminSet, addUuid)
		}

		adminSetByte, _ := json.Marshal(adminSet)
		adminSetJson = string(adminSetByte)

		_, err = database.Exec("UPDATE "+roomType+" SET adminSet = ? where roomUuid = ?", adminSetJson, roomUuid)
		if err != nil {
			// log.Printf("Playersendmsg err : %+v\n", err)
			return nil
		}

		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		// log.Printf("adminSetJson : %+v\n", adminSetJson)

		//更新列表
		data := map[string]string{}
		data["roomType"] = roomType
		data["roomUuid"] = roomUuid
		dataJson, _ := json.Marshal(data)
		pubData := common.Syncdata{Synctype: "roomsInfoSync", Data: string(dataJson)}
		pubDataJson, _ := json.Marshal(pubData)
		common.Redispubdata("sync", string(pubDataJson))
		return nil
	case "DCN": //changeNickname

		argument := regexp.MustCompile(" +-").Split(strings.Trim(packetSendShell.Payload.Message, " "), -1)

		if len(shellCmd) < 3 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		var platformUuid = userPlatform.Platformuuid
		for i := 1; i < len(argument); i++ {
			value := regexp.MustCompile(" +").Split(argument[i], 2)
			if len(value) != 2 {
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connCore, SendMsgJson)
				return nil
			}
			switch value[0] {
			case "u":

				row := database.QueryRow("select platformUuid from users where uuid = ?", value[1])
				err := row.Scan(&platformUuid)
				if err != nil {
					SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SELECT_USER_ERROR", userPlatform.Useruuid, err)}}
					SendMsgJson, _ := json.Marshal(SendMsg)
					common.Sendmessage(connCore, SendMsgJson)
					return nil
				}

				break
			default:
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SHELL_ERROR", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connCore, SendMsgJson)
				return nil
			}
		}
		var data = url.Values{"nickname": {shellCmd[2]}, "uuid": {platformUuid}}
		req, err := http.NewRequest("PATCH", "https://mmt.zanxingbctv.com/v0.3/auth/user/nickname", strings.NewReader(data.Encode()))
		if os.Getenv("environmentId") == "Online" {
			req, err = http.NewRequest("PATCH", "https://wklapi.zanstartv.com/v0.3/auth/user/nickname", strings.NewReader(data.Encode()))
		}
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PATCH_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_REQUEST_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return nil

	case "DPS": //訊息輸出

		if len(shellCmd) < 4 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return nil
		}

		second, err := strconv.ParseFloat(shellCmd[2], 64)

		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_TIME_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return err
		}

		frequency, err := strconv.ParseFloat(shellCmd[3], 64)

		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_TIME_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return err
		}

		var duration = int(1000 * 1 / frequency)
		var i int64
		for i = 0; i < int64(second*frequency); i++ {
			historyUuid := common.Getid().Hexstring()
			timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			message := timeUnix + "-" + strconv.FormatInt(i, 10)
			roomInfo, _ := common.Roomsinforead(packetSendShell.Payload.Chattarget)
			chatMessage := socket.Chatmessage{
				Historyuuid: historyUuid,
				From:        userPlatform,
				Stamp:       timeUnix,
				Message:     message,
				Style:       packetSendShell.Payload.Style,
			}
			sendMsgBroadcast := socket.Cmd_b_player_speak{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_SPEAK, Stamp: timeUnix}}
			sendMsgBroadcast.Payload.Chatmessage = chatMessage
			sendMsgBroadcast.Payload.Chattarget = roomInfo.Roomcore.Roomuuid
			sendMsgBroadcastJson, _ := json.Marshal(sendMsgBroadcast)

			common.Redispubroomdata(roomInfo.Roomcore.Roomuuid, sendMsgBroadcastJson)

			chatMessageHsitory := common.Chathistory{
				Historyuuid:    historyUuid,
				Chattarget:     roomInfo.Roomcore.Roomuuid,
				Myuuid:         userPlatform.Useruuid,
				Myplatformuuid: userPlatform.Platformuuid,
				Myplatform:     userPlatform.Platform,
				Stamp:          timeUnix,
				Message:        message,
				Style:          packetSendShell.Payload.Style,
			}
			// Index a second tweet (by string)
			chatMessageJson, _ := json.Marshal(chatMessageHsitory)

			// log.Printf("chatMessageHsitory : %+v\n", chatMessageHsitory)

			err = common.Esinsert(os.Getenv(roomInfo.Roomcore.Roomtype), string(chatMessageJson))
			time.Sleep(time.Duration(duration) * time.Millisecond)
			// time.Sleep(time.Duration(250) * time.Millisecond)
		}
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return nil

	case "C8763": //重啟server
		var releasseMap = map[string]string{}
		releasseMap["ReleaseId"] = os.Getenv("releaseId")
		releasseMap["EnvironmentId"] = os.Getenv("environmentId")
		releasseMapJson, _ := json.Marshal(releasseMap)
		var url = "http://192.168.20.169/api/Spaces-1/deployments"
		req, err := http.NewRequest("POST", url, strings.NewReader(string(releasseMapJson)))
		if err != nil {
			// log.Printf("C8763 err : %+v\n", err)
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_OCTOPUS_URL_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connCore, SendMsgJson)
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Octopus-ApiKey", "API-JWFWV0KGDF2AGXATCYWUNKT0LW")
		_, err = http.DefaultClient.Do(req)
		// log.Printf("C8763 resp : %+v\n", resp)
		// log.Printf("C8763 err : %+v\n", err)
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return nil
	default:
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendShell.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_SHELLDEMO_SHELL_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}
}
