package commonFunc

import (
	"log"
	"net"
	"strings"

	_ "net/http/pprof"

	"../commonData"
	"../commonData/game"
	"../commonData/homeBanner"
	"../commonData/launchScreen"
	"../commonData/league"
	"../commonData/team"
	"../external/database"
	"../external/ginEngine"
	"../external/stamp"

	"github.com/gin-gonic/gin"
)

func MyIp() string {

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	// log.Printf("localAddr : %+v\n", localAddr)

	idx := strings.LastIndex(localAddr, ":")

	myIp := localAddr[0:idx]

	return myIp
}

func SessionCheck(context *gin.Context) (string, bool) {

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		return "", false
	}

	return userUuid, true
}

func UserInfoSearch(userUuid string) (commonData.UserInfo, bool, commonData.Exception) {

	exception := commonData.Exception{}
	userInfo, ok := GetRedisUserInfo(userUuid)

	if !ok {

		row := database.QueryRow("select uuid,phone,nickname,icon from account where uuid = ? ",
			userUuid,
		)
		err := row.Scan(&userInfo.UserUuid, &userInfo.Phone, &userInfo.Nickname, &userInfo.Icon)

		if err != nil {
			exception = Exception("COMMON_USERINFOSEARCH_USERINFO_READ_ERROR", userUuid, err)
			return commonData.UserInfo{}, false, exception
		}

		SetRedisUserInfo(userUuid, userInfo)
	}
	return userInfo, true, exception
}

func Exception(msg string, userUuid string, err error) commonData.Exception {

	if msg == "" {
		return commonData.Exception{}
	}

	var code = EsSysErrorLog(msg, userUuid, err)

	return commonData.Exception{Code: code, Message: msg}
}

func ClientLeagueInfoSearch(leagueUuid string) (commonData.ClientLeagueInfo, bool, commonData.Exception) {

	clientLeagueInfo := commonData.ClientLeagueInfo{}
	clientLeagueInfo.TeamList = make(map[string]commonData.TeamCore)
	exception := commonData.Exception{}

	leagueInfo, ok := league.GetLeagueInfo(leagueUuid)

	if !ok {
		exception = Exception("COMMON_CLIENTCLASSINFOSEARCH_CLASS_UUID_ERROR", "", nil)
		return commonData.ClientLeagueInfo{}, false, exception

	}
	clientLeagueInfo.LeagueCore = leagueInfo.LeagueCore

	log.Printf("ClientLeagueInfoSearch leagueInfo.TeamList: %+v\n", leagueInfo.TeamList)

	for _, teamCore := range leagueInfo.TeamList {
		clientLeagueInfo.TeamList[teamCore.TeamUuid] = teamCore
	}

	return clientLeagueInfo, true, exception
}

func ClientTeamInfoSearch(teamUuid string) (commonData.TeamInfo, bool, commonData.Exception) {

	exception := commonData.Exception{}

	teamInfo, ok := team.GetTeamInfo(teamUuid)

	if !ok {
		exception = Exception("COMMON_CLIENTTEAMINFOSEARCH_CLASS_UUID_ERROR", "", nil)
		return commonData.TeamInfo{}, false, exception

	}

	return teamInfo, true, exception
}

func ClientLeagueList() ([]commonData.ClientLeagueList, bool, commonData.Exception) {

	exception := commonData.Exception{}
	clientLeagueList := []commonData.ClientLeagueList{}

	leagueList := league.GetLeagueList()

	for _, leagueInfo := range leagueList {
		clientLeagueInfo := commonData.ClientLeagueList{}
		clientLeagueInfo.LeagueCore = leagueInfo.LeagueCore
		clientLeagueInfo.Sequence = leagueInfo.Sequence
		clientLeagueInfo.TeamCount = len(leagueInfo.TeamList)

		clientLeagueList = append(clientLeagueList, clientLeagueInfo)
	}

	return clientLeagueList, true, exception
}

func ClientTeamList() ([]commonData.ClientTeamList, bool, commonData.Exception) {

	exception := commonData.Exception{}
	clientTeamList := []commonData.ClientTeamList{}

	teamList := team.GetTeamList()

	for _, teamInfo := range teamList {
		clientTeamInfo := commonData.ClientTeamList{}
		clientTeamInfo.TeamCore = teamInfo.TeamCore
		clientTeamInfo.PlayerCount = len(teamInfo.PlayerList)

		clientTeamList = append(clientTeamList, clientTeamInfo)
	}

	return clientTeamList, true, exception
}

func GetLaunchScreen() []commonData.Announcement {

	return launchScreen.GetLaunchScreen()
}

func GetHomeBanner() []commonData.Announcement {

	return homeBanner.GetHomeBanner()
}

func GetSchedule(date int64, leagueUuid string) []commonData.GameInfo {

	gameList, ok := game.GetGameList(date, leagueUuid)
	if !ok {
		queryStartTime := stamp.Date(date)
		queryEndTime := queryStartTime + stamp.Day
		QuerySchedule(queryStartTime, queryEndTime)
		gameList, _ = game.GetGameList(date, leagueUuid)
	}

	return gameList
}

func GetAllSchedule(leagueUuid string) []commonData.GameInfo {

	allGameList := []commonData.GameInfo{}
	today := stamp.Today()

	date := today - stamp.Day
	gameList, ok := game.GetGameList(date, leagueUuid)
	if !ok {
		queryStartTime := date
		queryEndTime := queryStartTime + stamp.Day
		QuerySchedule(queryStartTime, queryEndTime)
		gameList, _ = game.GetGameList(date, leagueUuid)
	}
	for _, gameInfo := range gameList {
		if gameInfo.State != "完场" {
			allGameList = append(allGameList, gameInfo)
		}
	}

	for i := 0; i < 2; i++ {
		date = today + int64(i)*stamp.Day
		gameList, ok = game.GetGameList(date, leagueUuid)
		if !ok {
			queryStartTime := date
			queryEndTime := queryStartTime + stamp.Day
			QuerySchedule(queryStartTime, queryEndTime)
			gameList, _ = game.GetGameList(date, leagueUuid)
		}
		allGameList = append(allGameList, gameList...)
	}

	return allGameList
}

func GetHotSchedule() []commonData.GameInfo {

	return GetAllSchedule("")
}
