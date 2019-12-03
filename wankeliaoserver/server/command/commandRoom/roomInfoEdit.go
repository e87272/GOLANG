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

	"../../common"
	"../../database"
	"../../socket"
	"github.com/gorilla/websocket"
)

func Roominfoedit(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendRoomInfoEdit := socket.Cmd_r_room_info_edit{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_ROOM_INFO_EDIT,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	maxLength := 30

	var packetRoomInfoEdit socket.Cmd_c_room_info_edit

	if err := json.Unmarshal([]byte(msg), &packetRoomInfoEdit); err != nil {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_JSON_ERROR", userUuid, err)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return err
	}
	sendRoomInfoEdit.Base_R.Idem = packetRoomInfoEdit.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_GUEST", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}

	if !common.Checkadmin("", userUuid, "RoomEdit") {
		//block處理
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_NOT_ADMIN", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}
	if packetRoomInfoEdit.Payload.Roomname == "" || utf8.RuneCountInString(packetRoomInfoEdit.Payload.Roomname) > maxLength {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_NAME_ERROR", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}

	isDirtyWord, _ := common.Matchdirtyword(packetRoomInfoEdit.Payload.Roomname, maxLength)
	if isDirtyWord {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_NAME_DIRTY_WORD", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}

	roomIconLink := ""
	if packetRoomInfoEdit.Payload.Roomicon != "" {
		//因為base64轉byte是4:3，檔案限制大小為3MB，所以阻擋要用4MB阻擋，先阻擋來減少解析時間
		if len(packetRoomInfoEdit.Payload.Roomicon) > 4*1024*1024 {
			sendRoomInfoEdit.Base_R.Result = "err"
			sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_ICON_TOO_LARGE", userUuid, nil)
			sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
			common.Sendmessage(connect, sendRoomInfoEditJson)
			return nil
		}

		// 对上面的编码结果进行base64解码
		decodeBytes, err := base64.StdEncoding.DecodeString(packetRoomInfoEdit.Payload.Roomicon)

		if err != nil {
			sendRoomInfoEdit.Base_R.Result = "err"
			sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_ICON_ERROR", userUuid, err)
			sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
			common.Sendmessage(connect, sendRoomInfoEditJson)
			return nil
		}

		thisTime := time.Now()
		fileName := "/" + common.Getid().Hexstring() + ".png"
		uploadPath := os.Getenv("uploadPath") + strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month()))
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		roomIconLink = strconv.Itoa(thisTime.Year()) + strconv.Itoa(int(thisTime.Month())) + fileName
		ioutil.WriteFile(uploadPath+fileName, decodeBytes, os.ModePerm)
	}

	roomInfo, ok := common.Roomsinforead(packetRoomInfoEdit.Payload.Roomcore.Roomuuid)

	if !ok {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMINFOEDIT_ROOM_UUID_ERROR", userUuid, nil)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}

	var err error
	if roomIconLink != "" {
		_, err = database.Exec("UPDATE "+roomInfo.Roomtype+" SET roomName = ? , roomIcon = ? where roomUuid = ?", packetRoomInfoEdit.Payload.Roomname, roomIconLink, roomInfo.Roomuuid)
		roomInfo.Roomicon = roomIconLink
	} else {
		_, err = database.Exec("UPDATE "+roomInfo.Roomtype+" SET roomName = ? where roomUuid = ?", packetRoomInfoEdit.Payload.Roomname, roomInfo.Roomuuid)
	}
	if err != nil {
		sendRoomInfoEdit.Base_R.Result = "err"
		sendRoomInfoEdit.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_UPDATE_GROUP_ERROR", userUuid, err)
		sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
		common.Sendmessage(connect, sendRoomInfoEditJson)
		return nil
	}
	roomInfo.Roomname = packetRoomInfoEdit.Payload.Roomname

	common.Setredisroominfo(roomInfo.Roomuuid, roomInfo)
	common.Roomsinfoinsert(packetRoomInfoEdit.Payload.Roomcore.Roomuuid, roomInfo)

	sendRoomInfoEdit.Base_R.Result = "ok"
	sendRoomInfoEditJson, _ := json.Marshal(sendRoomInfoEdit)
	common.Sendmessage(connect, sendRoomInfoEditJson)

	targetUserAry := strings.Split(packetRoomInfoEdit.Payload.Targetuuid, ",")

	for _, targetUuid := range targetUserAry {
		targetUserInfo, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userUuid, targetUuid)
		if !ok {
			common.Essyserrorlog(exception.Message, "userUuid : "+userUuid+"  targetUserUuid : "+targetUserInfo.Userplatform.Useruuid, nil)
			continue
		}

		if roomInfo.Roomtype == "privateGroup" {

			privateGroupMap := map[string]string{}
			privateGroupArray := strings.Split(targetUserInfo.Privategroup, ",")
			for _, roomUuid := range privateGroupArray {
				privateGroupMap[roomUuid] = roomUuid
			}

			privateGroupMapJson, _ := json.Marshal(privateGroupMap)
			common.Essyslog("Roominfoedit privateGroupMap : "+string(privateGroupMapJson), loginUuid, userPlatform.Useruuid)

			_, ok := privateGroupMap[roomInfo.Roomuuid]
			if ok {
				common.Essyserrorlog("COMMAND_ROOMINFOEDIT_TARGET_IN_ROOM", userPlatform.Useruuid, nil)
				continue
			}

			userListName := roomInfo.Roomtype + "UserList"
			uuid := common.Getid().Hexstring()
			_, err = database.Exec(
				"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
				uuid,
				roomInfo.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)
			if err != nil {
				common.Essyserrorlog("COMMAND_ROOMINFOEDIT_INSERT_DB_ERROR", userPlatform.Useruuid, nil)
				continue
			}

			privateGroupMap[roomInfo.Roomuuid] = roomInfo.Roomuuid
			common.Setredisfirstenterroom(roomInfo.Roomuuid+targetUserInfo.Userplatform.Useruuid, userPlatform.Useruuid)

		} else if roomInfo.Roomtype == "vipGroup" {

			vipGroupMap := map[string]string{}
			vipGroupArray := strings.Split(targetUserInfo.Vipgroup, ",")
			for _, roomUuid := range vipGroupArray {
				vipGroupMap[roomUuid] = roomUuid
			}

			vipGroupMapJson, _ := json.Marshal(vipGroupMap)
			common.Essyslog("Roominfoedit vipGroupMap : "+string(vipGroupMapJson), loginUuid, userPlatform.Useruuid)

			_, ok := vipGroupMap[roomInfo.Roomuuid]
			if ok {
				common.Essyserrorlog("COMMAND_ROOMINFOEDIT_TARGET_IN_ROOM", userPlatform.Useruuid, nil)
				continue
			}

			userListName := roomInfo.Roomtype + "UserList"
			uuid := common.Getid().Hexstring()
			_, err = database.Exec(
				"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
				uuid,
				roomInfo.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)
			if err != nil {
				common.Essyserrorlog("COMMAND_ROOMINFOEDIT_INSERT_DB_ERROR", userPlatform.Useruuid, nil)
				continue
			}

			vipGroupMap[roomInfo.Roomuuid] = roomInfo.Roomuuid
			common.Setredisfirstenterroom(roomInfo.Roomuuid+targetUserInfo.Userplatform.Useruuid, userPlatform.Useruuid)

		} else {
			common.Essyserrorlog("COMMAND_ROOMINFOEDIT_ROOM_TYPE_ERROR", userPlatform.Useruuid, nil)
			continue
		}

		roomCoreList := []socket.Roomcore{}
		roomCoreList = append(roomCoreList, socket.Roomcore{Roomuuid: roomInfo.Roomuuid, Roomtype: roomInfo.Roomtype})
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
	sendRoomInfoBroadcast["roomType"] = roomInfo.Roomtype
	sendRoomInfoBroadcast["roomUuid"] = roomInfo.Roomuuid
	sendRoomInfoBroadcastJson, _ := json.Marshal(sendRoomInfoBroadcast)

	pubData := common.Syncdata{Synctype: "roomsInfoSync", Data: string(sendRoomInfoBroadcastJson)}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	return nil
}
