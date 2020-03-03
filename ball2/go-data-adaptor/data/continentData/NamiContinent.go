package continentData

import (
	"sync"

	"../../external/database"
)

var mutexContinent = new(sync.Mutex)

var namiContinentMap = make(map[int]string) //key: namiID  val:continentUuid

func getContinentUuidByNami(namiId int) (string, bool) {
	mutexContinent.Lock()
	defer mutexContinent.Unlock()
	continentUuid, ok := namiContinentMap[namiId]
	if !ok {
		return "", false
	}
	return continentUuid, true
}

func setContinentUuidByNami(namiId int, continentUuid string) {
	mutexContinent.Lock()
	defer mutexContinent.Unlock()
	namiContinentMap[namiId] = continentUuid
}

func SearchContinentUuidByNami(namiId int) (string, error) {

	var continentUuid string

	continentUuid, ok := getContinentUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'continent' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&continentUuid)
		if err != nil {
			return "", err
		}

		setContinentUuidByNami(namiId, continentUuid)
	}

	return continentUuid, nil
}
