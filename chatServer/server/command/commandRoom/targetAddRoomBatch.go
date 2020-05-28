package commandRoom

import (
	"encoding/json"
	"strconv"
	"time"

	"server/common"
	"server/socket"
)

func Targetaddroombatch(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendTargetAddRoomBatch := socket.Cmd_r_target_add_room_batch{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_TARGET_ADD_ROOM_BATCH,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetTargetAddRoomBatch socket.Cmd_c_target_add_room_batch
	err := json.Unmarshal([]byte(msg), &packetTargetAddRoomBatch)
	if err != nil {
		sendTargetAddRoomBatch.Base_R.Result = "err"
		sendTargetAddRoomBatch.Base_R.Exp = common.Exception("COMMAND_TARGETADDROOMBATCH_JSON_ERROR", userUuid, err)
		sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
		common.Sendmessage(connCore, sendTargetAddRoomBatchJson)
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
			common.Sendmessage(connCore, sendTargetAddRoomBatchJson)
			return nil
		}
	}

	targetUserInfo, ok, exception := common.Hierarchytargetinfosearch(loginUuid, userPlatform.Useruuid, packetTargetAddRoomBatch.Payload.Targetuuid)
	if !ok {
		sendTargetAddRoomBatch.Base_R.Result = "err"
		sendTargetAddRoomBatch.Base_R.Exp = exception
		sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
		common.Essyslog(string(sendTargetAddRoomBatchJson), loginUuid, userPlatform.Useruuid)
		common.Sendmessage(connCore, sendTargetAddRoomBatchJson)
		return nil
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

		switch roomCore.Roomtype {
		case "privateGroup", "vipGroup":
			ok, code := common.Roominsertuser(userPlatform, targetUserInfo, roomCore)
			if !ok {
				sendTargetAddRoomBatch.Payload[key].Result = code
				continue
			}
		default:
			code := common.Essyserrorlog("COMMAND_TARGETADDROOMBATCH_ROOM_TYPE_ERROR", userPlatform.Useruuid, nil)
			sendTargetAddRoomBatch.Payload[key].Result = code
			continue
		}

		roomCoreList = append(roomCoreList, roomCore)
	}
	sendTargetAddRoomBatch.Base_R.Result = "ok"
	sendTargetAddRoomBatchJson, _ := json.Marshal(sendTargetAddRoomBatch)
	common.Sendmessage(connCore, sendTargetAddRoomBatchJson)

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
