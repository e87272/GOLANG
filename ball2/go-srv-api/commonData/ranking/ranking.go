package ranking

import (
	"sync"

	commonData ".."
)

var mutexRankMap = new(sync.Mutex)
var rank = make(map[string][]commonData.RankInfo)

func SetRankMap(rankMap map[string][]commonData.RankInfo) {
	mutexRankMap.Lock()
	defer mutexRankMap.Unlock()

	rank = rankMap
}

func GetRank(rankType string) ([]commonData.RankInfo, bool) {
	mutexRankMap.Lock()
	defer mutexRankMap.Unlock()

	list, ok := rank[rankType]
	return list, ok
}

func SetRank(rankType string, list []commonData.RankInfo) {
	mutexRankMap.Lock()
	defer mutexRankMap.Unlock()

	rank[rankType] = list
}

func DeleteRank(rankType string) {
	mutexRankMap.Lock()
	defer mutexRankMap.Unlock()

	delete(rank, rankType)
}
