package teamData

import (
	"sync"

	data ".."
	"../../external/database"
)

var mutexTeam = new(sync.Mutex)

var namiTeamMap = make(map[int]string) //key: namiID  val:teamUuid

var namiTeamList data.NamiTeamList

func getTeamUuidByNami(namiId int) (string, bool) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	teamUuid, ok := namiTeamMap[namiId]
	if !ok {
		return "", false
	}
	return teamUuid, true
}

func setTeamUuidByNami(namiId int, teamUuid string) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	namiTeamMap[namiId] = teamUuid
}

func SearchTeamUuidByNami(namiId int) (string, error) {

	var teamUuid string

	teamUuid, ok := getTeamUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'team' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&teamUuid)
		if err != nil {
			return "", err
		}

		setTeamUuidByNami(namiId, teamUuid)
	}

	return teamUuid, nil
}

func GetNamiTeamListData(namiId int) (data.NamiTeamListInfo, bool) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()

	for _, teamInfo := range namiTeamList {
		if teamInfo.Id == namiId {
			return teamInfo, true
		}
	}
	return data.NamiTeamListInfo{}, false
}

func SetNamiTeamListData(teamList data.NamiTeamList) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	namiTeamList = teamList
}
