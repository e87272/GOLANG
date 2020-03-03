package gameData

import (
	"sync"

	data ".."
	"../../external/database"
)

var mutexGame = new(sync.Mutex)

var namiGameMap = make(map[int]string)       //key: namiID  val:gameUuid
var gameAnaysisMap = make(map[string]string) //key: gameUuid  val:gameUuid

var namiGameListMap = make(map[string]data.NamiGameList)

func getGameUuidByNami(namiId int) (string, bool) {
	mutexGame.Lock()
	defer mutexGame.Unlock()
	gameUuid, ok := namiGameMap[namiId]
	if !ok {
		return "", false
	}
	return gameUuid, true
}

func setGameUuidByNami(namiId int, gameUuid string) {
	mutexGame.Lock()
	defer mutexGame.Unlock()
	namiGameMap[namiId] = gameUuid
}

func SearchGameUuidByNami(namiId int) (string, error) {

	var gameUuid string

	gameUuid, ok := getGameUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'game' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&gameUuid)
		if err != nil {
			return "", err
		}

		setGameUuidByNami(namiId, gameUuid)
	}

	return gameUuid, nil
}

func GetNamiGameListData(dateStr string) data.NamiGameList {
	mutexGame.Lock()
	defer mutexGame.Unlock()

	return namiGameListMap[dateStr]
}

func SetNamiGameListData(dateStr string, gameList data.NamiGameList) {
	mutexGame.Lock()
	defer mutexGame.Unlock()
	namiGameListMap[dateStr] = gameList
}

func getGameAnaysisUuid(gameUuid string) bool {
	mutexGame.Lock()
	defer mutexGame.Unlock()
	_, ok := gameAnaysisMap[gameUuid]
	if !ok {
		return false
	}
	return true
}

func setGameAnaysisUuid(gameUuid string) {
	mutexGame.Lock()
	defer mutexGame.Unlock()
	gameAnaysisMap[gameUuid] = gameUuid
}

func CheckGameAnaysisUuid(gameUuid string) bool {

	ok := getGameAnaysisUuid(gameUuid)

	if !ok {
		var uuid string
		row := database.QueryRow("SELECT `game_uuid` FROM `game_analysis` WHERE `game_uuid` = ?",
			gameUuid,
		)
		err := row.Scan(&uuid)
		if err != nil {
			return false
		}

		setGameAnaysisUuid(gameUuid)
	}

	return true
}
