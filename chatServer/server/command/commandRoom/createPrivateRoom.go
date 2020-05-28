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

func Createprivateroom(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendCreatePrivateRoom := socket.Cmd_r_create_private_room{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_CREATE_PRIVATE_ROOM,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	maxLength := 30

	var packetCreatePrivateRoom socket.Cmd_c_create_private_room

	err := json.Unmarshal([]byte(msg), &packetCreatePrivateRoom)
	if err != nil {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_JSON_ERROR", userUuid, err)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return err
	}
	sendCreatePrivateRoom.Base_R.Idem = packetCreatePrivateRoom.Base_C.Idem

	if loginUuid == userUuid && false {
		//block處理
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_GUEST", userUuid, nil)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	clientIp, ok := common.Iplistread(loginUuid)
	if !ok {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_IP_READ_ERROR", userUuid, nil)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	if !common.Checkadmin("", userUuid, "CreateGroup") {
		//block處理
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_NOT_ADMIN", userUuid, nil)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	if packetCreatePrivateRoom.Payload.Roomname == "" || utf8.RuneCountInString(packetCreatePrivateRoom.Payload.Roomname) > maxLength {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_ROOM_NAME_ERROR", userUuid, nil)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	isDirtyWord, _ := common.Matchdirtyword(packetCreatePrivateRoom.Payload.Roomname, maxLength)
	if isDirtyWord {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_DIRTY_WORD", userUuid, nil)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	roomIconLink := ""
	if packetCreatePrivateRoom.Payload.Roomicon != "" {
		//因為base64轉byte是4:3，檔案限制大小為3MB，所以阻擋要用4MB阻擋，先阻擋來減少解析時間
		if len(packetCreatePrivateRoom.Payload.Roomicon) > 4*1024*1024 {
			sendCreatePrivateRoom.Base_R.Result = "err"
			sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_ROOM_ICON_TOO_LARGE", userUuid, nil)
			sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
			common.Sendmessage(connCore, sendCreatePrivateRoomJson)
			return nil
		}

		// 对上面的编码结果进行base64解码
		decodeBytes, err := base64.StdEncoding.DecodeString(packetCreatePrivateRoom.Payload.Roomicon)

		if err != nil {
			sendCreatePrivateRoom.Base_R.Result = "err"
			sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_ROOM_ICON_ERROR", userUuid, err)
			sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
			common.Sendmessage(connCore, sendCreatePrivateRoomJson)
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

	roomUuid := common.Getid().Hexstring()
	_, err = database.Exec(
		"INSERT INTO `privateGroup` (roomUuid, roomName, roomIcon, owner) VALUES (?, ?, ?, ?)",
		roomUuid,
		packetCreatePrivateRoom.Payload.Roomname,
		roomIconLink,
		userUuid,
	)
	if err != nil {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_INSERT_DB_ERROR", userUuid, err)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	chatMessage := socket.Chatmessage{
		Historyuuid: roomUuid,
		From:        userPlatform,
		Stamp:       timeUnix,
		Message:     "create room",
		Style:       "sys",
		Ip:          clientIp,
	}
	common.Setredisroomlastmessage(roomUuid, chatMessage)

	chatMessageHsitory := common.Chathistory{
		Historyuuid:    roomUuid,
		Chattarget:     roomUuid,
		Myuuid:         userUuid,
		Myplatformuuid: userPlatform.Platformuuid,
		Myplatform:     userPlatform.Platform,
		Stamp:          timeUnix,
		Message:        "create room",
		Style:          "sys",
		Ip:             clientIp,
	}
	chatMessageHsitoryJson, _ := json.Marshal(chatMessageHsitory)
	err = common.Esinsert(os.Getenv("privateGroup"), string(chatMessageHsitoryJson))
	if err != nil {
		common.Essyslog("Esinsert "+os.Getenv("privateGroup")+" err: "+err.Error(), loginUuid, userUuid)
	}

	roleUuid := os.Getenv("roleUuidAdmin")

	uuid := common.Getid().Hexstring()
	_, err = database.Exec(
		"INSERT INTO `privateGroupUserList` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
		uuid,
		roomUuid,
		userUuid,
		roleUuid,
	)

	if err != nil {
		sendCreatePrivateRoom.Base_R.Result = "err"
		sendCreatePrivateRoom.Base_R.Exp = common.Exception("COMMAND_CREATEPRIVATEROOM_INSERT_USER_LIST_ERROR", userUuid, err)
		sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
		common.Sendmessage(connCore, sendCreatePrivateRoomJson)
		return nil
	}

	adminSet := map[string]string{}
	adminSet[userUuid] = roleUuid
	adminSetJson, _ := json.Marshal(adminSet)
	roomInfo := socket.Roominfo{
		Roomcore: socket.Roomcore{
			Roomuuid: roomUuid,
			Roomtype: "privateGroup",
		},
		Roomname:      packetCreatePrivateRoom.Payload.Roomname,
		Roomicon:      roomIconLink,
		Adminset:      string(adminSetJson),
		Ownerplatform: userPlatform,
	}

	common.Setredisroominfo(roomUuid, roomInfo)

	userInfo, _ := common.Usersinforead(userUuid)
	if userInfo.Privategroup == "" {
		userInfo.Privategroup = roomUuid
	} else {
		userInfo.Privategroup = userInfo.Privategroup + "," + roomUuid
	}

	common.Usersinfoinsert(userUuid, userInfo)

	sendCreatePrivateRoom.Base_R.Result = "ok"
	sendCreatePrivateRoom.Payload.Roomuuid = roomUuid
	sendCreatePrivateRoom.Payload.Roomtype = "privateGroup"
	sendCreatePrivateRoomJson, _ := json.Marshal(sendCreatePrivateRoom)
	common.Sendmessage(connCore, sendCreatePrivateRoomJson)

	//自己也要收到71預防多登狀況
	roomCoreList := []socket.Roomcore{}
	roomCoreList = append(roomCoreList, socket.Roomcore{Roomuuid: roomUuid, Roomtype: "privateGroup"})
	targetAddRoomMessage := socket.Cmd_b_target_add_room{Base_B: socket.Base_B{
		Cmd:   socket.CMD_B_NOTIFY_ENTER_ROOM,
		Stamp: timeUnix,
	}}
	targetAddRoomMessage.Payload = roomCoreList
	targetAddRoomMessageJson, _ := json.Marshal(targetAddRoomMessage)

	userMessage := common.Redispubsubuserdata{
		Useruuid: userUuid,
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

	// log.Printf("Createprivateroom pubDataJson : %+v\n", string(pubDataJson))

	targetUserAry := strings.Split(packetCreatePrivateRoom.Payload.Targetuuid, ",")

	for _, targetUuid := range targetUserAry {
		targetUserInfo, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userUuid, targetUuid)
		if !ok {
			common.Essyserrorlog(exception.Message, "userUuid : "+userUuid+"  targetUserUuid : "+targetUserInfo.Userplatform.Useruuid, nil)
			continue
		}

		ok, _ = common.Roominsertuser(userPlatform, targetUserInfo, roomInfo.Roomcore)
		if !ok {
			continue
		}

		roomCoreList := []socket.Roomcore{}
		roomCoreList = append(roomCoreList, socket.Roomcore{Roomuuid: roomUuid, Roomtype: "privateGroup"})
		targetAddRoomMessage := socket.Cmd_b_target_add_room{Base_B: socket.Base_B{
			Cmd:   socket.CMD_B_NOTIFY_ENTER_ROOM,
			Stamp: timeUnix,
		}}
		targetAddRoomMessage.Payload = roomCoreList
		targetAddRoomMessageJson, _ := json.Marshal(targetAddRoomMessage)

		userMessage := common.Redispubsubuserdata{
			Useruuid: targetUuid,
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

		// log.Printf("Createprivateroom pubDataJson : %+v\n", string(pubDataJson))
	}

	return nil
}
