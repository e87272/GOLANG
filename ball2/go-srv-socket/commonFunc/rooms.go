package commonFunc

import (
	_ "net/http/pprof"

	"../commonData"
)

func RoomsInsert(roomUuid string, client map[string]commonData.RoomClient) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	commonData.Rooms[roomUuid] = client

}

func RoomsRead(roomUuid string) (map[string]commonData.RoomClient, bool) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	roomClient, ok := commonData.Rooms[roomUuid]

	if !ok {
		return make(map[string]commonData.RoomClient), false
	}
	return roomClient, true

}

func RoomsDelete(roomUuid string) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	delete(commonData.Rooms, roomUuid)

}
func RoomsClientInsert(roomUuid string, connCore commonData.ConnCore, userUuid string) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	var roomClient = commonData.RoomClient{}
	roomClient.ConnCore = connCore
	roomClient.UserUuid = userUuid
	_, ok := commonData.Rooms[roomUuid]
	if !ok {
		var roomMap = make(map[string]commonData.RoomClient)
		roomMap[connCore.LoginUuid] = roomClient
		commonData.Rooms[roomUuid] = roomMap
	}
	commonData.Rooms[roomUuid][connCore.LoginUuid] = roomClient

	EsSysLog(commonData.Rooms[roomUuid], connCore.LoginUuid, userUuid)

	return

}
func RoomsClientRead(roomUuid string, loginUuid string) (commonData.RoomClient, bool) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	roomClient, ok := commonData.Rooms[roomUuid][loginUuid]
	if !ok {
		return commonData.RoomClient{}, false
	}

	return roomClient, true

}
func RoomsClientDelete(roomUuid string, loginUuid string) {
	commonData.MutexRooms.Lock()
	defer commonData.MutexRooms.Unlock()

	delete(commonData.Rooms[roomUuid], loginUuid)

	return

}

func CheckInRoom(roomUuid string, loginUuid string) bool {

	roomClient, ok := commonData.Rooms[roomUuid]

	EsSysLog(commonData.Rooms[loginUuid], loginUuid, loginUuid)

	if !ok {
		return false
	}

	_, ok = roomClient[loginUuid]

	if !ok {
		return false
	}

	return true
}
