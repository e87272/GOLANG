package common

import (
	"context"
	"encoding/json"
	"math/rand"

	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"../database"
	"../socket"
)

func Sendmessage(conn Conncore, msg []byte) {

	//加鎖 加鎖 加鎖
	conn.Connmutex.Lock()
	// log.Printf("Mutexconnect :SendmessageLock\n")
	defer func() {
		conn.Connmutex.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexconnect :SendmessageUNLock\n")
	}()
	err := conn.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		// Handle error
		// log.Printf("WriteMessage : %+v\n", err)
	}

	// log.Printf("exit Sendmessage\n")
	// log.Printf("conn : %+v\n", conn)

}

func copyConn(roomUuid string) []Conncore {

	// log.Printf("roomUuid : %s\n", roomUuid)
	Mutexrooms.Lock()
	// log.Printf("Mutexrooms :Broadcast\n")
	defer func() {
		Mutexrooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :BroadcastUNLock\n")
	}()

	targetroom := Rooms[roomUuid]

	var connArray []Conncore

	for loginUuid := range targetroom {
		// log.Printf("loginUuid : %+v\n", loginUuid)
		connArray = append(connArray, targetroom[loginUuid].Conncore)

	}
	return connArray
}

func Broadcast(roomUuid string, msg []byte, packetStamp int64) {

	connArray := copyConn(roomUuid)
	
	if len(connArray) > 0{
		index := rand.Intn(len(connArray))

		ok := Dropmsg(connArray[index], msg, packetStamp)
		if ok {
			timeUnix := time.Now().UnixNano() / int64(time.Millisecond)
			Essyslog(string(msg), "timeUnix : "+strconv.FormatInt(timeUnix, 10)+" packetStamp : "+strconv.FormatInt(packetStamp, 10), "Broadcast")
			return
		}
	}

	for _, conn := range connArray {
		// log.Printf("loginUuid : %+v\n", loginUuid)
		Sendmessage(conn, msg)
	}
	// log.Printf("-------------------------end-------------------------\n")
	return
}

func Dropmsg(connCore Conncore, msg []byte, packetStamp int64) bool {

	connCore.Connmutex.Lock()
	defer connCore.Connmutex.Unlock()

	//訊息過久直接略過
	timeUnix := time.Now().UnixNano() / int64(time.Millisecond)
	if timeUnix-packetStamp > Packetdroptime {
		return true
	}

	return false
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

func Isnewuser(userUuid string) bool {

	// log.Printf("Isnewuser : %+v\n", userUuid)
	cdTime, ok := Blocknewuserlistread(userUuid)
	// log.Printf("Isnewuser cdTime : %+v\n", cdTime)
	if !ok {
		return false
	}

	if time.Now().UnixNano()/int64(time.Millisecond) > cdTime {
		Blocknewuserlistdelete(userUuid)
		return false
	}

	return true
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
			Sendmessage(targetroom[uuid].Conncore, msg)
		}
	}
}

var mutexGlobalMessage = new(sync.Mutex)
var globalMessageList = map[string]GlobalMessage{}
var globalMessageTimer = time.NewTimer(0)

func Queryglobalmessage() {
	// log.Printf("Queryglobalmessage start\n")

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
	rows.Close()

	if len(globalMessageList) > 0 {
		globalMessageTimer.Reset(time.Duration((nextTickTime - now) * int64(time.Millisecond)))
	}
	// log.Printf("Queryglobalmessage end\n")
}

func Servertick() {
	for {
		select {
		case <-globalMessageTimer.C:
			// log.Printf("Servertick start\n")
			globalMessageTimerTick()
			// log.Printf("Servertick end\n")
		}
	}
}

func globalMessageTimerTick() {
	// log.Printf("globalMessageTimerTick start\n")

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
	// log.Printf("globalMessageTimerTick end\n")
}

func copyAllConn() []Conncore {

	// log.Printf("copyAllConn start\n")
	Mutexclients.Lock()
	defer Mutexclients.Unlock()

	var connArray []Conncore
	for _, client := range Clients {
		connArray = append(connArray, client.Conncore)
	}

	// log.Printf("copyAllConn end\n")
	return connArray
}

func globalMessageBroadcast(station string, msg string) {
	// log.Printf("globalMessageBroadcast start\n")

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

	connArray := copyAllConn()
	// log.Printf("connArray count : %+v\n", len(connArray))
	for _, conn := range connArray {
		Sendmessage(conn, sendBroadcastJson)
	}
	// log.Printf("globalMessageBroadcast end\n")
}

