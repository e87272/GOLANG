package commonData

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Exception struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SysLog struct {
	UserUuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Stamp    string `json:"stamp"`
}
type SysErrorLog struct {
	UserUuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Stamp    string `json:"stamp"`
}

type ConnCore struct {
	Conn      *websocket.Conn
	Connmutex *sync.Mutex
	LoginUuid string
}

type UserInfo struct {
	UserUuid string `json:"userUuid"`
}

type Client struct {
	Room     map[string]string `json:"room"` //key:roomUuid val:roomUuid
	ConnCore ConnCore          `json:"connCore"`
	UserInfo UserInfo          `json:"userInfo"`
}

type RoomClient struct {
	ConnCore ConnCore `json:"connCore"`
	UserUuid string   `json:"userUuid"`
}

type ChatMessage struct {
	HistoryUuid string   `json:"historyUuid"`
	From        UserInfo `json:"from"`
	Stamp       string   `json:"stamp"`
	Message     string   `json:"message"`
	Style       string   `json:"style"`
}

type SendMessage struct {
	ChatTarget string `json:"chatTarget"`
	Style      string `json:"style"`
	Message    string `json:"message"`
}

type ChatHistory struct {
	HistoryUuid string `json:"historyUuid"`
	Chattarget  string `json:"chatTarget"`
	MyUuid      string `json:"myUuid"`
	Style       string `json:"style"`
	Message     string `json:"message"`
	Stamp       string `json:"stamp"`
}
