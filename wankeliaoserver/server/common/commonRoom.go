package common

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"../database"
	"../socket"
)

func Checkinroom(roomUuid string, uuid string) bool {

	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	client := Clients[uuid]

	// log.Printf("client : %+v\n", client)

	for _, clientRoom := range client.Room {
		if roomUuid == clientRoom.Roomuuid {
			return true
		}
	}
	return false
}

func Roomsread(roomUuid string) map[string]Roomclient {
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsread\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsreadUnlock\n")
	}()
	rooms := Rooms[roomUuid]
	return rooms
}

func Roomsreadmemberlist(roomUuid string) []socket.Userplatform {

	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsread\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsreadUnlock\n")
	}()

	memberList := []socket.Userplatform{}
	for _, v := range Rooms[roomUuid] {
		memberList = append(memberList, v.Userplatform)
	}
	return memberList
}

func Roomsmembercount(roomUuid string) int {

	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsread\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsreadUnlock\n")
	}()

	return len(Rooms[roomUuid])
}

func Roomsinsert(roomUuid string, roomClient map[string]Roomclient) {
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsinsert\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsinsertUnlock\n")
	}()
	Rooms[roomUuid] = roomClient
}

func Roomsclientread(roomUuid string, loginUuid string) Roomclient {
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsclientread\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsclientreadUnlock\n")
	}()
	roomsClient := Rooms[roomUuid][loginUuid]
	return roomsClient
}

func Roomsclientinsert(roomUuid string, loginUuid string, roomClient Roomclient) {
	//加鎖 加鎖 加鎖
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsclientinsert\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsclientinsertUnlock\n")
	}()
	Rooms[roomUuid][loginUuid] = roomClient
}

func Roomsdelete(roomUuid string) {
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsdelete\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsdeleteUnlock\n")
	}()
	delete(Rooms, roomUuid)
}

func Roomsclientdelete(roomUuid string, uuid string) {
	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Roomsclientdelete\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :RoomsclientdeleteUnlock\n")
	}()
	delete(Rooms[roomUuid], uuid)
}

func Roomsinforead(roomUuid string) (socket.Roominfo, bool) {
	Mutexroomsinfo.Lock()
	defer Mutexroomsinfo.Unlock()
	roomsinfo, ok := Roomsinfo[roomUuid]
	return roomsinfo, ok
}

func Roomsinfolist() []socket.Roominfo {
	Mutexroomsinfo.Lock()
	defer Mutexroomsinfo.Unlock()

	roomList := make([]socket.Roominfo, 0, len(Roomsinfo))
	for _, roominfo := range Roomsinfo {
		roomList = append(roomList, roominfo)
	}

	return roomList
}

func Roomsinfoinsert(roomUuid string, roomInfo socket.Roominfo) {
	Mutexroomsinfo.Lock()
	defer Mutexroomsinfo.Unlock()
	Roomsinfo[roomUuid] = roomInfo
}

func Roomsinfodelete(roomUuid string) {
	Mutexroomsinfo.Lock()
	defer Mutexroomsinfo.Unlock()
	delete(Roomsinfo, roomUuid)
}

func Queryroominfo(clientName string, roomType string, roomUuid string) string {

	var roomName string
	var roomIcon string
	var owner string
	if !Checkword(roomType) {
		//block處理
		Essyserrorlog("COMMON_QUERYROOMINFO_ERROR", clientName, nil)
		return "COMMON_QUERYROOMINFO_ERROR"
	}

	row := database.QueryRow("select roomName,roomIcon,owner from "+roomType+" where roomUuid = ?",
		roomUuid,
	)
	err := row.Scan(&roomName, &roomIcon, &owner)
	if err != nil {
		Essyserrorlog("COMMON_QUERYROOMINFO_ROOM_READ_ERROR", clientName, err)
		return "COMMON_QUERYROOMINFO_ROOM_READ_ERROR"
	}

	ownerUser, ok, exception := Hierarchytargetinfosearch("", "clientName", owner)
	if !ok {
		return exception.Message
	}

	adminSet := map[string]string{}
	userListName := roomType + "UserList"
	rows, _ := database.Query("select userUuid,roleSet from "+userListName+" where roomUuid = ? and roleSet != ?",
		roomUuid,
		"",
	)
	for rows.Next() {
		var userUuid string
		var roleSet string
		rows.Scan(&userUuid, &roleSet)
		adminSet[userUuid] = roleSet
	}
	rows.Close()
	adminSetJson, _ := json.Marshal(adminSet)

	roomInfo := socket.Roominfo{Roomuuid: roomUuid, Roomname: roomName, Roomtype: roomType, Roomicon: roomIcon, Adminset: string(adminSetJson), Ownerplatform: ownerUser.Userplatform}

	Setredisroominfo(roomUuid, roomInfo)

	return "ok"
}

