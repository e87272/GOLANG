package common

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"

	"../socket"
)

var Elasticclient *elastic.Client

var Redisclient *redis.Client

var Mutexiplist = new(sync.Mutex)

var Mutexrooms = new(sync.Mutex)

var Mutexroomsinfo = new(sync.Mutex)

var Mutexroomsstation = new(sync.Mutex)

var Mutexusersinfo = new(sync.Mutex)

var Mutexroomspopulation = new(sync.Mutex)

var Mutexuserfriendlist = new(sync.Mutex)

var Mutexclients = new(sync.Mutex)

var Mutexclientsconnect = new(sync.Mutex)

var Mutexconnect = new(sync.Mutex)

var Mutexelastic = new(sync.Mutex)

var Mutexredis = new(sync.Mutex)

var Mutexblockchatlist = new(sync.Mutex)

var MutexblockipList = new(sync.Mutex)

var Mutexblocknewuserlist = new(sync.Mutex)

var Mutexspeakcdtime = new(sync.Mutex)

var Mutexproclamationlist = new(sync.Mutex)

var Mutexdirtywordlist = new(sync.Mutex)

var Mutexfunctionmanagement = new(sync.Mutex)

var Mutexmutilangerrormsg = new(sync.Mutex)

const Maxchathistory int = 30

const Oncechathistorylong int64 = 1000 * 60 * 60 * 24 * 30

const Speakcdtime int64 = 1000 * 1

const Packetdroptime int64 = 1000 * 1

const Newusercdtime int64 = 1000 * 60 * 60 * 24

var Iplist = make(map[string]string) //key1 : loginUuid val:ip

var Rooms = make(map[string]map[string]Roomclient)        //key1 : roomUuid  key2 : loginUuid
var Roomsinfo = make(map[string]socket.Roominfo)          //key : roomUuid
var Roomsstation = make(map[string]string)                //key1 : station + "_" + Memberuuid val:roomUuid
var Usersinfo = make(map[string]socket.User)              //key : userUuid
var Clients = make(map[string]Client)                     //key : loginUuid
var Clientsconnect = make(map[string]map[string]Conncore) //key1 : userUuid  key2 : loginUuid
var Speakcdlist = make(map[string]int64)                  //key1 : loginUuid

var Blocklist = make(map[string]map[string]string) //key1 : userUuid  key2 : roomUuid
var Grouplist = make(map[string]map[string]string)

var Roomspopulation = make(map[string]map[string]int) //key1 : ip  key2 : roomUuid

var BlockchatList = make(map[string]map[string]int64) //key1 : blockUserUuid  key2 : blocktargetUuid

var BlockipList = make(map[string]int64) //key1 : ip

var Blocknewuserlist = make(map[string]int64) //key1 : UserUuid  key2 : timestamp

var Proclamationlist = make(map[string]map[string]socket.Proclamation) //key1 : roomUuid  key2 : proclamationUuid

var Dirtywordlist = make(map[string]string) //key1 : dirtyWordUuid

var Functionmanagement = make(map[string]string) //key1 : functionName

var UserfriendList = make(map[string]map[string]socket.Friendplatform) //key1 : userUuid  key2 : targetUserUuid

var Mutilangerrormsg = make(map[string]map[string]string) //key1 : lang  key2 : uuid
