package commonFunc

import (
	_ "net/http/pprof"

	"../commonData"

	"github.com/gorilla/websocket"
)

func SendMessage(conn commonData.ConnCore, msg []byte) {

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

	EsSysLog(string(msg), conn.LoginUuid, conn.LoginUuid)

	// log.Printf("exit Sendmessage\n")
	// log.Printf("conn : %+v\n", conn)

}

func copyConn(roomUuid string) []commonData.ConnCore {

	// log.Printf("roomUuid : %s\n", roomUuid)
	commonData.MutexRooms.Lock()
	// log.Printf("Mutexrooms :Broadcast\n")
	defer func() {
		commonData.MutexRooms.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexrooms :BroadcastUNLock\n")
	}()

	targetroom := commonData.Rooms[roomUuid]

	var connArray []commonData.ConnCore

	for loginUuid := range targetroom {
		// log.Printf("loginUuid : %+v\n", loginUuid)
		connArray = append(connArray, targetroom[loginUuid].ConnCore)

	}
	return connArray
}

func Broadcast(roomUuid string, msg []byte) {

	connArray := copyConn(roomUuid)

	for _, conn := range connArray {
		// log.Printf("loginUuid : %+v\n", loginUuid)
		SendMessage(conn, msg)
	}
	// log.Printf("-------------------------end-------------------------\n")
	return
}
