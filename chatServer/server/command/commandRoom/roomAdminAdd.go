package commandRoom

import (
	"encoding/json"
	"strconv"
	"time"

	"../../common"
	"../../database"
	"../../socket"
)

func Roomadminadd(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendRoomAdminAdd := socket.Cmd_r_room_admin_add{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_ROOM_ADMIN_ADD,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetRoomAdminAdd socket.Cmd_c_room_admin_add

	if err := json.Unmarshal([]byte(msg), &packetRoomAdminAdd); err != nil {
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_JSON_ERROR", userUuid, err)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return err
	}
	sendRoomAdminAdd.Base_R.Idem = packetRoomAdminAdd.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_GUEST", userUuid, nil)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetRoomAdminAdd.Payload.Roomuuid)
	if !ok {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_ROOMUUID_ERROR", userUuid, nil)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}
	// log.Printf("ROOMADMINADD roomInfo : %+v\n", roomInfo)

	if !common.Checkadmin(roomInfo.Roomcore.Roomuuid, userPlatform.Useruuid, "AddAdmin") {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_NOT_ADMIN", userUuid, nil)
		// log.Printf("sendRoomAdminAdd : %+v\n", sendRoomAdminAdd)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}

	targetUser, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetRoomAdminAdd.Payload.Targetuuid)

	if !ok {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = exception
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}

	var roleUuid string
	row := database.QueryRow("SELECT roleUuid FROM role where roleName = ?", packetRoomAdminAdd.Payload.Role)
	err := row.Scan(&roleUuid)
	if err != nil {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_SELECT_ROLE_ERROR", userUuid, nil)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}

	var targetRoleSet string
	userListName := roomInfo.Roomcore.Roomtype + "UserList"
	row = database.QueryRow("SELECT roleSet FROM "+userListName+" where roomUuid = ? and userUuid = ?", roomInfo.Roomcore.Roomuuid, targetUser.Userplatform.Useruuid)
	err = row.Scan(&targetRoleSet)
	if err == nil {
		targetRoleSet = targetRoleSet + "," + roleUuid

		_, err = database.Exec("UPDATE "+userListName+" SET roleSet = ? where roomUuid = ? and userUuid = ?", targetRoleSet, roomInfo.Roomcore.Roomuuid, targetUser.Userplatform.Useruuid)
		if err != nil {
			sendRoomAdminAdd.Base_R.Result = "err"
			sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_UPDATE_ROLE_ERROR", userUuid, nil)
			sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
			common.Sendmessage(connCore, sendRoomAdminAddJson)
			return nil
		}

	} else if err == database.ErrNoRows {
		uuid := common.Getid().Hexstring()
		_, err = database.Exec(
			"INSERT INTO `"+userListName+"` (uuid, roomUuid, userUuid, roleSet) VALUES (?, ?, ?, ?)",
			uuid,
			roomInfo.Roomcore.Roomuuid,
			userUuid,
			roleUuid,
		)
		if err != nil {
			sendRoomAdminAdd.Base_R.Result = "err"
			sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_INSERT_ROLE_ERROR", userUuid, nil)
			sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
			common.Sendmessage(connCore, sendRoomAdminAddJson)
			return nil
		}
	} else if err != nil {
		//block處理
		sendRoomAdminAdd.Base_R.Result = "err"
		sendRoomAdminAdd.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_SELECT_ROLE_ERROR", userUuid, nil)
		sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
		common.Sendmessage(connCore, sendRoomAdminAddJson)
		return nil
	}

	sendRoomAdminAdd.Base_R.Result = "ok"
	sendRoomAdminAddJson, _ := json.Marshal(sendRoomAdminAdd)
	common.Sendmessage(connCore, sendRoomAdminAddJson)

	common.Queryroominfo(userUuid, roomInfo.Roomcore.Roomtype, roomInfo.Roomcore.Roomuuid)

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
