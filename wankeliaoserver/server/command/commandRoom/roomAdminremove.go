package commandRoom

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"../../common"
	"../../database"
	"../../socket"
	"github.com/gorilla/websocket"
)

func Roomadminremove(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendRoomAdminRemove := socket.Cmd_r_room_admin_remove{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_ROOM_ADMIN_REMOVE,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetRoomAdminAdd socket.Cmd_c_room_admin_add

	if err := json.Unmarshal([]byte(msg), &packetRoomAdminAdd); err != nil {
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_JSON_ERROR", userUuid, err)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return err
	}
	sendRoomAdminRemove.Base_R.Idem = packetRoomAdminAdd.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_GUEST", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(packetRoomAdminAdd.Payload.Roomuuid)
	if !ok {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_ROOMUUID_ERROR", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}
	// log.Printf("ROOMADMINADD roomInfo : %+v\n", roomInfo)

	if !common.Checkadmin(roomInfo.Roomuuid, userPlatform.Useruuid, "AddAdmin") {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_NOT_ADMIN", userUuid, nil)
		// log.Printf("sendRoomAdminRemove : %+v\n", sendRoomAdminRemove)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	targetUser, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetRoomAdminAdd.Payload.Targetuuid)

	if !ok {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = exception
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	var roleUuid string
	row := database.QueryRow("SELECT roleUuid FROM role where roleName = ?", packetRoomAdminAdd.Payload.Role)
	err := row.Scan(&roleUuid)
	if err != nil {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_SELECT_ROLE_ERROR", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	var targetRoleSet string
	userListName := roomInfo.Roomtype + "UserList"
	row = database.QueryRow("SELECT roleSet FROM "+userListName+" where roomUuid = ? and userUuid = ?", roomInfo.Roomuuid, targetUser.Userplatform.Useruuid)
	err = row.Scan(&targetRoleSet)
	if err != nil {
		//block處理
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_SELECT_ROLE_ERROR", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	targetRoleArray := strings.Split(targetRoleSet, ",")
	targetRoleMap := map[string]bool{}
	for _, value := range targetRoleArray {
		targetRoleMap[value] = true
	}
	_, ok = targetRoleMap[roleUuid]
	if !ok {
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_USER_ROLE_ERROR", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	delete(targetRoleMap, roleUuid)
	targetRoleSet = ""
	for key, _ := range targetRoleMap {
		if targetRoleSet == "" {
			targetRoleSet = key
		} else {
			targetRoleSet = targetRoleSet + "," + key
		}
	}

	_, err = database.Exec("UPDATE "+userListName+" SET roleSet = ? where roomUuid = ? and userUuid = ?", targetRoleSet, roomInfo.Roomuuid, targetUser.Userplatform.Useruuid)
	if err != nil {
		sendRoomAdminRemove.Base_R.Result = "err"
		sendRoomAdminRemove.Base_R.Exp = common.Exception("COMMAND_ROOMADMINADD_UPDATE_ROLE_ERROR", userUuid, nil)
		sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
		common.Sendmessage(connect, sendRoomAdminRemoveJson)
		return nil
	}

	sendRoomAdminRemove.Base_R.Result = "ok"
	sendRoomAdminRemoveJson, _ := json.Marshal(sendRoomAdminRemove)
	common.Sendmessage(connect, sendRoomAdminRemoveJson)

	common.Queryroominfo(userUuid, roomInfo.Roomtype, roomInfo.Roomuuid)

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
