package player

import (
	"sync"

	commonData ".."
)

var mutexPlayerList = new(sync.Mutex)

var playerList = make(map[string]commonData.PlayerInfo)

func GetPlayerInfo(playerUuid string) (commonData.PlayerInfo, bool) {

	mutexPlayerList.Lock()
	defer mutexPlayerList.Unlock()

	playerInfo, ok := playerList[playerUuid]

	return playerInfo, ok
}

func SetPlayerInfo(playerInfo commonData.PlayerInfo) {

	mutexPlayerList.Lock()
	defer mutexPlayerList.Unlock()

	playerList[playerInfo.PlayerUuid] = playerInfo

	return
}
