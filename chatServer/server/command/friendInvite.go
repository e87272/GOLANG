package command

import (
	"encoding/json"
	"strconv"
	"time"

	"../common"
	"../database"
	"../socket"
)

func Friendinvite(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendInvite := socket.Cmd_r_friend_invite{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_FRIEND_INVITE,
		Stamp: timeUnix,
	}}
	client, _ := common.Clientsread(loginUuid)
	userPlatform := client.Userplatform
	userUuid := userPlatform.Useruuid
	

	var packetInvite socket.Cmd_c_friend_invite
	err := json.Unmarshal([]byte(msg), &packetInvite)
	if err != nil {
		sendInvite.Base_R.Result = "err"
		sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_JSON_ERROR", userUuid, err)
		sendInviteJson, _ := json.Marshal(sendInvite)
		common.Sendmessage(connCore, sendInviteJson)
		return err
	}
	sendInvite.Base_R.Idem = packetInvite.Base_C.Idem

	if loginUuid == userUuid {
		sendInvite.Base_R.Result = "err"
		sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_GUEST", userUuid, nil)
		sendInviteJson, _ := json.Marshal(sendInvite)
		common.Sendmessage(connCore, sendInviteJson)
		return nil
	}

	targetUserUuid := packetInvite.Payload
	isFriend := false
	targetFriend, ok := common.Userfriendlistuserread(userUuid, targetUserUuid)
	if ok {
		switch targetFriend.State {
		case "friend":
			// 已經是好友
			sendInvite.Base_R.Result = "err"
			sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_TARGET_IS_FRIEND", userUuid, nil)
			sendInviteJson, _ := json.Marshal(sendInvite)
			common.Sendmessage(connCore, sendInviteJson)
			return nil
		case "inviteTo":
			// 已發送過邀請
			sendInvite.Base_R.Result = "err"
			sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_HAVE_ALREADY_INVITED", userUuid, nil)
			sendInviteJson, _ := json.Marshal(sendInvite)
			common.Sendmessage(connCore, sendInviteJson)
			return nil
		case "inviteFrom":
			// 同意加好友
			isFriend = true
			targetFriend.State = "friend"
		default:
			// 未知狀態
			sendInvite.Base_R.Result = "err"
			sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_FRIEND_STATE_ERROR", userUuid, nil)
			sendInviteJson, _ := json.Marshal(sendInvite)
			common.Sendmessage(connCore, sendInviteJson)
			return nil
		}
	} else {

		if !common.Checkadmin("", userUuid, "AddFriend") {
			//block處理
			sendInvite.Base_R.Result = "err"
			sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_NOT_ADMIN", userUuid, nil)
			sendInviteJson, _ := json.Marshal(sendInvite)
			common.Essyslog(string(sendInviteJson), loginUuid, userUuid)
			common.Sendmessage(connCore, sendInviteJson)
			return nil
		}

		row := database.QueryRow("SELECT platformUuid,platform FROM users WHERE uuid = ?", targetUserUuid)
		err = row.Scan(&targetFriend.Userplatform.Platformuuid, &targetFriend.Userplatform.Platform)
		if err != nil {
			sendInvite.Base_R.Result = "err"
			sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_USER_UUID_ERROR", userUuid, nil)
			sendInviteJson, _ := json.Marshal(sendInvite)
			common.Sendmessage(connCore, sendInviteJson)
			return nil
		}

		targetFriend.Userplatform.Useruuid = targetUserUuid
		targetFriend.State = "inviteTo"
	}

	friendUuid := ""
	row := database.QueryRow("SELECT friendUuid FROM friendList WHERE fromUuid = ? AND toUuid = ?",
		userUuid,
		targetUserUuid,
	)
	err = row.Scan(&friendUuid)
	if err == nil {
		common.Userfriendlistuserinsert(userUuid, targetUserUuid, targetFriend)
		sendInvite.Base_R.Result = "err"
		sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_FRIEND_LIST_ERROR", userUuid, nil)
		sendInviteJson, _ := json.Marshal(sendInvite)
		common.Sendmessage(connCore, sendInviteJson)
		return nil
	} else if err != database.ErrNoRows {
		sendInvite.Base_R.Result = "err"
		sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_SELECT_DB_ERROR", userUuid, nil)
		sendInviteJson, _ := json.Marshal(sendInvite)
		common.Sendmessage(connCore, sendInviteJson)
		return nil
	}

	friendUuid = common.Getid().Hexstring()
	_, err = database.Exec(
		"INSERT INTO friendList (friendUuid, fromUuid, toUuid) VALUES (?, ?, ?)",
		friendUuid,
		userUuid,
		targetUserUuid,
	)
	if err != nil {
		sendInvite.Base_R.Result = "err"
		sendInvite.Base_R.Exp = common.Exception("COMMAND_FRIENDINVITE_INSERT_DB_ERROR", userUuid, nil)
		sendInviteJson, _ := json.Marshal(sendInvite)
		common.Sendmessage(connCore, sendInviteJson)
		return nil
	}

	sendInvite.Base_R.Result = "ok"
	sendInviteJson, _ := json.Marshal(sendInvite)
	common.Sendmessage(connCore, sendInviteJson)

	userFriend := socket.Friendplatform{}
	userFriend.Userplatform = userPlatform
	if isFriend {
		userFriend.State = "friend"
	} else {
		userFriend.State = "inviteFrom"
	}

	sendInviteBroadcast := socket.Cmd_b_friend{}
	sendInviteBroadcast.Base_B.Stamp = timeUnix
	sendInviteBroadcast.Base_B.Cmd = socket.CMD_B_FRIEND
	sendInviteBroadcast.Payload = userFriend
	sendInviteBroadcastJson, _ := json.Marshal(sendInviteBroadcast)
	common.Redispubfriendupdatestate(userFriend, targetFriend, sendInviteBroadcastJson)

	return nil
}