func Sendclearusermsg(data string) {
	packetStamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeUnix := strconv.FormatInt(packetStamp, 10)
	clearUserMsgData := Redispubsubclearusermsgdata{}
	json.Unmarshal([]byte(data), &clearUserMsgData)

	sendBroadcast := socket.Cmd_b_clear_user_msg{
		Base_B: socket.Base_B{
			Cmd:   socket.CMD_B_CLEAR_USER_MSG,
			Stamp: timeUnix,
		},
		Payload: socket.Clearusermsg{
			Roomuuid:   clearUserMsgData.Roomuuid,
			Targetuuid: clearUserMsgData.Targetuuid,
		},
	}
	sendBroadcastJson, _ := json.Marshal(sendBroadcast)

	Broadcast(clearUserMsgData.Roomuuid, sendBroadcastJson, packetStamp)
}

func Membercountbroadcast(data string) {

	var roomCountBroadcast socket.Cmd_b_room_member_count

	err := json.Unmarshal([]byte(data), &roomCountBroadcast)
	if err != nil {
		panic(err)
	}

	packetStamp, _ := strconv.ParseInt(roomCountBroadcast.Stamp, 10, 64)

	go Broadcast(roomCountBroadcast.Payload.Roomuuid, []byte(data), packetStamp)

}

func Hierarchyroomlastmessage(loginUuid string, userUuid string, roomCore socket.Roomcore) socket.Chatmessage {

	lastMessage := socket.Chatmessage{}

	lastMessage, ok := Getredisroomlastmessage(roomCore.Roomuuid)

	if !ok {
		boolQ := elastic.NewBoolQuery()
		boolQ.Filter(elastic.NewMatchQuery("chatTarget", roomCore.Roomuuid))

		// Search with a term query
		searchResult, err := Elasticclient.Search(os.Getenv(roomCore.Roomtype)).Query(boolQ).Sort("historyUuid", false).Size(1).Do(context.Background()) // execute

		if err != nil {
			Exception("COMMON_HIERARCHYROOMLASTMESSAGE_ES_SEARCH_ERROR", userUuid, err)
			return lastMessage
		}

		// Here's how you iterate through results with full control over each step.
		if searchResult.Hits.TotalHits.Value > 0 {

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index
				// log.Printf("hit : %+v\n", hit)
				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var chatHistory Chathistory

				err := json.Unmarshal(hit.Source, &chatHistory)
				if err != nil {
					// Deserialization failed
				}

				// Work with tweet
				// log.Printf("ChatMessage : %+v\n", chatHistory)

				lastMessage.Historyuuid = chatHistory.Historyuuid
				lastMessage.From.Useruuid = chatHistory.Myuuid
				lastMessage.From.Platformuuid = chatHistory.Myplatformuuid
				lastMessage.From.Platform = chatHistory.Myplatform
				lastMessage.Stamp = chatHistory.Stamp
				lastMessage.Message = chatHistory.Message
				lastMessage.Style = chatHistory.Style
				break
			}
		}
		Setredisroomlastmessage(roomCore.Roomuuid, lastMessage)
	}

	return lastMessage
}

func Hierarchysidetextlastmessage(loginUuid string, userUuid string, sideTextUuid string) socket.Chatmessage {

	lastMessage := socket.Chatmessage{}

	lastMessage, ok := Getredissidetextlastmessage(sideTextUuid)

	if !ok {
		boolQ := elastic.NewBoolQuery()
		boolQ.Filter(elastic.NewMatchQuery("chatTarget", sideTextUuid))
		searchResult, err := Elasticclient.Search(os.Getenv("sideText")).Query(boolQ).Sort("historyUuid", false).Size(1).Do(context.Background())

		if err != nil {
			Exception("COMMON_HIERARCHYSIDETEXTLASTMESSAGE_ES_SEARCH_ERROR", userUuid, err)
			return lastMessage
		}

		// Here's how you iterate through results with full control over each step.
		if searchResult.Hits.TotalHits.Value > 0 {

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index
				// log.Printf("hit : %+v\n", hit)
				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var chatHistory Chathistory

				err := json.Unmarshal(hit.Source, &chatHistory)
				if err != nil {
					// Deserialization failed
				}

				// Work with tweet
				// log.Printf("ChatMessage : %+v\n", chatHistory)

				lastMessage.Historyuuid = chatHistory.Historyuuid
				lastMessage.From.Useruuid = chatHistory.Myuuid
				lastMessage.From.Platformuuid = chatHistory.Myplatformuuid
				lastMessage.From.Platform = chatHistory.Myplatform
				lastMessage.Stamp = chatHistory.Stamp
				lastMessage.Message = chatHistory.Message
				lastMessage.Style = chatHistory.Style
				break
			}
		}
		Setredissidetextlastmessage(sideTextUuid, lastMessage)
	}

	return lastMessage
}