func Syncroominfo(dataJson string) {
	var data map[string]string
	err := json.Unmarshal([]byte(dataJson), &data)
	if err != nil {
		Essyserrorlog("COMMON_QUERYROOMSINFO_ERROR", "", err)
		return
	}

	roomInfo, ok := Roomsinforead(data["roomUuid"])
	if ok {

		roomInfo, ok = Getredisroominfo(data["roomUuid"])
		if ok {
			Roomsinfoinsert(data["roomUuid"], roomInfo)
			// log.Printf("Syncroominfo roomInfo : %+v\n", roomInfo)
			roomInfo.Roomicon = os.Getenv("linkPath") + roomInfo.Roomicon
			timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			sendRoomInfoBroadcast := socket.Cmd_b_room_info_update{Base_B: socket.Base_B{Cmd: socket.CMD_B_ROOM_INFO_UPDATE, Stamp: timeUnix}, Payload: roomInfo}
			sendRoomInfoBroadcastJson, _ := json.Marshal(sendRoomInfoBroadcast)
			Broadcast(data["roomUuid"], sendRoomInfoBroadcastJson)
		}
	}
}

func Hierarchyroominfosearch(loginUuid string, client Client, roomType string, roomUuid string) (socket.Roominfo, bool, socket.Exception) {

	roomInfo := socket.Roominfo{}
	exception := socket.Exception{}

	if roomUuid == "" {
		exception = Exception("COMMON_HIERARCHYROOMINFOSEARCH_ROOM_UUID_NULL", client.Userplatform.Useruuid, nil)
		return roomInfo, false, exception
	}

	roomInfo, ok := Roomsinforead(roomUuid)

	if !ok {
		roomInfo, ok = Getredisroominfo(roomUuid)
		if ok {
			var adduser = make(map[string]Roomclient)
			var roomClient = Roomclient{}
			roomClient.Conn = client.Conn
			roomClient.Userplatform = client.Userplatform
			adduser[loginUuid] = roomClient
			Roomsinsert(roomUuid, adduser)
			Roomsinfoinsert(roomUuid, roomInfo)
			// log.Printf("Hierarchyroominfosearch Getredisroominfo roomInfo : %+v\n", roomInfo)
			return roomInfo, true, exception
		}

		var roomName string
		var roomIcon string
		var ownerUuid string
		if !Checkword(roomType) {
			//block處理
			exception = Exception("COMMON_HIERARCHYROOMINFOSEARCH_ROOM_TYPE_NOT_WORD", client.Userplatform.Useruuid, nil)
			return roomInfo, false, exception
		}

		row := database.QueryRow("select roomName,roomIcon,owner from "+roomType+" where roomUuid = ?",
			roomUuid,
		)
		err := row.Scan(&roomName, &roomIcon, &ownerUuid)
		if err != nil {
			exception = Exception("COMMON_HIERARCHYROOMINFOSEARCH_ROOM_READ_ERROR", client.Userplatform.Useruuid, err)
			return roomInfo, false, exception
		}

		ownerUser, _, _ := Hierarchytargetinfosearch(loginUuid, client.Userplatform.Useruuid, ownerUuid)
		log.Printf("Hierarchyroominfosearch Getredisroominfo ownerUser : %+v\n", ownerUser)

		adminSet := map[string]string{}
		userListName := roomType + "UserList"
		rows, _ := database.Query("select userUuid,roleSet from "+userListName+" where roomUuid = ? and roleSet != ?",
			roomUuid,
			"",
		)
		for rows.Next() {
			var userUuid string
			var roleSet string
			rows.Scan(&userUuid, &roleSet)
			adminSet[userUuid] = roleSet
		}
		rows.Close()
		adminSetJson, _ := json.Marshal(adminSet)

		roomInfo = socket.Roominfo{Roomuuid: roomUuid, Roomname: roomName, Roomtype: roomType, Roomicon: roomIcon, Adminset: string(adminSetJson), Ownerplatform: ownerUser.Userplatform}
		var adduser = make(map[string]Roomclient)
		var roomClient = Roomclient{}
		roomClient.Conn = client.Conn
		roomClient.Userplatform = client.Userplatform
		adduser[loginUuid] = roomClient
		Roomsinsert(roomUuid, adduser)
		Roomsinfoinsert(roomUuid, roomInfo)

		roomInfoJson, _ := json.Marshal(roomInfo)
		Essyslog("Hierarchyroominfosearch roomInfo : "+string(roomInfoJson), loginUuid, client.Userplatform.Useruuid)

		Setredisroominfo(roomUuid, roomInfo)
	} else {
		//房間已存在不能覆蓋掉原本的client
		roomClient := Roomclient{}
		roomClient.Conn = client.Conn
		roomClient.Userplatform = client.Userplatform
		Roomsclientinsert(roomUuid, loginUuid, roomClient)
	}
	roomInfo.Roomicon = os.Getenv("linkPath") + roomInfo.Roomicon
	return roomInfo, true, exception
}

