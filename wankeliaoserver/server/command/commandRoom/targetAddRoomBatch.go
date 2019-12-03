package commandRoom

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../database"
	"../../socket"
)

func Targetaddroombatch(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendTargetAddRoomBatch := socket.Cmd_r_target_add_room_batch{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_TARGET_ADD_ROOM_BATCH,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetTargetAddRoomBatch socket.Cmd_c_target_add_room_batch
	err := json.Unmarshal([]byte(msg), &packetTargetAddRoomBatch)
	if err != nil {
		sendTargetAddRoomBatch.Base_R.Result = "err"
		sendTargetAddRoomBatch.Base_R.Exp = common.Exception("COMMAND_TARGETADDROOMBATCH_JSON_ERROR", userUuid, err)
		sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
		common.Sendmessage(connect, sendTargetAddRoomBatchJson)
		return err
	}
	sendTargetAddRoomBatch.Base_R.Idem = packetTargetAddRoomBatch.Base_C.Idem

	for _, roomCore := range packetTargetAddRoomBatch.Payload.Room {
		// roomCoreJson, _ := json.Marshal(roomCore)
		// common.Essyslog(string(roomCoreJson), loginUuid, userPlatform.Useruuid)

		if !common.Checkadmin(roomCore.Roomuuid, userPlatform.Useruuid, "AddPartner") {
			sendTargetAddRoomBatch.Base_R.Result = "err"
			sendTargetAddRoomBatch.Base_R.Exp = common.Exception("COMMAND_TARGETADDROOMBATCH_NOT_ADMIN", userPlatform.Useruuid, nil)
			sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
			common.Essyslog(string(sendTargetAddRoomBatchJson), loginUuid, userPlatform.Useruuid)
			common.Sendmessage(connect, sendTargetAddRoomBatchJson)
			return nil
		}
	}

	targetUserInfo, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetTargetAddRoomBatch.Payload.Targetuuid)
	if !ok {
		sendTargetAddRoomBatch.Base_R.Result = "err"
		sendTargetAddRoomBatch.Base_R.Exp = exception
		sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
		common.Essyslog(string(sendTargetAddRoomBatchJson), loginUuid, userPlatform.Useruuid)
		common.Sendmessage(connect, sendTargetAddRoomBatchJson)
		return nil
	}

	vipGroupMap := map[string]string{}
	vipGroupArray := strings.Split(targetUserInfo.Vipgroup, ",")
	for _, roomUuid := range vipGroupArray {
		vipGroupMap[roomUuid] = roomUuid
	}

	privateGroupMap := map[string]string{}
	privateGroupArray := strings.Split(targetUserInfo.Privategroup, ",")
	for _, roomUuid := range privateGroupArray {
		privateGroupMap[roomUuid] = roomUuid
	}

	roomCoreList := []socket.Roomcore{}
	for key, roomCore := range packetTargetAddRoomBatch.Payload.Room {
		// roomCoreJson, _ := json.Marshal(roomCore)
		// common.Essyslog(string(roomCoreJson), loginUuid, userPlatform.Useruuid)

		sendTargetAddRoomBatch.Payload = append(sendTargetAddRoomBatch.Payload, struct {
			Result   string          `json:"result"`
			Roomcore socket.Roomcore `json:"roomCore"`
		}{})
		sendTargetAddRoomBatch.Payload[key].Roomcore = roomCore

		if roomCore.Roomtype == "privateGroup" {

			_, ok := privateGroupMap[roomCore.Roomuuid]
			if ok {
				code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_TARGET_IN_ROOM", userPlatform.Useruuid, nil)
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

			userListName := roomCore.Roomtype + "UserList"
			uuid := common.Getid().Hexstring()
			_, err = database.Exec(
				"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
				uuid,
				roomCore.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)
			if err != nil {
				code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_INSERT_DB_ERROR", userPlatform.Useruuid, nil)
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

			log.Printf("INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (%+v, %+v, %+v, %+v)",
				uuid,
				roomCore.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)

			privateGroupMap[roomCore.Roomuuid] = roomCore.Roomuuid
			common.Setredisfirstenterroom(roomCore.Roomuuid+targetUserInfo.Userplatform.Useruuid, userPlatform.Useruuid)

			code := memberCount(userListName, roomCore.Roomuuid, userPlatform.Useruuid)
			if code != "" {
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

		} else if roomCore.Roomtype == "vipGroup" {

			_, ok := vipGroupMap[roomCore.Roomuuid]
			if ok {
				code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_TARGET_IN_ROOM", userPlatform.Useruuid, nil)
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

			userListName := roomCore.Roomtype + "UserList"
			uuid := common.Getid().Hexstring()
			_, err = database.Exec(
				"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
				uuid,
				roomCore.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)
			if err != nil {
				code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_INSERT_DB_ERROR", userPlatform.Useruuid, nil)
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

			log.Printf("INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (%+v, %+v, %+v, %+v)",
				uuid,
				roomCore.Roomuuid,
				targetUserInfo.Userplatform.Useruuid,
				"",
			)

			vipGroupMap[roomCore.Roomuuid] = roomCore.Roomuuid
			common.Setredisfirstenterroom(roomCore.Roomuuid+targetUserInfo.Userplatform.Useruuid, userPlatform.Useruuid)

			code := memberCount(userListName, roomCore.Roomuuid, userPlatform.Useruuid)
			if code != "" {
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}

		} else {
			code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_ROOM_TYPE_ERROR", userPlatform.Useruuid, nil)
			sendTargetAddRoomBatch.Payload[key].Result = code
			continue
		}

		roomCoreList = append(roomCoreList, roomCore)
	}
	sendTargetAddRoomBatch.Base_R.Result = "ok"
	sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
	common.Sendmessage(connect, sendTargetAddRoomBatchJson)

	if len(roomCoreList) > 0 {

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
		// common.Essyslog(string(pubDataJson), loginUuid, userPlatform.Useruuid)
		common.Redispubdata("sync", string(pubDataJson))

	}

	return nil
}

func memberCount(userListName string, roomUuid string, userUuid string) string {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	row := database.QueryRow("SELECT count(*) FROM users RIGHT JOIN "+userListName+" ON users.uuid="+userListName+".userUuid WHERE "+userListName+".roomUuid = ?",
		roomUuid,
	)
	memberCount := 0
	err := row.Scan(&memberCount)
	if err != nil {
		code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_QUERY_MEMBER_ERROR", userUuid, err)
		return code
	}

	common.Setredismembercount(roomUuid, memberCount)

	roomCountBroadcast := socket.Cmd_b_room_member_count{Base_B: socket.Base_B{Cmd: socket.CMD_B_ROOM_MEMBER_COUNT, Stamp: timeUnix}}
	roomCountBroadcast.Payload.Count = memberCount
	roomCountBroadcast.Payload.Roomuuid = roomUuid
	roomCountBroadcastJson, _ := json.Marshal(roomCountBroadcast)
	common.Redispubroomdata(roomUuid, roomCountBroadcastJson)
	return ""
}
