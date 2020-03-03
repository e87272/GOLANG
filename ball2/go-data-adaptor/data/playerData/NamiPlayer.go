package playerData

import (
	"sync"

	"../../external/database"
)

var mutexPlayer = new(sync.Mutex)

var namiPlayerMap = make(map[int]string) //key: namiID  val:playerUuid

func getPlayerUuidByNami(namiId int) (string, bool) {
	mutexPlayer.Lock()
	defer mutexPlayer.Unlock()
	playerUuid, ok := namiPlayerMap[namiId]
	if !ok {
		return "", false
	}
	return playerUuid, true
}

func setPlayerUuidByNami(namiId int, playerUuid string) {
	mutexPlayer.Lock()
	defer mutexPlayer.Unlock()
	namiPlayerMap[namiId] = playerUuid
}

func SearchPlayerUuidByNami(namiId int) (string, error) {

	var playerUuid string

	playerUuid, ok := getPlayerUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'player' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&playerUuid)
		if err != nil {
			return "", err
		}

		setPlayerUuidByNami(namiId, playerUuid)
	}

	return playerUuid, nil
}