func Roomspopulationcount(data Redispubsubroomsinfo) {

	Mutexroomspopulation.Lock()
	defer Mutexroomspopulation.Unlock()
	if Roomspopulation[data.Ip] == nil {
		var rooms = make(map[string]int)
		rooms[data.RoomUuid] = data.Usercount
		Roomspopulation[data.Ip] = rooms
	} else {
		if data.Usercount == 0 {
			delete(Roomspopulation, data.RoomUuid)
		} else {
			Roomspopulation[data.Ip][data.RoomUuid] = data.Usercount
		}
	}
	numbers := 0
	for _, rooms := range Roomspopulation {
		values, ok := rooms[data.RoomUuid]
		if ok {
			numbers += values
		}
	}

}

func Roomspopulationcountread(roomUuid string) int {

	var numbers int = 0
	Mutexroomspopulation.Lock()
	defer Mutexroomspopulation.Unlock()
	for _, rooms := range Roomspopulation {
		values, ok := rooms[roomUuid]
		if ok {
			numbers = numbers + values
		}
	}
	return numbers
}

func Hierarchymembercount(loginUuid string, client Client, roomType string, roomUuid string) (int, bool, socket.Exception) {

	memberCount := 0
	exception := socket.Exception{}

	memberCount, ok := Getredismembercount(roomUuid)

	if !ok {

		userListName := roomType + "UserList"

		row := database.QueryRow("SELECT count(*) FROM users RIGHT JOIN "+userListName+" ON users.uuid="+userListName+".userUuid WHERE "+userListName+".roomUuid = ?",
			roomUuid,
		)
		err := row.Scan(&memberCount)

		// log.Printf("SELECT users.uuid, users.platform, users.platformUuid FROM users RIGHT JOIN " + userListName + " ON users.uuid=" + userListName + ".userUuid WHERE " + userListName + ".roomUuid = " + packetGetMemberList.Payload.Roomuuid + "\n")

		if err != nil {
			exception = Exception("COMMON_HIERARCHYMEMBERCOUNT_QUERY_DB_ERROR", client.Userplatform.Useruuid, err)
			return memberCount, false, exception
		}

		Setredismembercount(roomUuid, memberCount)
	}
	return memberCount, true, exception
}
