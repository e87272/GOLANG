package command

import (
	"encoding/json"
	"strconv"
	"time"

	"../common"
	"../database"
	"../socket"
	"github.com/gorilla/websocket"
)

func Getfriendlist(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendFriendList := socket.Cmd_r_get_friend_list{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_FRIEND_LIST,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetFriendList socket.Cmd_c_get_friend_list
	err := json.Unmarshal([]byte(msg), &packetFriendList)
	if err != nil {
		sendFriendList.Base_R.Result = "err"
		sendFriendList.Base_R.Exp = common.Exception("COMMAND_GETFRIENDLIST_JSON_ERROR", userUuid, err)
		sendFriendListJson, _ := json.Marshal(sendFriendList)
		common.Sendmessage(connect, sendFriendListJson)
		return err
	}
	sendFriendList.Base_R.Idem = packetFriendList.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendFriendList.Base_R.Result = "err"
		sendFriendList.Base_R.Exp = common.Exception("COMMAND_GETFRIENDLIST_GUEST", userUuid, nil)
		sendFriendListJson, _ := json.Marshal(sendFriendList)
		common.Sendmessage(connect, sendFriendListJson)
		return nil
	}

	friendMap, ok := common.Userfriendlistread(userUuid)
	if !ok {
		friendMap = map[string]socket.Friendplatform{}

		rows, _ := database.Query(
			"SELECT users.uuid,platformUuid,platform FROM users RIGHT JOIN friendList ON users.uuid = friendList.toUuid WHERE friendList.fromUuid = ?",
			userUuid,
		)
		for rows.Next() {
			var target socket.Userplatform
			rows.Scan(&target.Useruuid, &target.Platformuuid, &target.Platform)
			targetFriend := socket.Friendplatform{}
			targetFriend.Userplatform = target
			targetFriend.State = "inviteTo"
			friendMap[target.Useruuid] = targetFriend
		}
		rows.Close()

		rows, _ = database.Query(
			"SELECT users.uuid,platformUuid,platform FROM users RIGHT JOIN friendList ON users.uuid = friendList.fromUuid WHERE friendList.toUuid = ?",
			userUuid,
		)
		for rows.Next() {
			var target socket.Userplatform
			rows.Scan(&target.Useruuid, &target.Platformuuid, &target.Platform)
			targetFriend, ok := friendMap[target.Useruuid]
			if !ok {
				targetFriend = socket.Friendplatform{}
				targetFriend.Userplatform = target
				targetFriend.State = "inviteFrom"
			} else if targetFriend.State == "inviteTo" {
				targetFriend.State = "friend"
			}
			friendMap[target.Useruuid] = targetFriend
		}
		rows.Close()

		common.Userfriendlistinsert(userUuid, friendMap)
	}

	friendList := []socket.Userplatform{}
	inviteFromList := []socket.Userplatform{}
	inviteToList := []socket.Userplatform{}
	for _, friend := range friendMap {
		switch friend.State {
		case "friend":
			friendList = append(friendList, friend.Userplatform)
		case "inviteFrom":
			inviteFromList = append(inviteFromList, friend.Userplatform)
		case "inviteTo":
			inviteToList = append(inviteToList, friend.Userplatform)
		}
	}

	sendFriendList.Base_R.Result = "ok"
	sendFriendList.Payload.Friendlist = friendList
	sendFriendList.Payload.Invitefromlist = inviteFromList
	sendFriendList.Payload.Invitetolist = inviteToList
	sendFriendListJson, _ := json.Marshal(sendFriendList)
	common.Sendmessage(connect, sendFriendListJson)

	return nil
}
