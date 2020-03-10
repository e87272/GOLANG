package commonData

import "sync"

var MutexIpList = new(sync.Mutex)

var IpList = make(map[string]string)

var MutexRooms = new(sync.Mutex)

var Rooms = make(map[string]map[string]RoomClient) //key1 : roomUuid  key2 : loginUuid

var MutexUsersInfo = new(sync.Mutex)

var UsersInfo = make(map[string]UserInfo) //key : userUuid

var MutexClients = new(sync.Mutex)

var Clients = make(map[string]Client) //key : loginUuid

var MutexClientsConnections = new(sync.Mutex)

var ClientsConnections = make(map[string]map[string]ConnCore) //key1 : userUuid  key2 : loginUuid
