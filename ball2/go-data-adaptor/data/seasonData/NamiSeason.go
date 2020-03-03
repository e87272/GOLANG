package seasonData

import (
	"sync"

	data ".."
	"../../external/database"
)

var mutexSeason = new(sync.Mutex)

var mutexSeasonRank = new(sync.Mutex)

var namiSeasonMap = make(map[int]string) //key: namiID  val:seasonUuid

var namiSeasonRankMap = make(map[int]string) //key: namiID  val:seasonUuid

func getSeasonUuidByNami(namiId int) (string, bool) {
	mutexSeason.Lock()
	defer mutexSeason.Unlock()
	seasonUuid, ok := namiSeasonMap[namiId]
	if !ok {
		return "", false
	}
	return seasonUuid, true
}

func setSeasonUuidByNami(namiId int, seasonUuid string) {
	mutexSeason.Lock()
	defer mutexSeason.Unlock()
	namiSeasonMap[namiId] = seasonUuid
}

func SearchSeasonUuidByNami(namiId int) (string, error) {

	var seasonUuid string

	seasonUuid, ok := getSeasonUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'season' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&seasonUuid)
		if err != nil {
			return "", err
		}

		setSeasonUuidByNami(namiId, seasonUuid)
	}

	return seasonUuid, nil
}

var mutexPromotion = new(sync.Mutex)

var promotionList = make(map[int]data.NamiPromotion)

func GetPromotionData(promotionId int) (data.NamiPromotion, bool) {
	mutexPromotion.Lock()
	defer mutexPromotion.Unlock()
	promotion, ok := promotionList[promotionId]
	if !ok {
		return promotion, false
	}
	return promotion, true
}

func SetPromotionData(promotionId int, promotion data.NamiPromotion) {
	mutexPromotion.Lock()
	defer mutexPromotion.Unlock()
	promotionList[promotionId] = promotion
}

func getSeasonRankUuidByNami(namiId int) (string, bool) {
	mutexSeason.Lock()
	defer mutexSeason.Unlock()
	seasonUuid, ok := namiSeasonMap[namiId]
	if !ok {
		return "", false
	}
	return seasonUuid, true
}

func setSeasonRankUuidByNami(namiId int, SeasonRankUuid string) {
	mutexSeasonRank.Lock()
	defer mutexSeasonRank.Unlock()
	namiSeasonRankMap[namiId] = SeasonRankUuid
}

func SearchSeasonRankUuidByNami(namiId int) (string, error) {

	var SeasonRankUuid string

	SeasonRankUuid, ok := getSeasonRankUuidByNami(namiId)

	if !ok {
		row := database.QueryRow("SELECT `my_id` FROM `corporation_id` WHERE `type` = 'SeasonRank' AND `source_name` = 'nami' AND `source_id` = ?",
			namiId,
		)
		err := row.Scan(&SeasonRankUuid)
		if err != nil {
			return "", err
		}

		setSeasonRankUuidByNami(namiId, SeasonRankUuid)
	}

	return SeasonRankUuid, nil
}
