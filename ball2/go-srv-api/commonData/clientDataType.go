package commonData

type ClientLeagueList struct {
	LeagueCore LeagueCore `json:"leagueCore"`
	Sequence   int        `json:"sequence"`
	TeamCount  int        `json:"teamCount"`
}

type ClientLeagueInfo struct {
	LeagueCore LeagueCore          `json:"leagueCore"`
	TeamList   map[string]TeamCore `json:"teamList"`
}

type ClientTeamList struct {
	TeamCore    TeamCore `json:"teamCore"`
	PlayerCount int      `json:"playerCount"`
}
