package commonFunc

import (
	_ "net/http/pprof"

	"../commonData"
)

func IpListRead(loginUuid string) (string, bool) {
	commonData.MutexIpList.Lock()
	defer commonData.MutexIpList.Unlock()
	ip, ok := commonData.IpList[loginUuid]
	return ip, ok
}

func IpListInsert(loginUuid string, ip string) {
	commonData.MutexIpList.Lock()
	defer commonData.MutexIpList.Unlock()
	commonData.IpList[loginUuid] = ip
}

func IpListDelete(loginUuid string) {
	commonData.MutexIpList.Lock()
	defer commonData.MutexIpList.Unlock()
	delete(commonData.IpList, loginUuid)
}

func ClientsInsert(loginUuid string, client commonData.Client) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	commonData.Clients[loginUuid] = client

}

func ClientsRead(loginUuid string) (commonData.Client, bool) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	client, ok := commonData.Clients[loginUuid]
	if !ok {
		return commonData.Client{}, false
	}

	return client, true

}
func ClientsDelete(loginUuid string) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	delete(commonData.Clients, loginUuid)
}

func ClientsRoomInsert(loginUuid string, roomUuid string) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	commonData.Clients[loginUuid].Room[roomUuid] = roomUuid

	EsSysLog(commonData.Clients[loginUuid], loginUuid, loginUuid)

	return
}

func ClientsRoomRead(loginUuid string) (map[string]string, bool) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	client, ok := commonData.Clients[loginUuid]
	if !ok {
		return make(map[string]string), false
	}

	return client.Room, true
}
func ClientsRoomDelete(loginUuid string, roomUuid string) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()

	delete(commonData.Clients[loginUuid].Room, roomUuid)

	return
}

func ClientsUserInfoRead(loginUuid string) (commonData.UserInfo, bool) {
	commonData.MutexClients.Lock()
	defer commonData.MutexClients.Unlock()
	client, ok := commonData.Clients[loginUuid]
	if !ok {
		return commonData.UserInfo{}, false
	}
	userInfo := client.UserInfo
	return userInfo, true

}

func ClientsConnectionsInsert(userUuid string, connection map[string]commonData.ConnCore) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	commonData.ClientsConnections[userUuid] = connection

	return
}

func ClientsConnectionsRead(userUuid string) (map[string]commonData.ConnCore, bool) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	connections, ok := commonData.ClientsConnections[userUuid]
	if !ok {
		return make(map[string]commonData.ConnCore), false
	}
	return connections, true
}
func ClientsConnectionsDelete(userUuid string) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	delete(commonData.ClientsConnections, userUuid)

	return
}

func ClientsConnectionsLoginUuidInsert(userUuid string, loginUuid string, connCore commonData.ConnCore) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	commonData.ClientsConnections[userUuid][loginUuid] = connCore

	return
}

func ClientsConnectionsLoginUuidRead(userUuid string, loginUuid string) (commonData.ConnCore, bool) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	connCore, ok := commonData.ClientsConnections[userUuid][loginUuid]
	if !ok {
		return commonData.ConnCore{}, false
	}
	return connCore, true
}

func ClientsConnectionsLoginUuidDelete(userUuid string, loginUuid string) {
	commonData.MutexClientsConnections.Lock()
	defer commonData.MutexClientsConnections.Unlock()

	delete(commonData.ClientsConnections[userUuid], loginUuid)

	return
}
