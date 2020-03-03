package countryData

import (
	"sync"

	"../../external/database"
)

var mutexCountry = new(sync.Mutex)

var namiCountryMap = make(map[int]string) //key: namiID  val:countryUuid

func getCountryUuidByNami(namiId int) (string, bool) {
	mutexCountry.Lock()
	defer mutexCountry.Unlock()
	countryUuid, ok := namiCountryMap[namiId]
	if !ok {
		return "", false
	}
	return countryUuid, true
}

func setCountryUuidByNami(namiId int, countryUuid string) {
	mutexCountry.Lock()
	defer mutexCountry.Unlock()
	namiCountryMap[namiId] = countryUuid
}

func SearchCountryUuidByNami(namiId int) (string, error) {

	var countryUuid string

	countryUuid, ok := getCountryUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'country' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&countryUuid)
		if err != nil {
			return "", err
		}

		setCountryUuidByNami(namiId, countryUuid)
	}

	return countryUuid, nil
}
