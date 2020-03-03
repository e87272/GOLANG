package lineupData

import (
	"sync"

	"../../external/database"
)

var mutexLineup = new(sync.Mutex)

var lineupMap = make(map[int]string) //key: namiID  val:lineupUuid

func getLineupUuidByNami(namiId int) (string, bool) {
	mutexLineup.Lock()
	defer mutexLineup.Unlock()
	lineupUuid, ok := lineupMap[namiId]
	if !ok {
		return "", false
	}
	return lineupUuid, true
}

func setLineupUuidByNami(namiId int, lineupUuid string) {
	mutexLineup.Lock()
	defer mutexLineup.Unlock()
	lineupMap[namiId] = lineupUuid
}

func SearchLineupUuidByNami(namiId int) (string, error) {

	var lineupUuid string

	lineupUuid, ok := getLineupUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'lineup' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&lineupUuid)
		if err != nil {
			return "", err
		}

		setLineupUuidByNami(namiId, lineupUuid)
	}

	return lineupUuid, nil
}
