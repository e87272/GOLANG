package command

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../database"
	"../socket"
)

func Sidetextdelete(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendSidetextDelete := socket.Cmd_r_player_side_text_delete{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_SIDETEXT_DELETE,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetSideTextDelete socket.Cmd_c_side_text_delete

	if err := json.Unmarshal([]byte(msg), &packetSideTextDelete); err != nil {
		sendSidetextDelete.Base_R.Result = "err"
		sendSidetextDelete.Base_R.Exp = common.Exception("COMMAND_SIDETEXTDELETE_JSON_ERROR", userUuid, err)
		sendSidetextDeleteJson, _ := json.Marshal(sendSidetextDelete)
		common.Sendmessage(connect, sendSidetextDeleteJson)
		return err
	}
	sendSidetextDelete.Base_R.Idem = packetSideTextDelete.Base_C.Idem

	//禁止訪客發話
	if loginUuid == userPlatform.Useruuid {
		sendSidetextDelete.Base_R.Result = "err"
		sendSidetextDelete.Base_R.Exp = common.Exception("COMMAND_SIDETEXTDELETE_GUEST", userUuid, nil)
		sendSidetextDeleteJson, _ := json.Marshal(sendSidetextDelete)
		common.Sendmessage(connect, sendSidetextDeleteJson)
		return nil
	}

	if !common.Checkadmin("", userPlatform.Useruuid, "KillMessage") {
		//block處理
		sendSidetextDelete.Base_R.Result = "err"
		sendSidetextDelete.Base_R.Exp = common.Exception("COMMAND_SIDETEXTDELETE_KILLMESSAGE_NOT_ADMIN", userUuid, nil)
		sendSidetextDeleteJson, _ := json.Marshal(sendSidetextDelete)
		common.Sendmessage(connect, sendSidetextDeleteJson)
		return nil
	}

	targetPlatform, ok := common.Clientssidetextuserread(loginUuid, packetSideTextDelete.Payload)

	if !ok {
		//block處理
		sendSidetextDelete.Base_R.Result = "err"
		sendSidetextDelete.Base_R.Exp = common.Exception("COMMAND_SIDETEXTDELETE_NOT_SIDETEXT", userUuid, nil)
		sendSidetextDeleteJson, _ := json.Marshal(sendSidetextDelete)
		common.Sendmessage(connect, sendSidetextDeleteJson)
		return nil
	}

	common.Clientssidetextuserdelete(loginUuid, packetSideTextDelete.Payload)

	// log.Printf("Sidetextdelete targetPlatform.Sidetextuuid : %+v\n", targetPlatform.Sidetextuuid)

	_, err := database.Exec(
		"DELETE FROM `sideText` WHERE sideTextUuid = ?",
		targetPlatform.Sidetextuuid,
	)

	if err != nil {
		sendSidetextDelete.Base_R.Result = "err"
		sendSidetextDelete.Base_R.Exp = common.Exception("COMMAND_SIDETEXTDELETE_DELETE_ERROR", userUuid, nil)
		sendSidetextDeleteJson, _ := json.Marshal(sendSidetextDelete)
		common.Sendmessage(connect, sendSidetextDeleteJson)
		return nil
	}

	SendMsg := socket.Cmd_r_player_side_text_delete{Base_R: socket.Base_R{Cmd: socket.CMD_R_SIDETEXT_DELETE, Idem: packetSideTextDelete.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
	SendMsgJson, _ := json.Marshal(SendMsg)
	common.Sendmessage(connect, SendMsgJson)

	deleteMsg := socket.Cmd_b_sidetext_delete{Base_B: socket.Base_B{Cmd: socket.CMD_B_SIDETEXT_DELETE, Stamp: timeUnix}, Payload: packetSideTextDelete.Payload}
	deleteMsgJson, _ := json.Marshal(deleteMsg)

	sidetextData := common.Redispubsubsidetextdata{Useruuid: userPlatform.Useruuid, Targetuuid: targetPlatform.Userplatform.Useruuid, Datajson: string(deleteMsgJson)}
	sidetextDataJson, _ := json.Marshal(sidetextData)

	pubData := common.Syncdata{Synctype: "sideTextDeleteSync", Data: string(sidetextDataJson)}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	return nil
}
