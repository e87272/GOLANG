package shell

import (
	"encoding/json"

	"regexp"
	"strconv"
	"strings"
	"time"

	"../common"
	"../database"
	"../socket"
)

func blockUser(connCore common.Conncore, userPlatform socket.Userplatform, packetSendMsg socket.Cmd_c_player_send_shell, timeUnix string) error {

	argument := regexp.MustCompile(" +-").Split(strings.Trim(packetSendMsg.Payload.Message, " "), -1)
	shellCmd := regexp.MustCompile(" +").Split(argument[0], -1)

	if len(shellCmd) < 4 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_PARAMETER_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	// log.Printf("blockuser  shellCmd: %+v\n", shellCmd)
	blockUserUuid := shellCmd[2]
	roomUuid := packetSendMsg.Payload.Chattarget
	if roomUuid == "" || blockUserUuid == "" {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_UUID_NULL", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	var blockUserPlatform socket.Userplatform
	row := database.QueryRow("select uuid,platformUuid,platform from users where uuid = ? ", blockUserUuid)
	err := row.Scan(&blockUserPlatform.Useruuid, &blockUserPlatform.Platformuuid, &blockUserPlatform.Platform)

	if err != nil {
		// log.Printf("Pubsudoresult select targetUuid err : %+v\n", err)
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_USER_UUID_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	roomInfo, ok := common.Roomsinforead(roomUuid)
	if !ok {
		// log.Printf("Pubsudoresult select targetUuid err : %+v\n", err)
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_ROOM_UUID_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	if !common.Checkadmin(roomUuid, userPlatform.Useruuid, "Quiet") {
		//block處理
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_NOT_ADMIN", userPlatform.Useruuid, nil)}}
		// log.Printf("sendChatblock : %+v\n", sendChatblock)
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}

	if common.Checkadmin(roomUuid, blockUserUuid, "Quiet") {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_TARGET_IS_ADMIN", userPlatform.Useruuid, nil)}}
		// log.Printf("sendChatblock : %+v\n", sendChatblock)
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return nil
	}
	blockType := "room"
	blockTime, err := strconv.ParseInt(shellCmd[3], 10, 64)
	blockTime = blockTime*60*1000 + time.Now().UnixNano()/int64(time.Millisecond)

	if err != nil {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_TIME_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)
		return err
	}

	var blockUuid = common.Getid().Hexstring()
	_, err = database.Exec(
		"DELETE FROM `chatBlock` WHERE blockUserUuid = ? and blocktarget = ? ",
		blockUserPlatform.Useruuid,
		roomInfo.Roomcore.Roomuuid,
	)
	// log.Printf("DELETE FROM `chatBlock` err: %+v\n", err)
	if err != nil {

		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_DB_DELETE_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return err
	}

	_, err = database.Exec(
		"INSERT INTO chatBlock (blockUuid, blockUserUuid, blocktarget , blockType, platformUuid, platform, timeStamp) VALUES (?, ? , ? , ? , ? , ? , ? )",
		blockUuid,
		blockUserPlatform.Useruuid,
		roomInfo.Roomcore.Roomuuid,
		blockType,
		blockUserPlatform.Platformuuid,
		blockUserPlatform.Platform,
		blockTime,
	)

	// log.Printf("INSERT INTO chatBlock err: %+v\n", err)
	if err != nil {

		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_BLOCKUSER_DB_INSERT_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connCore, SendMsgJson)

		return err
	}

	SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
	SendMsgJson, _ := json.Marshal(SendMsg)
	common.Sendmessage(connCore, SendMsgJson)

	//更新列表
	pubData := common.Syncdata{Synctype: "blockSync", Data: ""}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	common.Pubsudoresult(roomInfo.Roomcore.Roomuuid, userPlatform, shellCmd, blockUserPlatform, "Quiet")

	return nil
}
