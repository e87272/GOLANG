package namiAdaptor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"../commonFunc"
	"../data"
	"../data/gameData"
	"../data/leagueData"
	"../data/seasonData"
	"../data/stageData"
	"../data/teamData"
	"../external/database"
)

func seasonInfoUpdate(seasonId int, season string) (string, error) {

	var seasonInfo data.NamiSeasonInfo
	dataMap := make(map[string]string)
	dataMap["id"] = strconv.FormatInt(int64(seasonId), 10)

	resp, err := commonFunc.GetNamiApi("/season/detail", dataMap)
	if err != nil {
		log.Printf("seasonInfoUpdate err : %+v\n", err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return "", err
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &seasonInfo); err != nil {
		log.Printf("gameList json.Unmarshal err : %+v\n", err)
		return "", err
	}

	// log.Printf("Stages : %+v\n", len(seasonInfo.Stages))
	// log.Printf("Teams : %+v\n", len(seasonInfo.Teams))
	// log.Printf("Matches : %+v\n", len(seasonInfo.Matches))

	for _, stageInfo := range seasonInfo.Stages {
		stageData.SetStageData(stageInfo.Id, stageInfo)
	}

	leagueUuid, err := leagueData.SearchLeagueUuidByNami(seasonInfo.Competition.Id)

	if err != nil {
		log.Printf("SearchLeagueUuidByNami leagueUuid err : %+v\n", err)
		return "", err
	}

	seasonUuid, err := seasonData.SearchSeasonUuidByNami(seasonId)

	if err != nil {

		if err != database.ErrNoRows {
			log.Printf("seasonInfoUpdate select err : %+v\n", err)
			return "", err
		}

		seasonUuid = commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `season` (`uuid` ,`league_uuid`,`season_year`) VALUES (?, ?, ?)",
			seasonUuid, leagueUuid, season,
		)
		if err != nil {
			log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
			return "", err
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "season", seasonUuid, "nami", seasonId,
		)

		if err != nil {
			log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
			return "", err
		}

		if err != nil {
			log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
			return "", err
		}
	}

	for _, teamInfo := range seasonInfo.Teams {
		teamUuid, err := teamData.SearchTeamUuidByNami(teamInfo.Id)
		if err != nil {
			log.Printf("teamInfo : %+v\n", teamInfo.Id)
			teamInfoUpdate(teamInfo.Id, seasonInfo.Competition.Id)

			teamUuid, err = teamData.SearchTeamUuidByNami(teamInfo.Id)
			if err != nil {
				log.Printf("teamInfo : %+v\n", teamInfo.Id)
				return "", err
			}
		} else {
			go teamInfoUpdate(teamInfo.Id, seasonInfo.Competition.Id)
		}

		seasonTeamUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `season_team_list` (`uuid` ,`seasonUuid` ,`teamUuid` ) VALUES (?, ?, ?)",
			seasonTeamUuid, seasonUuid, teamUuid,
		)
		time.Sleep(time.Duration(1) * time.Second)
	}

	for _, gameInfo := range seasonInfo.Matches {
		_, err := gameData.SearchGameUuidByNami(gameInfo.Id)
		if err != nil {
			// log.Printf("seasonInfo.Matches gameInfo.Id : %+v\n", gameInfo.Id)
			gameInfoUpdate(gameInfo.Id, seasonUuid, gameInfo.Round.Stage_id, seasonInfo.Competition.Short_name_zh)
		} else {
			go gameInfoUpdate(gameInfo.Id, seasonUuid, gameInfo.Round.Stage_id, seasonInfo.Competition.Short_name_zh)
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	// for _, promotion := range seasonInfo.Table.Promotions {
	// 	seasonData.SetPromotionData(promotion.Id, promotion)
	// }
	// isPromotion := true
	// if len(seasonInfo.Table.Promotions) == 0 {
	// 	isPromotion = false
	// }
	// for _, seasonRankInfo := range seasonInfo.Table.Tables {
	// 	var class string
	// 	if !isPromotion {
	// 		class = seasonInfo.Competition.Short_name_zht + data.GroupNum[seasonRankInfo.Group-1] + "組"
	// 	}

	// }

	return seasonUuid, nil

}
