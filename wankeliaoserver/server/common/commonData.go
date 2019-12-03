package common

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"../socket"
)

var Elasticclient *elastic.Client

var Redisclient *redis.Client

var Mutexrooms sync.Mutex

var Mutexroomsinfo sync.Mutex

var Mutexusersinfo sync.Mutex

var Mutexroomspopulation sync.Mutex

var Mutexuserfriendlist sync.Mutex

var Mutexclients sync.Mutex

var Mutexclientsconnect sync.Mutex

var Mutexconnect sync.Mutex

var Mutexelastic sync.Mutex

var Mutexredis sync.Mutex

var Mutexblockchatlist sync.Mutex

var Mutexspeakcdtime sync.Mutex

var Mutexproclamationlist sync.Mutex

var Mutexdirtywordlist sync.Mutex

var Mutexfunctionmanagement sync.Mutex

var Mutexmutilangerrormsg sync.Mutex

const Maxchathistory int = 30

const Oncechathistorylong int64 = 1000 * 60 * 60 * 24 * 30

const Speakcdtime int64 = 1000 * 1

var Rooms = make(map[string]map[string]Roomclient)               //key1 : roomUuid  key2 : loginUuid
var Roomsinfo = make(map[string]socket.Roominfo)                 //key : roomUuid
var Usersinfo = make(map[string]socket.User)                     //key : userUuid
var Clients = make(map[string]Client)                            //key : loginUuid
var Clientsconnect = make(map[string]map[string]*websocket.Conn) //key1 : userUuid  key2 : loginUuid
var Speakcdlist = make(map[string]int64)                         //key1 : loginUuid

var Blocklist = make(map[string]map[string]string) //key1 : userUuid  key2 : roomUuid
var Grouplist = make(map[string]map[string]string)

var Roomspopulation = make(map[string]map[string]int) //key1 : ip  key2 : roomUuid

var BlockchatList = make(map[string]map[string]int64) //key1 : blockUserUuid  key2 : blocktargetUuid

var Proclamationlist = make(map[string]map[string]socket.Proclamation) //key1 : roomUuid  key2 : proclamationUuid

var Dirtywordlist = make(map[string]string) //key1 : dirtyWordUuid

var Functionmanagement = make(map[string]string) //key1 : functionName

var UserfriendList = make(map[string]map[string]socket.Friendplatform) //key1 : userUuid  key2 : targetUserUuid

var Mutilangerrormsg = make(map[string]map[string]string) //key1 : lang  key2 : uuid
