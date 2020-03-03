package namiAdaptor

import (
	"encoding/json"
	"fmt"
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

func seasonList() {

	fmt.Println(" seasonList...")

	for {
		now := time.Now()
		// 计算下一天
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 13, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		seasonListUpdate()
		seasonList()
	}

}

func seasonListUpdate() {

	var seasonList data.NamiSeasonList

	resp, err := commonFunc.GetNamiApi("/season/list", nil)
	if err != nil {
		log.Printf("seasonListTicker err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &seasonList); err != nil {
		log.Printf("seasonList json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("Areas : %+v\n", len(seasonList.Areas))
	log.Printf("Countries : %+v\n", len(seasonList.Countries))
	log.Printf("Competitions : %+v\n", len(seasonList.Competitions))
	log.Printf("Updated_at : %+v\n", seasonList.Updated_at)

	// areaInfoUpdate(seasonList.Areas)

	// countryInfoUpdate(seasonList.Countries)

	seasonListInfoUpdate(seasonList.Competitions)
}

func seasonListInfoUpdate(seasonList []data.NamiSeasonListInfo) {

	for _, seasonListInfo := range seasonList {

		log.Printf("seasonInfoUpdate seasonListInfo : %+v\n", seasonListInfo.Name_zh)
		for _, seasonCore := range seasonListInfo.Seasons {

			log.Printf("seasonInfoUpdate seasonCore year: %+v\n", seasonCore.Season)

			var seasonInfo data.NamiSeasonInfo
			data := make(map[string]string)
			data["id"] = strconv.FormatInt(int64(seasonCore.Id), 10)

			resp, err := commonFunc.GetNamiApi("/season/detail", data)
			if err != nil {
				log.Printf("seasonInfoUpdate err : %+v\n", err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("ioutil.ReadAll err : %+v\n", err)
				continue
			}
			// log.Printf("body : %+v\n", string(body))

			if err := json.Unmarshal(body, &seasonInfo); err != nil {
				log.Printf("gameList json.Unmarshal err : %+v\n", err)
				continue
			}

			log.Printf("Stages : %+v\n", len(seasonInfo.Stages))
			log.Printf("Teams : %+v\n", len(seasonInfo.Teams))
			log.Printf("Matches : %+v\n", len(seasonInfo.Matches))

			for _, stageInfo := range seasonInfo.Stages {
				stageData.SetStageData(stageInfo.Id, stageInfo)
			}

			leagueUuid, err := leagueData.SearchLeagueUuidByNami(seasonInfo.Competition.Id)

			if err != nil {
				log.Printf("SearchLeagueUuidByNami leagueUuid err : %+v\n", err)
				continue
			}

			seasonUuid, err := seasonData.SearchSeasonUuidByNami(seasonCore.Id)

			if err != nil {

				if err != database.ErrNoRows {
					log.Printf("seasonInfoUpdate select err : %+v\n", err)
					continue
				}

				seasonUuid = commonFunc.GetUuid()
				_, err := database.Exec("INSERT INTO `season` (`uuid` ,`league_uuid`,`season_year`) VALUES (?, ?, ?)",
					seasonUuid, leagueUuid, seasonCore.Season,
				)
				if err != nil {
					log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
					continue
				}

				corporationUuid := commonFunc.GetUuid()
				_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
					corporationUuid, "season", seasonUuid, "nami", seasonCore.Id,
				)

				if err != nil {
					log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
					continue
				}

				if err != nil {
					log.Printf("seasonInfoUpdate INSERT DB err : %+v\n", err)
					continue
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
						continue
					}
				} else {
					// go teamInfoUpdate(teamInfo.Id, seasonInfo.Competition.Id)
				}

				seasonTeamUuid := commonFunc.GetUuid()
				_, err = database.Exec("INSERT INTO `season_team_list` (`uuid` ,`seasonUuid` ,`teamUuid` ) VALUES (?, ?, ?)",
					seasonTeamUuid, seasonUuid, teamUuid,
				)
			}

			for _, gameInfo := range seasonInfo.Matches {
				_, err := gameData.SearchGameUuidByNami(gameInfo.Id)
				if err != nil {
					log.Printf("gameInfo : %+v\n", gameInfo.Id)
					gameInfoUpdate(gameInfo.Id, seasonUuid, gameInfo.Round.Stage_id, seasonInfo.Competition.Short_name_zh)
				} else {
					// go gameInfoUpdate(gameInfo.Id, seasonUuid, gameInfo.Round.Stage_id, seasonInfo.Competition.Short_name_zh)
				}
			}

		}
	}
}
