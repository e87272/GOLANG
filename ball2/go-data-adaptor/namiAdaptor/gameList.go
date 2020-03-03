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
)

func gameList() {

	fmt.Println(" gameList...")

	for {
		now := time.Now()
		// 计算下一天
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 13, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		gameListUpdate()
		gameList()
	}

}

func gameListUpdate() {
	for i := -1; i < 8; i++ {

		var gameList data.NamiGameList
		t := time.Now().AddDate(0, 0, i)
		year, month, day := t.Date()
		yearStr := strconv.FormatInt(int64(year), 10)
		monthStr := strconv.FormatInt(int64(month+100), 10)[1:3]
		dayStr := strconv.FormatInt(int64(day+100), 10)[1:3]
		dateStr := yearStr + monthStr + dayStr
		data := make(map[string]string)
		data["date"] = dateStr
		log.Printf("gameListTicker dateStr : %+v\n", dateStr)

		resp, err := commonFunc.GetNamiApi("/match/list", data)
		if err != nil {
			log.Printf("gameListTicker err : %+v\n", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err : %+v\n", err)
			return
		}
		// log.Printf("body : %+v\n", string(body))

		if err := json.Unmarshal(body, &gameList); err != nil {
			log.Printf("gameList json.Unmarshal err : %+v\n", err)
			return
		}
		// gameData.SetNamiGameListData(dateStr, gameList)
		log.Printf("Teams : %+v\n", len(gameList.Teams))
		log.Printf("Events : %+v\n", len(gameList.Events))
		log.Printf("Matches : %+v\n", len(gameList.Matches))
		log.Printf("Stages : %+v\n", len(gameList.Stages))

		for _, stageInfo := range gameList.Stages {
			stageData.SetStageData(stageInfo.Id, stageInfo)
		}

		for _, leagueInfo := range gameList.Events {
			leagueData.SetLeagueData(leagueInfo.Id, leagueInfo)
		}

		for _, teamInfo := range gameList.Teams {
			_, err := teamData.SearchTeamUuidByNami(teamInfo.Id)
			if err != nil {
				log.Printf("teamInfo : %+v\n", teamInfo.Id)
				teamInfoUpdate(teamInfo.Id, teamInfo.Matchevent_id)
			} else {
				//已有資料背景同步
				// go teamInfoUpdate(teamInfo.Id, teamInfo.Matchevent_id)
			}
		}

		for _, gameInfo := range gameList.Matches {
			// log.Printf("gameInfo : %+v\n", gameInfo)
			// log.Printf("gameInfo : %+v\n", gameInfo[0])
			namiGameId, ok := gameInfo[0].(float64)
			if !ok {
				log.Printf("gameInfo[0] not ok for type float64 : %+v\n", gameInfo[0])
				continue
			}

			namiLeagueId, ok := gameInfo[1].(float64)
			if !ok {
				log.Printf("gameInfo[1] not ok for type float64 : %+v\n", gameInfo[1])
				continue
			}
			leagueInfo, ok := leagueData.GetLeagueData(int(namiLeagueId))
			if !ok {
				log.Printf("GetLeagueData namiLeagueId err : %+v\n", namiLeagueId)
				continue
			}

			namiSeasonList, ok := gameInfo[8].([]interface{})
			if !ok {
				log.Printf("gameInfo[8] not ok for type []interface{} : %+v\n", gameInfo[8])
				continue
			}
			// log.Printf("gameInfo[8] not ok for type []interface{} : %+v\n", gameInfo[8])
			namiSeasonId, ok := namiSeasonList[0].(float64)
			if !ok {
				log.Printf("namiSeasonList[0] not ok for type float64 : %+v\n", namiSeasonList[0])
				continue
			}
			namiSeasonYear, ok := namiSeasonList[1].(string)
			if !ok {
				log.Printf("namiSeasonList[1] not ok for type string : %+v\n", namiSeasonList[1])
				continue
			}

			namiStageList, ok := gameInfo[9].([]interface{})
			if !ok {
				log.Printf("gameInfo[9] not ok for type []interface{}  : %+v\n", gameInfo[9])
				continue
			}
			namiStageId, ok := namiStageList[0].(float64)
			if !ok {
				log.Printf("namiStageList[8] not ok for type float64 : %+v\n", namiStageList[0])
				continue
			}

			seasonUuid, err := seasonData.SearchSeasonUuidByNami(int(namiSeasonId))
			if err != nil {
				log.Printf("namiSeasonId : %+v\n", int(namiSeasonId))
				seasonUuid, err = seasonInfoUpdate(int(namiSeasonId), namiSeasonYear)
				if err != nil {
					log.Printf("seasonInfoUpdate err : %+v\n", err)
					continue
				}
			} else {
				go seasonInfoUpdate(int(namiSeasonId), namiSeasonYear)
			}

			_, err = gameData.SearchGameUuidByNami(int(namiGameId))
			// log.Printf("namiGameId : %+v\n", int(namiGameId))
			if err != nil {
				gameInfoUpdate(int(namiGameId), seasonUuid, int(namiStageId), leagueInfo.Short_name_zh)
			} else {
				go gameInfoUpdate(int(namiGameId), seasonUuid, int(namiStageId), leagueInfo.Short_name_zh)
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
