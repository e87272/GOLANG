package league

import (
	"sync"

	commonData ".."
)

var mutexLeagueList = new(sync.Mutex)

var leagueList = make(map[string]commonData.LeagueInfo)

func GetLeagueList() map[string]commonData.LeagueInfo {

	mutexLeagueList.Lock()
	defer mutexLeagueList.Unlock()

	return leagueList
}

func GetLeagueInfo(leagueUuid string) (commonData.LeagueInfo, bool) {

	mutexLeagueList.Lock()
	defer mutexLeagueList.Unlock()

	leagueInfo, ok := leagueList[leagueUuid]

	return leagueInfo, ok
}

func SetLeagueInfo(leagueInfo commonData.LeagueInfo) {

	mutexLeagueList.Lock()
	defer mutexLeagueList.Unlock()

	leagueList[leagueInfo.LeagueCore.LeagueUuid] = leagueInfo

	return
}

func SetLeagueInfoTeam(leagueUuid string, teamCore commonData.TeamCore) {

	mutexLeagueList.Lock()
	defer mutexLeagueList.Unlock()

	leagueInfo, ok := leagueList[leagueUuid]

	if ok {
		leagueInfo.TeamList[teamCore.TeamUuid] = teamCore
	}

	// log.Printf("SetLeagueInfoTeam leagueList: %+v\n", leagueList)

	return
}
