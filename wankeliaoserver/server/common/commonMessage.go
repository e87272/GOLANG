package common

import (
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"../database"
	"../socket"
)

func Sendmessage(conn *websocket.Conn, msg []byte) {

	//加鎖 加鎖 加鎖
	Mutexconnect.Lock()
	// log.Printf("Mutexconnect :SendmessageLock\n")
	defer func() {
		Mutexconnect.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexconnect :SendmessageUNLock\n")
	}()
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		// Handle error
		// log.Printf("WriteMessage : %+v\n", err)
	}

	// log.Printf("exit Sendmessage\n")
	// log.Printf("conn : %+v\n", conn)

}

func Broadcast(roomUuid string, msg []byte) {

	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Broadcast\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :BroadcastUNLock\n")
	}()

	var data socket.Cmd_b_player_speak

	err := json.Unmarshal(msg, &data)
	if err != nil {
		panic(err)
	}

	_, ok := Usersinforead(data.Payload.Chatmessage.From.Useruuid)

	if ok {
		Isspeakcd(data.Payload.Chatmessage.From.Useruuid, data.Payload.Chatmessage.Stamp)
	}

	// log.Printf("roomUuid : %s\n", roomUuid)
	targetroom := Rooms[roomUuid]
	// log.Printf("targetroom : %+v\n", len(targetroom))

	for loginUuid := range targetroom {
		Sendmessage(targetroom[loginUuid].Conn, msg)
	}
}

func Isspeakcd(userUuid string, timeUnix string) bool {
	Mutexspeakcdtime.Lock()
	defer Mutexspeakcdtime.Unlock()

	now, _ := strconv.ParseInt(timeUnix, 10, 64)

	cdTime, ok := Speakcdlist[userUuid]
	if !ok {
		Speakcdlist[userUuid] = now + Speakcdtime
		return false
	}
	if now < cdTime {
		return true
	}
	Speakcdlist[userUuid] = now + Speakcdtime
	return false
}

func BroadcastAdmin(roomUuid string, msg []byte, functionName string) {

	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :BroadcastAdmin\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :BroadcastAdminUNLock\n")
	}()

	// log.Printf("roomUuid : %s\n", roomUuid)
	targetroom := Rooms[roomUuid]
	// log.Printf("targetroom : %+v\n", targetroom)

	for uuid := range targetroom {
		// log.Printf("targetroom uuid: %+v\n", uuid)
		if Checkadmin(roomUuid, targetroom[uuid].Userplatform.Useruuid, functionName) {
			// log.Printf("BroadcastAdmin msg : %s\n", string(msg))
			Sendmessage(targetroom[uuid].Conn, msg)
		}
	}
}

var mutexGlobalMessage sync.Mutex
var globalMessageList = map[string]GlobalMessage{}
var globalMessageTimer = time.NewTimer(0)

func Queryglobalmessage() {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	rows, err := database.Query("SELECT uuid,station,content,startTime,endTime,timeInterval FROM globalMessage WHERE endTime > ?", now)
	if err != nil {
		Essyserrorlog("COMMON_QUERYGLOBALMESSAGE_ERROR", "", err)
		return
	}

	mutexGlobalMessage.Lock()
	defer mutexGlobalMessage.Unlock()

	globalMessageTimer.Stop()
	globalMessageList = map[string]GlobalMessage{}
	nextTickTime := int64(math.MaxInt64)

	for rows.Next() {
		uuid := ""
		msg := GlobalMessage{}
		rows.Scan(&uuid, &msg.Station, &msg.Content, &msg.Ticktime, &msg.Endtime, &msg.Timeinterval)

		if now >= msg.Ticktime && msg.Timeinterval > 0 {
			msg.Ticktime += ((now-msg.Ticktime)/msg.Timeinterval + 1) * msg.Timeinterval
		}

		if msg.Ticktime > now && msg.Ticktime <= msg.Endtime {
			globalMessageList[uuid] = msg
			if msg.Ticktime < nextTickTime {
				nextTickTime = msg.Ticktime
			}
		}
	}

	if len(globalMessageList) > 0 {
		globalMessageTimer.Reset(time.Duration((nextTickTime - now) * int64(time.Millisecond)))
	}
}

func Servertick() {
	for {
		select {
		case <-globalMessageTimer.C:
			globalMessageTimerTick()
		}
	}
}

func globalMessageTimerTick() {

	mutexGlobalMessage.Lock()
	defer mutexGlobalMessage.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)
	nextTickTime := int64(math.MaxInt64)

	for uuid, msg := range globalMessageList {

		if now >= msg.Ticktime {
			go globalMessageBroadcast(msg.Station, msg.Content)
			if msg.Timeinterval > 0 {
				msg.Ticktime += ((now-msg.Ticktime)/msg.Timeinterval + 1) * msg.Timeinterval
			}
		}

		if msg.Ticktime > now && msg.Ticktime <= msg.Endtime {
			if msg.Ticktime < nextTickTime {
				nextTickTime = msg.Ticktime
			}
		} else {
			delete(globalMessageList, uuid)
		}
	}

	if len(globalMessageList) > 0 {
		globalMessageTimer.Reset(time.Duration((nextTickTime - now) * int64(time.Millisecond)))
	}
}

func globalMessageBroadcast(station string, msg string) {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendBroadcast := socket.Cmd_b_global_message{
		Base_B: socket.Base_B{
			Cmd:   socket.CMD_B_GLOBAL_MESSAGE,
			Stamp: timeUnix,
		},
		Payload: socket.Globalmessage{
			Historyuuid: Getid().Hexstring(),
			Station:     station,
			Message:     msg,
		},
	}
	sendBroadcastJson, _ := json.Marshal(sendBroadcast)

	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	for _, client := range Clients {
		Sendmessage(client.Conn, sendBroadcastJson)
	}
}
