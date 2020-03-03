package game

import (
	"sync"

	commonData ".."
	"../../external/stamp"
)

var mutexGame = new(sync.Mutex)
var gameMap = map[string]commonData.GameInfo{}           //key1: gameUuid
var allGameUuidList = map[int64][]string{}               //key1: date
var leagueGameUuidList = map[string]map[int64][]string{} //key1: leagueUuid, key2 : date

func getAllGameList(date int64) ([]commonData.GameInfo, bool) {

	gameList := []commonData.GameInfo{}

	uuidList, ok := allGameUuidList[stamp.Date(date)]
	if !ok {
		return gameList, false
	}

	for _, gameUuid := range uuidList {
		gameList = append(gameList, gameMap[gameUuid])
	}
	return gameList, true
}

func getLeagueGameList(date int64, leagueUuid string) ([]commonData.GameInfo, bool) {

	gameList := []commonData.GameInfo{}

	uuidMap, ok := leagueGameUuidList[leagueUuid]
	if !ok {
		return gameList, false
	}
	uuidList, ok := uuidMap[stamp.Date(date)]
	if !ok {
		return gameList, false
	}

	for _, gameUuid := range uuidList {
		gameList = append(gameList, gameMap[gameUuid])
	}
	return gameList, true
}

func setGame(gameInfo commonData.GameInfo) {

	_, ok := gameMap[gameInfo.Uuid]
	gameMap[gameInfo.Uuid] = gameInfo
	if ok {
		return
	}
	date := stamp.Date(gameInfo.StartTime)

	uuidList, ok := allGameUuidList[date]
	if !ok {
		uuidList = []string{}
	}
	uuidList = append(uuidList, gameInfo.Uuid)
	allGameUuidList[date] = uuidList

	uuidMap, ok := leagueGameUuidList[gameInfo.LeagueUuid]
	if !ok {
		uuidMap = map[int64][]string{}
	}
	uuidList, ok = uuidMap[date]
	if !ok {
		uuidList = []string{}
	}
	uuidList = append(uuidList, gameInfo.Uuid)
	uuidMap[date] = uuidList
	leagueGameUuidList[gameInfo.LeagueUuid] = uuidMap
}

func GetGameInfo(gameUuid string) (commonData.GameInfo, bool) {

	mutexGame.Lock()
	defer mutexGame.Unlock()

	gameInfo, ok := gameMap[gameUuid]
	return gameInfo, ok
}

func GetGameList(date int64, leagueUuid string) ([]commonData.GameInfo, bool) {

	mutexGame.Lock()
	defer mutexGame.Unlock()

	if leagueUuid == "" {
		return getAllGameList(date)
	}
	return getLeagueGameList(date, leagueUuid)
}

func SetGameInfo(gameInfo commonData.GameInfo) {

	mutexGame.Lock()
	defer mutexGame.Unlock()

	setGame(gameInfo)
}

func SetGameInfoBatch(gameList []commonData.GameInfo) {

	mutexGame.Lock()
	defer mutexGame.Unlock()

	for _, gameInfo := range gameList {
		setGame(gameInfo)
	}
}

func ResetGameMap() {

	mutexGame.Lock()
	defer mutexGame.Unlock()

	gameMap = map[string]commonData.GameInfo{}
	allGameUuidList = map[int64][]string{}
	leagueGameUuidList = map[string]map[int64][]string{}
}
