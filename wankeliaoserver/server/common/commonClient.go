package common

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"../database"
	"../socket"
	"github.com/gorilla/websocket"
)

func Clientsread(loginUuid string) (Client, bool) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client, ok := Clients[loginUuid]
	return client, ok
}

func Clientsinsert(loginUuid string, client Client) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	Clients[loginUuid] = client
}

func Clientsuserplatformread(loginUuid string) (socket.Userplatform, bool) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client, ok := Clients[loginUuid]
	userPlatform := client.Userplatform
	return userPlatform, ok
}

func Clientsroomread(loginUuid string, roomUuid string) (socket.Roomcore, bool) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	roomInfo, ok := Clients[loginUuid].Room[roomUuid]
	return roomInfo, ok
}

func Clientsroominsert(loginUuid string, roomUuid string, roomCore socket.Roomcore) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	Clients[loginUuid].Room[roomUuid] = roomCore
}

func Clientsdelete(loginUuid string) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	delete(Clients, loginUuid)
}

func Clientsroomdelete(loginUuid string, roomUuid string) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	delete(Clients[loginUuid].Room, roomUuid)
}

func Clientssidetextmapread(loginUuid string) (map[string]Sidetextplatform, bool) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client, ok := Clients[loginUuid]
	return client.Sidetext, ok
}

func Clientssidetextinsert(loginUuid string, sideText map[string]Sidetextplatform) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	// log.Printf("Clientssidetextinsert loginUuid : %+v\n", loginUuid)
	// log.Printf("Clientssidetextinsert sideText : %+v\n", sideText)
	clients := Clients[loginUuid]
	clients.Sidetext = sideText
	Clients[loginUuid] = clients
}

func Clientssidetextuserread(loginUuid string, chatTargetUuid string) (Sidetextplatform, bool) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client, ok := Clients[loginUuid]
	if !ok {
		return Sidetextplatform{}, false
	}
	sideTextUser, ok := client.Sidetext[chatTargetUuid]
	return sideTextUser, ok
}

func Clientssidetextuserinsert(loginUuid string, chatTargetUuid string, sideTextUser Sidetextplatform) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client, ok := Clients[loginUuid]
	// sideTextUserJson, _ := json.Marshal(sideTextUser)
	// Essyslog("Clientssidetextuserinsert chatTargetUuid : "+chatTargetUuid, loginUuid, client.Userplatform.Useruuid)
	// Essyslog("Clientssidetextuserinsert sideTextUser : "+string(sideTextUserJson), loginUuid, client.Userplatform.Useruuid)
	if ok {
		client.Sidetext[chatTargetUuid] = sideTextUser
	}
}

func Clientssidetextuserdelete(loginUuid string, targetUuid string) {
	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	delete(Clients[loginUuid].Sidetext, targetUuid)
}

func Clientsconnectread(uuid string) (map[string]*websocket.Conn, bool) {
	Mutexclientsconnect.Lock()
	defer Mutexclientsconnect.Unlock()
	clientConnect, ok := Clientsconnect[uuid]
	return clientConnect, ok
}

func Clientsconnectinsert(uuid string, connect map[string]*websocket.Conn) {
	Mutexclientsconnect.Lock()
	defer Mutexclientsconnect.Unlock()
	Clientsconnect[uuid] = connect
}

func Clientsconnectdelete(uuid string) {
	Mutexclientsconnect.Lock()
	defer Mutexclientsconnect.Unlock()
	delete(Clientsconnect, uuid)
}

func Clientsconnectloginuuidinsert(uuid string, loginuuid string, connect *websocket.Conn) {
	Mutexclientsconnect.Lock()
	defer Mutexclientsconnect.Unlock()
	Clientsconnect[uuid][loginuuid] = connect
}

func Clientsconnectloginuuiddelete(uuid string, loginuuid string) {
	Mutexclientsconnect.Lock()
	defer Mutexclientsconnect.Unlock()
	delete(Clientsconnect[uuid], loginuuid)
}

func Usersinforead(userUuid string) (socket.User, bool) {
	Mutexusersinfo.Lock()
	defer Mutexusersinfo.Unlock()
	Usersinfo, ok := Usersinfo[userUuid]
	return Usersinfo, ok
}

func Usersinfoinsert(userUuid string, user socket.User) {
	Mutexusersinfo.Lock()
	defer Mutexusersinfo.Unlock()
	Usersinfo[userUuid] = user
}

func Usersinfodelete(userUuid string) {
	Mutexusersinfo.Lock()
	defer Mutexusersinfo.Unlock()
	delete(Usersinfo, userUuid)
}

func Userfriendlistread(userUuid string) (map[string]socket.Friendplatform, bool) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	friendList, ok := UserfriendList[userUuid]
	return friendList, ok
}

func Userfriendlistinsert(userUuid string, friendList map[string]socket.Friendplatform) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	UserfriendList[userUuid] = friendList
}

func Userfriendlistdelete(userUuid string) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	delete(UserfriendList, userUuid)
}

func Userfriendlistuserread(userUuid string, targetUuid string) (socket.Friendplatform, bool) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	friendList, ok := UserfriendList[userUuid]
	if !ok {
		return socket.Friendplatform{}, false
	}
	friend, ok := friendList[targetUuid]
	return friend, ok
}

func Userfriendlistuserinsert(userUuid string, targetUuid string, friend socket.Friendplatform) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	friendList, ok := UserfriendList[userUuid]
	if !ok {
		return
	}
	friendList[targetUuid] = friend
	UserfriendList[userUuid] = friendList
}

