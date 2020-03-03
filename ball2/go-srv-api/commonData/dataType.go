package commonData

type Exception struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UserInfo struct {
	UserUuid string `json:"userUuid"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Icon     string `json:"icon"`
}

type SysLog struct {
	ApiName string `json:"apiName"`
	Message string `json:"message"`
	Stamp   string `json:"stamp"`
}

type SysErrorLog struct {
	UserUuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Stamp    string `json:"stamp"`
}

type LeagueCore struct {
	LeagueUuid string `json:"leagueUuid"`
	LeagueName string `json:"leagueName"`
}

type LeagueInfo struct {
	LeagueCore LeagueCore          `json:"leagueCore"`
	Sequence   int                 `json:"sequence"`
	TeamList   map[string]TeamCore `json:"teamList"`
}

type TeamCore struct {
	TeamUuid string `json:"teamUuid"`
	Name     string `json:"name"`
}

type TeamInfo struct {
	TeamCore   TeamCore              `json:"teamCore"`
	Manager    string                `json:"manager"`
	Venue      string                `json:"venue"`
	Found      string                `json:"found"`
	PlayerList map[string]PlayerInfo `json:"playerList"`
}

type PlayerInfo struct {
	PlayerUuid  string `json:"playerUuid"`
	PlayerName  string `json:"playerName"`
	CountryUuid string `json:"countryUuid"`
}

type Announcement struct {
	Uuid     string `json:"uuid"`
	Type     string `json:"type"`
	Sequence string `json:"sequence"`
	Content  string `json:"content"`
	Url      string `json:"url"`
}

type GameInfo struct {
	Uuid          string `json:"uuid"`
	HomeTeamUuid  string `json:"homeTeamUuid"`
	GuestTeamUuid string `json:"guestTeamUuid"`
	StartTime     int64  `json:"startTime"`
	LeagueUuid    string `json:"leagueUuid"`
	HomeScore     int    `json:"homeScore"`
	GuestScore    int    `json:"guestScore"`
	Subtitle      string `json:"subtitle"`
	State         string `json:"state"`
}

type RankInfo struct {
	TeamCore TeamCore `json:"teamCore"`
	Score    int      `json:"score"`
}
