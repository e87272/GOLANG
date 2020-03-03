package namiBackup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"../commonFunc"
	"../data"
	"../external/namiDatabase"
)

func backupGameDate() {

	log.Printf("backupGameDate...\n")

	for i := -1; i < 8; i++ {

		var gameList data.NamiGameList
		t := time.Now().AddDate(0, 0, i)
		year, month, day := t.Date()
		yearStr := strconv.FormatInt(int64(year), 10)
		monthStr := strconv.FormatInt(int64(month+100), 10)[1:3]
		dayStr := strconv.FormatInt(int64(day+100), 10)[1:3]
		dateStr := yearStr + monthStr + dayStr
		apiData := make(map[string]string)
		apiData["date"] = dateStr
		log.Printf("backupGameDate dateStr : %+v\n", dateStr)

		resp, err := commonFunc.GetNamiApi("/match/list", apiData)
		if err != nil {
			log.Printf("backupGameDate err : %+v\n", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err : %+v\n", err)
			return
		}
		// log.Printf("body : %+v\n", string(body))

		if err := json.Unmarshal(body, &gameList); err != nil {
			log.Printf("backupGameDate json.Unmarshal err : %+v\n", err)
			return
		}
		// gameData.SetNamiGameListData(dateStr, gameList)
		log.Printf("Teams : %+v\n", len(gameList.Teams))
		log.Printf("Events : %+v\n", len(gameList.Events))
		log.Printf("Matches : %+v\n", len(gameList.Matches))
		log.Printf("Stages : %+v\n", len(gameList.Stages))

		tx, err := namiDatabase.Begin()
		if err != nil {
			log.Printf("backupGameDate stageInfo Begin err : %+v\n", err)
			return
		}
		stmt, err := namiDatabase.Prepare(tx, "INSERT INTO stages SET id=? ,mode=? ,group_count=? ,round_count=? , name_zh=? , name_zht=? , name_en=? "+
			"ON DUPLICATE KEY UPDATE mode=? ,group_count=? ,round_count=? , name_zh=? , name_zht=? , name_en=? ")
		if err != nil {
			log.Printf("backupGameDate stageInfo Prepare err : %+v\n", err)
			return
		}
		for _, stageInfo := range gameList.Stages {

			_, err = stmt.Exec(stageInfo.Id, stageInfo.Mode, stageInfo.Group_count, stageInfo.Round_count, stageInfo.Name_zh, stageInfo.Name_zht, stageInfo.Name_en,
				stageInfo.Mode, stageInfo.Group_count, stageInfo.Round_count, stageInfo.Name_zh, stageInfo.Name_zht, stageInfo.Name_en)
			if err != nil {
				log.Printf("backupGameDate stageInfo UPDATE DB err : %+v\n", err)
				continue
			}
		}
		namiDatabase.Commit(tx)

		tx, err = namiDatabase.Begin()
		if err != nil {
			log.Printf("backupGameDate stageInfo Begin err : %+v\n", err)
			return
		}
		stmt, err = namiDatabase.Prepare(tx, "INSERT INTO game SET id=? , home_team=? , away_team=? , tlive=? , stats=? , info=? , history=? , goal_distribution=? , injury=? , league_table=? "+
			"ON DUPLICATE KEY UPDATE home_team=? , away_team=? , tlive=? , stats=? , info=? , history=? , goal_distribution=? , injury=? , league_table=?")
		if err != nil {
			log.Printf("backupGameDate stageInfo Prepare err : %+v\n", err)
			return
		}
		for _, gameInfo := range gameList.Matches {
			// log.Printf("gameInfo : %+v\n", gameInfo)
			// log.Printf("gameInfo : %+v\n", gameInfo[0])
			namiGameId, ok := gameInfo[0].(float64)
			if !ok {
				log.Printf("gameInfo[0] not ok for type float64 : %+v\n", gameInfo[0])
				continue
			}

			var gameDetail data.NamiGameDetail

			apiAata := make(map[string]string)
			apiAata["id"] = strconv.FormatInt(int64(namiGameId), 10)

			resp, err := commonFunc.GetNamiApi("/match/detail", apiAata)
			if err != nil {
				log.Printf("backupGameDate err : %+v\n", err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("ioutil.ReadAll err : %+v\n", err)
				return
			}
			// log.Printf("body : %+v\n", string(body))

			if err := json.Unmarshal(body, &gameDetail); err != nil {
				log.Printf("backupGameDate json.Unmarshal err : %+v\n", err)
				return
			}

			var gameAnaysis data.NamiGameAnaysis

			apiAata = make(map[string]string)
			apiAata["id"] = strconv.FormatInt(int64(namiGameId), 10)

			resp, err = commonFunc.GetNamiApi("/match/analysis", apiAata)
			if err != nil {
				log.Printf("backupGameDate err : %+v\n", err)
				return
			}

			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("ioutil.ReadAll err : %+v\n", err)
				return
			}
			// log.Printf("body : %+v\n", string(body))

			if err := json.Unmarshal(body, &gameAnaysis); err != nil {
				log.Printf("backupGameDate json.Unmarshal err : %+v\n", err)
				return
			}

			homeTeamJson, _ := json.Marshal(gameDetail.Home_team)
			awayTeamJson, _ := json.Marshal(gameDetail.Away_team)
			tliveJson, _ := json.Marshal(gameDetail.Tlive)
			statsJson, _ := json.Marshal(gameDetail.Stats)
			infoJson, _ := json.Marshal(gameAnaysis.Info)
			historyJson, _ := json.Marshal(gameAnaysis.History)
			goalDistributionJson, _ := json.Marshal(gameAnaysis.Goal_distribution)
			injuryJson, _ := json.Marshal(gameAnaysis.Injury)
			leagueTableJson, _ := json.Marshal(gameAnaysis.Table)

			_, err = stmt.Exec(namiGameId, string(homeTeamJson), string(awayTeamJson), string(tliveJson), string(statsJson), string(infoJson), string(historyJson), string(goalDistributionJson), string(injuryJson), string(leagueTableJson),
				string(homeTeamJson), string(awayTeamJson), string(tliveJson), string(statsJson), string(infoJson), string(historyJson), string(goalDistributionJson), string(injuryJson), string(leagueTableJson))
			if err != nil {
				log.Printf("backupGameDate stageInfo UPDATE DB err : %+v\n", err)
				continue
			}

		}
		namiDatabase.Commit(tx)

	}
}
