package socket

type Base_C struct {
	Cmd  string `json:"cmd"`
	Idem string `json:"idem"`
}

type Base_R struct {
	Cmd    string    `json:"cmd"`
	Idem   string    `json:"idem"`
	Stamp  string    `json:"stamp"`
	Result string    `json:"result"`
	Exp    Exception `json:"exp"`
}

type Base_B struct {
	Cmd   string `json:"cmd"`
	Stamp string `json:"stamp"`
}

type Exception struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Userplatform Userplatform `json:"userPlatform"`
	Globalrole   string       `json:"globalRole"`
	Vipgroup     string       `json:"vipGroup"`
	Privategroup string       `json:"privateGroup"`
}

type Logininfo struct {
	Platformuuid string `json:"platformUuid"`
	Platform     string `json:"platform"`
	Token        string `json:"token"`
}

type Roominfo struct {
	Roomuuid      string       `json:"roomUuid"`
	Roomname      string       `json:"roomName"`
	Roomtype      string       `json:"roomType"`
	Roomicon      string       `json:"roomIcon"`
	Adminset      string       `json:"adminSet"`
	Ownerplatform Userplatform `json:"ownerPlatform"`
}

type Sendmessage struct {
	Chattarget string `json:"chatTarget"`
	Style      string `json:"style"`
	Message    string `json:"message"`
}

type Proclamation struct {
	Proclamationuuid string `json:"proclamationUuid"`
	Roomuuid         string `json:"roomUuid"`
	Type             string `json:"type"`
	Order            string `json:"order"`
	Apptype          string `json:"appType"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	Style            string `json:"style"`
	Url              string `json:"url"`
}

type Sudoresult struct {
	Shelltarget        string       `json:"shellTarget"`
	Userplatform       Userplatform `json:"userPlatform"`
	Cmd                []string     `json:"cmd"`
	Targetuserplatform Userplatform `json:"targetUserPlatform"`
}

type Userplatform struct {
	Useruuid     string `json:"userUuid"`
	Platformuuid string `json:"platformUuid"`
	Platform     string `json:"platform"`
}

type Chatmessage struct {
	Historyuuid string       `json:"historyUuid"`
	From        Userplatform `json:"from"`
	Stamp       string       `json:"stamp"`
	Message     string       `json:"message"`
	Style       string       `json:"style"`
}

type Roomcore struct {
	Roomuuid string `json:"roomUuid"`
	Roomtype string `json:"roomType"`
}

type Newsidetext struct {
	Targetuserplatform Userplatform  `json:"targetUserPlatform"`
	Message            []Chatmessage `json:"message"`
	Lastmessage        Chatmessage   `json:"lastMessage"`
}

type Friendplatform struct {
	Userplatform Userplatform `json:"userPlatform"`
	State        string       `json:"state"` //friend,inviteFrom,inviteTo
}

type Globalmessage struct {
	Historyuuid string `json:"historyUuid"`
	Station     string `json:"station"`
	Message     string `json:"message"`
}