func Userfriendlistuserdelete(userUuid string, targetUuid string) {
	Mutexuserfriendlist.Lock()
	defer Mutexuserfriendlist.Unlock()
	friendList, ok := UserfriendList[userUuid]
	if !ok {
		return
	}
	delete(friendList, targetUuid)
	UserfriendList[userUuid] = friendList
}

func Queryuserinfo(userUuid string) {

	client, ok := Clientsconnectread(userUuid)
	if !ok {
		return
	}

	var user socket.User
	row := database.QueryRow("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE uuid = ?",
		userUuid,
	)
	err := row.Scan(&user.Userplatform.Useruuid, &user.Userplatform.Platformuuid, &user.Userplatform.Platform, &user.Globalrole)
	if err != nil {
		return
	}

	// userInfoJson, _ := json.Marshal(user)
	// Essyslog("Queryuserinfo user : "+string(userInfoJson), "", userUuid)

	rows, _ := database.Query("select roomUuid from vipGroupUserList where userUuid = ?",
		userUuid,
	)
	for rows.Next() {
		var roomUuid string
		rows.Scan(&roomUuid)
		if user.Vipgroup == "" {
			user.Vipgroup = roomUuid
		} else {
			user.Vipgroup = user.Vipgroup + "," + roomUuid
		}
	}
	rows.Close()

	rows, _ = database.Query("select roomUuid from privateGroupUserList where userUuid = ?",
		userUuid,
	)
	for rows.Next() {
		var roomUuid string
		rows.Scan(&roomUuid)
		if user.Privategroup == "" {
			user.Privategroup = roomUuid
		} else {
			user.Privategroup = user.Privategroup + "," + roomUuid
		}
	}
	rows.Close()

	Usersinfoinsert(userUuid, user)
	// log.Printf("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE uuid = " + userUuid + "\n")
	// log.Printf("Queryuserinfo user : %+v\n", user)
	Setredisuserinfo(userUuid, user)

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendUserInfo := socket.Cmd_b_user_info_update{Base_B: socket.Base_B{Cmd: socket.CMD_B_USER_INFO_UPDATE, Stamp: timeUnix}, Payload: user}
	sendUserInfoJson, _ := json.Marshal(sendUserInfo)

	for _, connect := range client {
		Sendmessage(connect, sendUserInfoJson)
	}
}

func Userinfosyncandemit(dataJson string) {

	var data Redispubsubuserdata

	err := json.Unmarshal([]byte(dataJson), &data)
	if err != nil {
		panic(err)
	}

	// log.Printf("Userinfosyncandemit : %+v\n", data.Datajson)
	// log.Printf("Userinfosyncandemit time: %+v\n", time.Now())

	user, ok := Clientsconnectread(data.Useruuid)

	if ok {
		// Essyslog("Userinfosyncandemit user : "+dataJson, "", data.Useruuid)
		Queryuserinfo(data.Useruuid)
		for _, connect := range user {
			Sendmessage(connect, []byte(data.Datajson))

		}
	}
}

func Hierarchytargetinfosearch(loginUuid string, userUuid string, targetUuid string) (socket.User, bool, socket.Exception) {
	//因為需要Globalrole所以取User
	exception := socket.Exception{}
	targetInfo := socket.User{}
	if targetUuid == "" {
		exception = Exception("COMMON_HIERARCHYTARGETINFOSEARCH_TARGET_UUID_NULL", userUuid, nil)
		return targetInfo, false, exception
	}
	targetInfo, ok := Usersinforead(targetUuid)
	if !ok {
		targetInfo, ok = Getredisuserinfo(targetUuid)
		if ok {
			log.Printf("Getredisuserinfo targetUuid = "+targetUuid+"  targetInfoUser : %+v\n", targetInfo)
			return targetInfo, true, exception
		}

		row := database.QueryRow("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE uuid = ?",
			targetUuid,
		)
		err := row.Scan(&targetInfo.Userplatform.Useruuid, &targetInfo.Userplatform.Platformuuid, &targetInfo.Userplatform.Platform, &targetInfo.Globalrole)
		if err != nil {
			exception = Exception("COMMON_HIERARCHYTARGETINFOSEARCH_SELECT_USER_ERROR", userUuid, nil)
			return targetInfo, false, exception
		}

		// log.Printf("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE uuid = " + targetUuid + "\n")
		// log.Printf("Hierarchytargetinfosearch SELECT targetInfoUser : %+v\n", targetInfo)

		rows, _ := database.Query("select roomUuid from vipGroupUserList where userUuid = ?",
			targetUuid,
		)
		for rows.Next() {
			var roomUuid string
			rows.Scan(&roomUuid)
			if targetInfo.Vipgroup == "" {
				targetInfo.Vipgroup = roomUuid
			} else {
				targetInfo.Vipgroup = targetInfo.Vipgroup + "," + roomUuid
			}
		}
		rows.Close()

		// log.Printf("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE uuid = " + targetUuid + "\n")
		// log.Printf("Hierarchytargetinfosearch vipGroup targetInfoUser : %+v\n", targetInfo)

		rows, _ = database.Query("select roomUuid from privateGroupUserList where userUuid = ?",
			targetUuid,
		)
		for rows.Next() {
			var roomUuid string
			rows.Scan(&roomUuid)
			if targetInfo.Privategroup == "" {
				targetInfo.Privategroup = roomUuid
			} else {
				targetInfo.Privategroup = targetInfo.Privategroup + "," + roomUuid
			}
		}
		rows.Close()

		//log.Printf("targetInfo Usersinforead ok : %+v \n", ok)
		//log.Printf("Hierarchytargetinfosearch targetInfoUser : %+v\n", targetInfo)
		Setredisuserinfo(targetUuid, targetInfo)

	}

	return targetInfo, true, exception
}
