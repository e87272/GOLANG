package common

import (
	"sync"

	"../socket"
	"github.com/gorilla/websocket"
)

type Conncore struct {
	Conn      *websocket.Conn
	Connmutex *sync.Mutex
}

type Client struct {
	Room         map[string]socket.Roomcore
	Conncore     Conncore
	Userplatform socket.Userplatform
	Sidetext     map[string]Sidetextplatform
}

type Roomclient struct {
	Conncore     Conncore
	Userplatform socket.Userplatform
}

type Redispubsubroomdata struct {
	RoomUuid string
	Datajson string
}

type Redispubsubuserdata struct {
	Useruuid string
	Datajson string
}

type Chathistory struct {
	Historyuuid        string `json:"historyUuid"`
	Chattarget         string `json:"chatTarget"`
	Myuuid             string `json:"myUuid"`
	Myplatformuuid     string `json:"myPlatformUuid"`
	Myplatform         string `json:"myPlatform"`
	Stamp              string `json:"stamp"`
	Message            string `json:"message"`
	Style              string `json:"style"`
	Ip                 string `json:"ip"`
	Forwardchatmessage string `json:"forwardChatMessage"`
}

type Redispubsubsidetextdata struct {
	Useruuid   string
	Targetuuid string
	Datajson   string
}

type Redispubsubroomsinfo struct {
	Ip        string
	Roomuuid  string
	Usercount int
	Datajson  string
}

type Syncdata struct {
	Synctype string
	Data     string
}

type Redispubsubdeletesidetext struct {
	Useruuid   string
	Targetuuid string
	Datajson   string
}

type Syserrorlog struct {
	Useruuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Stamp    string `json:"stamp"`
}

type Sidetextplatform struct {
	Userplatform socket.Userplatform
	Sidetextuuid string
}

type Redispubsubinvitedata struct {
	Userfriend   socket.Friendplatform
	Targetfriend socket.Friendplatform
	Isfriend     bool
	Datajson     string
}

type Redispubsubfrienddeletedata struct {
	Useruuid   string
	Targetuuid string
}

type Redispubsubclearusermsgdata struct {
	Roomuuid   string
	Targetuuid string
}

type GlobalMessage struct {
	Station      string
	Content      string
	Ticktime     int64
	Endtime      int64
	Timeinterval int64
}
