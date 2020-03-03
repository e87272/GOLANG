package leagueData

import (
	"sync"

	data ".."
	"../../external/database"
)

var mutexLeague = new(sync.Mutex)

var namiLeagueMap = make(map[int]string) //key: namiID  val:leagueUuid

var leagueList = make(map[int]data.NamiLeague)

func getLeagueUuidByNami(namiId int) (string, bool) {
	mutexLeague.Lock()
	defer mutexLeague.Unlock()
	leagueUuid, ok := namiLeagueMap[namiId]
	if !ok {
		return "", false
	}
	return leagueUuid, true
}

func setLeagueUuidByNami(namiId int, leagueUuid string) {
	mutexLeague.Lock()
	defer mutexLeague.Unlock()
	namiLeagueMap[namiId] = leagueUuid
}

func SearchLeagueUuidByNami(namiId int) (string, error) {

	var leagueUuid string

	leagueUuid, ok := getLeagueUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'league' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&leagueUuid)
		if err != nil {
			return "", err
		}

		setLeagueUuidByNami(namiId, leagueUuid)
	}

	return leagueUuid, nil
}

func GetLeagueData(leagueId int) (data.NamiLeague, bool) {
	mutexLeague.Lock()
	defer mutexLeague.Unlock()
	league, ok := leagueList[leagueId]
	if !ok {
		return league, false
	}
	return league, true
}

func SetLeagueData(leagueId int, league data.NamiLeague) {
	mutexLeague.Lock()
	defer mutexLeague.Unlock()
	leagueList[leagueId] = league
}
