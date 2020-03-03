package team

import (
	"sync"

	commonData ".."
)

var mutexTeamList = new(sync.Mutex)

var teamList = make(map[string]commonData.TeamInfo)

func GetTeamList() map[string]commonData.TeamInfo {

	mutexTeamList.Lock()
	defer mutexTeamList.Unlock()

	return teamList
}

func GetTeamInfo(teamUuid string) (commonData.TeamInfo, bool) {

	mutexTeamList.Lock()
	defer mutexTeamList.Unlock()

	teamInfo, ok := teamList[teamUuid]

	return teamInfo, ok
}

func SetTeamInfo(teamUuid string, teamInfo commonData.TeamInfo) {

	mutexTeamList.Lock()
	defer mutexTeamList.Unlock()

	teamList[teamUuid] = teamInfo

	return
}

func SetTeamInfoPlayer(teamUuid string, playerUuid string, playerInfo commonData.PlayerInfo) {

	mutexTeamList.Lock()
	defer mutexTeamList.Unlock()

	teamInfo, ok := teamList[teamUuid]
	if ok {
		teamInfo.PlayerList[playerUuid] = playerInfo
	}

	return
}
