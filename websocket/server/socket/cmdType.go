package socket
type Base_C struct {
	Cmd		string   `json:"cmd"`
	Idem	string   `json:"idem"`
}

type Base_R struct {
	Cmd		string    `json:"cmd"`
	Idem	string    `json:"idem"`
	Stamp	string    `json:"stamp"`
	Result	string    `json:"result"`
	Exp	    Exception `json:"exp"`
}

type Base_B struct {
	Cmd		string      `json:"cmd"`
	Stamp	string      `json:"stamp"`
}

type Exception struct {
	Code	string    `json:"code"`
	Message	string    `json:"message"`
}

type User struct {
	Id			string    `json:"id"`
	Nickname	string    `json:"nickName"`
	Icon		string    `json:"icon"`
	Role		string    `json:"role"`
	Status		string    `json:"status"`
}

type Chat struct {
	Text	    string    `json:"text"`
	Style	    string    `json:"style"`
	Roominfo	RoomInfo    `json:"roomInfo"`
}

type ChatMessage struct {
	From	User      `json:"from"`
	Stamp	string    `json:"stamp"`
	Text	string    `json:"text"`
	Style	string    `json:"style"`
	Roominfo	RoomInfo    `json:"roomInfo"`
}

type Room struct {
	Users          []User    `json:"users"`
	ChatMessages   []ChatMessage    `json:"chatMessages"`
}

type LoginInfo struct{
	Roomname string   `json:"roomName"`
	Nickname string   `json:"nickName"`
}

type RoomInfo struct{
	Roomname string   `json:"roomName"`
}