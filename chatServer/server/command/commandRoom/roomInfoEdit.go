package commandRoom

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"server/common"
	"server/database"
	"server/socket"
)

func Roominfoedit(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendRoomInfoEdit := socket.Cmd_r_room_info_edit{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_ROOM_INFO_EDIT,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	maxLength := 30

	var packetRoomInfoEdit socket.Cmd_c_room_info_edit

	if err := json.Unmarshal([]byte(msg), &packetRoomInfoEdit); err != nil {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_JSON_ERROR", userUuid, err)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return err
	}
	sendRoomInfoEdit.Base_R.Idem = packetRoomInfoEdit.Base_C.Idem

	if loginUuid == userUuid && false {
		//block處理
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_GUEST", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}

	if !common.Checkadmin("", userUuid, "RoomEdit") {
		//block處理
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_NOT_ADMIN", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}
	if packetRoomInfoEdit.Payload.Roomname == "" || utf8.RuneCountInString(packetRoomInfoEdit.Payload.Roomname) > maxLength {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_NAME_ERROR", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}

	isDirtyWord, _ := common.Matchdirtyword(packetRoomInfoEdit.Payload.Roomname, maxLength)
	if isDirtyWord {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_NAME_DIRTY_WORD", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}

	roomIconLink := ""
	if packetRoomInfoEdit.Payload.Roomicon != "" {
		//因為base64轉byte是4:3，檔案限制大小為3MB，所以阻擋要用4MB阻擋，先阻擋來減少解析時間
		if len(packetRoomInfoEdit.Payload.Roomicon) > 4*1024*1024 {
			sendRoomInfoEdit.Base_R.Result = "err"
			sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_ICON_TOO_LARGE", userUuid, nil)
			sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
			common.Sendmessage(connCore, sendRoomInfoEditJson)
			return nil
		}

		// 对上面的编码结果进行base64解码
		decodeBytes, err := base64.StdEncoding.DecodeString(packetRoomInfoEdit.Payload.Roomicon)

		if err != nil {
			sendRoomInfoEdit.Base_R.Result = "err"
			sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_ICON_ERROR", userUuid, err)
			sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
			common.Sendmessage(connCore, sendRoomInfoEditJson)
			return nil
		}

		thisTime := time.Now()
		fileName := "/" + common.Getid().Hexstring() + ".png"
		uploadPath := os.Getenv("uploadPath") + strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month()))
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		roomIconLink = "/" + strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month())) + fileName
		ioutil.WriteFile(uploadPath+fileName, decodeBytes, os.ModePerm)
	}

	roomInfo, ok := common.Roomsinforead(packetRoomInfoEdit.Payload.Roomcore.Roomuuid)

	if !ok {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_UUID_ERROR", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}

	var err error
	if roomIconLink != "" {
		_, err = database.Exec("UPDATE "+roomInfo.Roomcore.Roomtype+" SET roomName = ? , roomIcon = ? where roomUuid = ?", packetRoomInfoEdit.Payload.Roomname, roomIconLink, roomInfo.Roomcore.Roomuuid)
		roomInfo.Roomicon = roomIconLink
	} else {
		_, err = database.Exec("UPDATE "+roomInfo.Roomcore.Roomtype+" SET roomName = ? where roomUuid = ?", packetRoomInfoEdit.Payload.Roomname, roomInfo.Roomcore.Roomuuid)
	}
	if err != nil {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_UPDATE_GROUP_ERROR", userUuid, err)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connCore, sendRoomInfoEditJson)
		return nil
	}
	roomInfo.Roomname = packetRoomInfoEdit.Payload.Roomname

	common.Setredisroominfo(roomInfo.Roomcore.Roomuuid, roomInfo)
	common.Roomsinfoinsert(packetRoomInfoEdit.Payload.Roomcore.Roomuuid, roomInfo)

	sendRoomInfoEdit.Base_R.Result = "ok"
	sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
	common.Sendmessage(connCore, sendRoomInfoEditJson)

	targetUserAry := strings.Split(packetRoomInfoEdit.Payload.Targetuuid, ",")

	for _, targetUuid := range targetUserAry {
		targetUserInfo, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userUuid, targetUuid)
		if !ok {
			common.Essyserrorlog(exception.Message, "userUuid : "+userUuid+"  targetUserUuid : "+targetUserInfo.Userplatform.Useruuid, nil)
			continue
		}

		// log.Printf("Roominfoedit roomInfo : %+v\n", roomInfo)
		switch roomInfo.Roomcore.Roomtype {
		case "privateGroup", "vipGroup":
			// log.Printf("Roominfoedit Roomtype : %+v\n", roomInfo.Roomcore.Roomtype)
			ok, _ := common.Roominsertuser(userPlatform, targetUserInfo, roomInfo.Roomcore)
			if !ok {
				continue
			}
		default:
			common.Essyserrorlog("COMMAND_ROOMINFOEDIT_ROOM_TYPE_ERROR", userPlatform.Useruuid, nil)
			continue
		}

		roomCoreList := []socket.Roomcore{}
		roomCoreList = append(roomCoreList, socket.Roomcore{Roomuuid: roomInfo.Roomcore.Roomuuid, Roomtype: roomInfo.Roomcore.Roomtype})
		targetAddRoomMessage := socket.Cmd_b_target_add_room{Base_B: socket.Base_B{
			Cmd:   socket.CMD_B_NOTIFY_ENTER_ROOM,
			Stamp: timeUnix,
		}}
		targetAddRoomMessage.Payload = roomCoreList
		targetAddRoomMessageJson, _ := json.Marshal(targetAddRoomMessage)

		userMessage := common.Redispubsubuserdata{
			Useruuid: targetUserInfo.Userplatform.Useruuid,
			Datajson: string(targetAddRoomMessageJson),
		}
		userMessageJson, _ := json.Marshal(userMessage)

		//更新列表
		pubData := common.Syncdata{
			Synctype: "userInfoSyncAndEmit",
			Data:     string(userMessageJson),
		}
		pubDataJson, _ := json.Marshal(pubData)
		// common.Essyslog(string(pubDataJson), loginUuid, userUuid)
		common.Redispubdata("sync", string(pubDataJson))
	}

	//更新列表
	sendRoomInfoBroadcast := map[string]string{}
	sendRoomInfoBroadcast["roomType"] = roomInfo.Roomcore.Roomtype
	sendRoomInfoBroadcast["roomUuid"] = roomInfo.Roomcore.Roomuuid
	sendRoomInfoBroadcastJson, _ := json.Marshal(sendRoomInfoBroadcast)

	pubData := common.Syncdata{Synctype: "roomsInfoSync", Data: string(sendRoomInfoBroadcastJson)}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	return nil
}
