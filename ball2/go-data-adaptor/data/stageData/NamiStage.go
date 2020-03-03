package stageData

import (
	"sync"

	data ".."
)

var mutexTeam = new(sync.Mutex)

var stageList = make(map[int]data.NamiStage)

func GetStageData(stageId int) (data.NamiStage, bool) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	stage, ok := stageList[stageId]
	if !ok {
		return stage, false
	}
	return stage, true
}

func SetStageData(stageId int, stage data.NamiStage) {
	mutexTeam.Lock()
	defer mutexTeam.Unlock()
	stageList[stageId] = stage
}
