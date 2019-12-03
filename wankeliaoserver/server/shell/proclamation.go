package shell

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"../common"
	"../database"
	"../socket"
	"github.com/gorilla/websocket"
)

func linkProclamation(connect *websocket.Conn, userPlatform socket.Userplatform, packetSendMsg socket.Cmd_c_player_send_shell, timeUnix string) error {

	var uuid = userPlatform.Useruuid
	maxLength := 50

	if !common.Checkadmin(packetSendMsg.Payload.Chattarget, uuid, "RoomSub") {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_ROLE_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	argument := regexp.MustCompile(" +-").Split(strings.Trim(packetSendMsg.Payload.Message, " "), -1)
	shellCmd := regexp.MustCompile(" +").Split(argument[0], -1)
	if len(shellCmd) != 3 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	proclamation := socket.Proclamation{}
	proclamation.Proclamationuuid = common.Getid().Hexstring()
	proclamation.Roomuuid = packetSendMsg.Payload.Chattarget
	proclamation.Type = "l"
	proclamation.Order = shellCmd[2]

	order, err := strconv.ParseInt(proclamation.Order, 10, 64)
	if err != nil || order < 1 || order > 3 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_ORDER_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return err
	}

	for i := 1; i < len(argument); i++ {
		value := regexp.MustCompile(" +").Split(argument[i], 2)
		if len(value) != 2 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
		switch value[0] {
		case "a":
			proclamation.Apptype = value[1]
			break
		case "t":
			proclamation.Title = value[1]
			if utf8.RuneCountInString(proclamation.Title) > 10 {
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_TITLE_TOOLONG", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connect, SendMsgJson)
				return nil
			}
			break
		case "c":
			proclamation.Content = value[1]
			if utf8.RuneCountInString(proclamation.Content) > maxLength {
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_CONTENT_TOOLONG", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connect, SendMsgJson)
				return nil
			}
			_, proclamation.Content = common.Matchdirtyword(proclamation.Content, maxLength)
			break
		case "s":
			proclamation.Style = value[1]
			break
		case "u":
			proclamation.Url = value[1]
			break
		default:
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
	}
	if proclamation.Content != "" {
		_, err = database.Exec(
			"INSERT INTO proclamation (proclamationUuid, roomUuid, type, clientOrder, appType, title, content, style, url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			proclamation.Proclamationuuid,
			proclamation.Roomuuid,
			proclamation.Type,
			proclamation.Order,
			proclamation.Apptype,
			proclamation.Title,
			proclamation.Content,
			proclamation.Style,
			proclamation.Url,
		)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_INSERT_DB_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
	}

	_, err = database.Exec(
		"DELETE FROM proclamation WHERE roomUuid = ? AND type = ? AND clientOrder = ? AND proclamationUuid < ?",
		proclamation.Roomuuid,
		proclamation.Type,
		proclamation.Order,
		proclamation.Proclamationuuid,
	)
	if err != nil {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_LINKPROCLAMATION_DELETE_DB_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
	SendMsgJson, _ := json.Marshal(SendMsg)
	common.Sendmessage(connect, SendMsgJson)

	pubData := common.Syncdata{Synctype: "proclamationSync", Data: proclamation.Roomuuid}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	common.Pubsudoresult(packetSendMsg.Payload.Chattarget, userPlatform, shellCmd, socket.Userplatform{}, "RoomPublish")

	return nil
}

func normalProclamation(connect *websocket.Conn, userPlatform socket.Userplatform, packetSendMsg socket.Cmd_c_player_send_shell, timeUnix string) error {

	var uuid = userPlatform.Useruuid
	maxLength := 50

	if !common.Checkadmin(packetSendMsg.Payload.Chattarget, uuid, "RoomPublish") {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_ROLE_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	argument := regexp.MustCompile(" +-").Split(strings.Trim(packetSendMsg.Payload.Message, " "), -1)
	shellCmd := regexp.MustCompile(" +").Split(argument[0], -1)
	if len(shellCmd) != 2 {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	proclamation := socket.Proclamation{}
	proclamation.Proclamationuuid = common.Getid().Hexstring()
	proclamation.Roomuuid = packetSendMsg.Payload.Chattarget
	proclamation.Type = "n"
	proclamation.Order = "1"

	for i := 1; i < len(argument); i++ {
		value := regexp.MustCompile(" +").Split(argument[i], 2)
		if len(value) != 2 {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
		switch value[0] {
		case "t":
			proclamation.Title = value[1]
			if utf8.RuneCountInString(proclamation.Title) > 10 {
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_TITLE_TOOLONG", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connect, SendMsgJson)
				return nil
			}
			break
		case "c":
			proclamation.Content = value[1]
			if utf8.RuneCountInString(proclamation.Content) > maxLength {
				SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_CONTENT_TOOLONG", userPlatform.Useruuid, nil)}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				common.Sendmessage(connect, SendMsgJson)
				return nil
			}
			_, proclamation.Content = common.Matchdirtyword(proclamation.Content, maxLength)
			break
		case "s":
			proclamation.Style = value[1]
			break
		case "u":
			proclamation.Url = value[1]
			break
		default:
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_SHELL_ERROR", userPlatform.Useruuid, nil)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
	}

	if proclamation.Content != "" {
		_, err := database.Exec(
			"INSERT INTO proclamation (proclamationUuid, roomUuid, type, clientOrder, appType, title, content, style, url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			proclamation.Proclamationuuid,
			proclamation.Roomuuid,
			proclamation.Type,
			proclamation.Order,
			proclamation.Apptype,
			proclamation.Title,
			proclamation.Content,
			proclamation.Style,
			proclamation.Url,
		)
		if err != nil {
			SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_INSERT_DB_ERROR", userPlatform.Useruuid, err)}}
			SendMsgJson, _ := json.Marshal(SendMsg)
			common.Sendmessage(connect, SendMsgJson)
			return nil
		}
	}
	_, err := database.Exec(
		"DELETE FROM proclamation WHERE roomUuid = ? AND type = ? AND clientOrder = ? AND proclamationUuid < ?",
		proclamation.Roomuuid,
		proclamation.Type,
		proclamation.Order,
		proclamation.Proclamationuuid,
	)
	if err != nil {
		SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: common.Exception("SHELL_NORMALPROCLAMATION_DELETE_DB_ERROR", userPlatform.Useruuid, err)}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		common.Sendmessage(connect, SendMsgJson)
		return nil
	}

	SendMsg := socket.Cmd_r_player_send_shell{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_SHELL, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "ok", Exp: common.Exception("", "", nil)}}
	SendMsgJson, _ := json.Marshal(SendMsg)
	common.Sendmessage(connect, SendMsgJson)

	pubData := common.Syncdata{Synctype: "proclamationSync", Data: proclamation.Roomuuid}
	pubDataJson, _ := json.Marshal(pubData)
	common.Redispubdata("sync", string(pubDataJson))

	common.Pubsudoresult(packetSendMsg.Payload.Chattarget, userPlatform, shellCmd, socket.Userplatform{}, "RoomPublish")

	return nil
}
