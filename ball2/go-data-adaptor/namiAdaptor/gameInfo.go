package namiAdaptor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"../commonFunc"
	"../data"
	"../data/gameData"
	"../data/leagueData"
	"../data/playerData"
	"../data/stageData"
	"../data/teamData"
	"../external/database"
)

func gameInfoUpdate(namiGameId int, seasonUuid string, stageId int, leagueShortName string) {

	stageInfo, ok := stageData.GetStageData(stageId)
	if !ok {
		// log.Printf("gameInfoUpdate GetStageData stageId : %+v\n", stageId)
		stageInfo = data.NamiStage{}
	}
	var subtitle string
	if stageInfo.Group_count != 0 && stageInfo.Round_count != 0 {
		subtitle = leagueShortName + stageInfo.Name_zh + strconv.Itoa(stageInfo.Round_count) + "輪" + "-" + data.GroupNum[stageInfo.Group_count-1] + "組"
	} else if stageInfo.Group_count != 0 {
		subtitle = leagueShortName + stageInfo.Name_zh + "-" + data.GroupNum[stageInfo.Group_count-1] + "組"
	} else if stageInfo.Round_count != 0 {
		subtitle = leagueShortName + stageInfo.Name_zh + strconv.Itoa(stageInfo.Round_count) + "輪"
	} else {
		subtitle = leagueShortName
	}

	var gameInfo data.NamiGameDetail

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(namiGameId), 10)

	resp, err := commonFunc.GetNamiApi("/match/detail", data)
	if err != nil {
		log.Printf("gameInfoUpdate err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &gameInfo); err != nil {
		log.Printf("gameInfoUpdate json.Unmarshal err : %+v\n", err)
		return
	}
	// log.Printf("gameInfo info : %+v\n", gameInfo.Info)
	// log.Printf("gameInfo matchevent : %+v\n", gameInfo.Matchevent)
	// log.Printf("gameInfo home_team : %+v\n", gameInfo.Home_team)
	// log.Printf("gameInfo away_team : %+v\n", gameInfo.Away_team)
	// log.Printf("gameInfo tlive : %+v\n", gameInfo.Tlive)
	// log.Printf("gameInfo stats : %+v\n", gameInfo.Stats)

	// log.Printf("areaInfoUpdate areaInfo : %+v\n", areaInfo)

	homeTeamUuid, err := teamData.SearchTeamUuidByNami(gameInfo.Home_team.Id)

	if err != nil {
		log.Printf("SearchTeamUuidByNami home team err : %+v\n", err)
		return
	}

	awayTeamUuid, err := teamData.SearchTeamUuidByNami(gameInfo.Away_team.Id)

	if err != nil {
		log.Printf("SearchTeamUuidByNami away team err : %+v\n", err)
		return
	}

	leagueUuid, err := leagueData.SearchLeagueUuidByNami(gameInfo.Matchevent.Id)

	if err != nil {
		log.Printf("SearchLeagueUuidByNami leagueUuid err : %+v\n", err)
		return
	}

	gameUuid, err := gameData.SearchGameUuidByNami(namiGameId)

	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("gameInfoUpdate select err : %+v\n", err)
			return
		}

		gameUuid, err = gameInsertSchedule(namiGameId, homeTeamUuid, awayTeamUuid, gameInfo, leagueUuid, seasonUuid, subtitle)

		if err != nil {
			log.Printf("gameInfoUpdate gameInsertSchedule err : %+v\n", err)
			return
		}
	} else {

		err = gameUpdateSchedule(gameUuid, homeTeamUuid, awayTeamUuid, gameInfo, leagueUuid, seasonUuid, subtitle)

		if err != nil {
			log.Printf("gameInfoUpdate gameUpdateSchedule err : %+v\n", err)
			return
		}
	}

	gameLineupUpdate(int(namiGameId))

	// log.Printf("gameInfoUpdate id : %+v\n", int64(namiGameId))
	// log.Printf("gameInfoUpdate scheduleUuid : %+v\n", gameUuid)

	return
}

func gameInsertSchedule(namiGameId int, homeTeamUuid string, awayTeamUuid string, gameInfo data.NamiGameDetail, leagueUuid string, seasonUuid string, subtitle string) (string, error) {

	gameUuid := commonFunc.GetUuid()
	_, err := database.Exec("INSERT INTO `schedule` (`uuid` ,`home_team`,`away_team`,`start_time`,`league_uuid`,`season_uuid`,`subtitle`,`state`) VALUES (?, ?,?, ?,?, ?,?, ?)",
		gameUuid, homeTeamUuid, awayTeamUuid, gameInfo.Info.Matchtime*1000, leagueUuid, seasonUuid, subtitle, gameInfo.Info.Statusid,
	)
	if err != nil {
		log.Printf("gameInfoUpdate INSERT DB err : %+v\n", err)
		return "", err
	}

	corporationUuid := commonFunc.GetUuid()
	_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
		corporationUuid, "game", gameUuid, "nami", namiGameId,
	)

	if err != nil {
		log.Printf("gameInfoUpdate INSERT DB err : %+v\n", err)
		return "", err
	}
	return gameUuid, nil

}

func gameUpdateSchedule(gameUuid string, homeTeamUuid string, awayTeamUuid string, gameInfo data.NamiGameDetail, leagueUuid string, seasonUuid string, subtitle string) error {

	_, err := database.Exec("UPDATE `schedule` SET "+
		"`home_team`=?,`away_team`=?,`start_time`=?,`league_uuid`=?,`season_uuid` = ? ,`subtitle` = ?,`state` = ?"+
		"  WHERE `uuid` = ?",
		homeTeamUuid, awayTeamUuid, gameInfo.Info.Matchtime*1000, leagueUuid, seasonUuid, subtitle, gameInfo.Info.Statusid, gameUuid)

	if err != nil {
		log.Printf("gameInfoUpdate UPDATE DB err : %+v\n", err)
		return err
	}
	return nil

}

func gameAnalysis(namiGameId int, gameUuid string) error {

	var gameAnalysis []interface{}
	var homeScoreArray [3]int
	var awayScoreArray [3]int

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(namiGameId), 10)

	resp, err := commonFunc.GetNamiApi("/match/analysis", data)
	if err != nil {
		log.Printf("gameInsertAnalysis err : %+v\n", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return err
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &gameAnalysis); err != nil {
		log.Printf("gameInsertAnalysis json.Unmarshal err : %+v\n", err)
		return err
	}

	homeScoreNami, ok := gameAnalysis[5].([9]float64)
	if !ok {
		log.Printf("gameAnalysis[5] not ok for type []interface{} : %+v\n", gameAnalysis[5])
		return nil
	}
	homeScoreArray[0] = int(homeScoreNami[2])
	homeScoreArray[1] = int(homeScoreNami[7])
	homeScoreArray[2] = int(homeScoreNami[8])

	awayScoreNami, ok := gameAnalysis[6].([9]float64)
	if !ok {
		log.Printf("gameAnalysis[6] not ok for type []interface{} : %+v\n", gameAnalysis[6])
		return nil
	}
	awayScoreArray[0] = int(awayScoreNami[2])
	awayScoreArray[1] = int(awayScoreNami[7])
	awayScoreArray[2] = int(awayScoreNami[8])

	homeScoreJson, _ := json.Marshal(homeScoreNami)
	awayScoreJson, _ := json.Marshal(awayScoreArray)

	ok = gameData.CheckGameAnaysisUuid(gameUuid)

	if !ok {
		gameInsertAnalysis(gameUuid, string(homeScoreJson), string(awayScoreJson))
	} else {
		gameUpdateAnalysis(gameUuid, string(homeScoreJson), string(awayScoreJson))
	}
	return nil
}

func gameInsertAnalysis(gameUuid string, homeScore string, awayScore string) {

	_, err := database.Exec("INSERT INTO `game_analysis` (`game_uuid` ,`home_score`,`away_score`) VALUES (?, ?, ?)",
		gameUuid, homeScore, awayScore,
	)
	if err != nil {
		log.Printf("gameInsertAnalysis INSERT DB err : %+v\n", err)
		return
	}

	return
}
func gameUpdateAnalysis(gameUuid string, homeScore string, awayScore string) {

	_, err := database.Exec("UPDATE `game_analysis` SET "+
		"`home_score`=?,`away_score`=?"+
		"  WHERE `game_uuid` = ?",
		homeScore, awayScore, gameUuid)

	if err != nil {
		log.Printf("gameUpdateAnalysis UPDATE DB err : %+v\n", err)
		return
	}
	return
}

func gameLineupUpdate(namiGameId int) {

	var gameLineup data.NamiGameLineup

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(namiGameId), 10)

	resp, err := commonFunc.GetNamiApi("/match/lineup", data)
	if err != nil {
		log.Printf("gameLineupUpdate err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &gameLineup); err != nil {
		log.Printf("body : %+v\n", string(body))
		log.Printf("gameLineupUpdate json.Unmarshal err : %+v\n", err)
		return
	}

	gameUuid, err := gameData.SearchGameUuidByNami(namiGameId)

	if err != nil {
		log.Printf("SearchGameUuidByNami gameUuid err : %+v\n", err)
		return
	}

	_, err = database.Exec("UPDATE `schedule` SET "+
		"`home_lineup`=?,`away_lineup`=?"+
		"  WHERE `uuid` = ?",
		gameLineup.Home_formation, gameLineup.Away_formation, gameUuid)

	if err != nil {
		log.Printf("gameLineupUpdate UPDATE err : %+v\n", err)
		return
	}

	playerLineupUpdate("home", gameUuid, gameLineup.Home)
	playerLineupUpdate("away", gameUuid, gameLineup.Away)

	// log.Printf("gameLineup schedule uuid: %+v\n", gameUuid)
}

func playerLineupUpdate(teamType string, gameUuid string, lineupList []data.NamiGamePlayerLineup) {

	for _, lineupInfo := range lineupList {

		var lineupUuid string
		row := database.QueryRow("SELECT `uuid` FROM `line_up` WHERE `player_uuid` = ? AND `schedule_uuid` = ? ",
			lineupInfo.Id, gameUuid,
		)
		err := row.Scan(&lineupUuid)
		if err != nil {
			if err != database.ErrNoRows {
				log.Printf("playerLineupUpdate select err : %+v\n", err)
				continue
			}

			playerUuid, err := playerData.SearchPlayerUuidByNami(lineupInfo.Id)

			if err != nil {
				log.Printf("playerLineupUpdate lineupInfo.Id err : %+v\n", lineupInfo.Id)
				log.Printf("playerLineupUpdate SearchPlayerUuidByNami err : %+v\n", err)
				playerUuid, err = playerInfoUpdate(lineupInfo.Id)
				if err != nil {
					log.Printf("playerLineupUpdate playerInfoUpdate err : %+v\n", err)
					continue
				}
			}

			lineupUuid = commonFunc.GetUuid()
			_, err = database.Exec("INSERT INTO `line_up` (`uuid` ,`player_uuid`,`schedule_uuid`,`type`,`first`,`shirt_number`,`position`,`x`,`y`) VALUES (?, ?,?, ?,?, ?,?, ?,?)",
				lineupUuid, playerUuid, gameUuid, teamType, lineupInfo.First, lineupInfo.Shirt_number, lineupInfo.Position, lineupInfo.X, lineupInfo.Y,
			)
			if err != nil {
				log.Printf("gameInfoUpdate INSERT DB lineupInfo.Position : %+v\n", lineupInfo.Position)
				log.Printf("gameInfoUpdate INSERT DB err : %+v\n", err)
				continue
			}
			continue
		}

		playerUuid, err := playerData.SearchPlayerUuidByNami(lineupInfo.Id)

		if err != nil {
			log.Printf("playerLineupUpdate lineupInfo.Id err : %+v\n", lineupInfo.Id)
			log.Printf("playerLineupUpdate SearchPlayerUuidByNami err : %+v\n", err)
			playerUuid, err = playerInfoUpdate(lineupInfo.Id)
			if err != nil {
				log.Printf("playerLineupUpdate playerInfoUpdate err : %+v\n", err)
				continue
			}
		}

		_, err = database.Exec("UPDATE `line_up` SET "+
			"`player_uuid`=?,`schedule_uuid`=?,`type`=?,`first`=?,`shirt_number`=?,`position`=?,`x`=?,`y`=?"+
			"  WHERE `uuid` = ?",
			playerUuid, gameUuid, teamType, lineupInfo.First, lineupInfo.Shirt_number, lineupInfo.Position, lineupInfo.X, lineupInfo.Y, lineupUuid)

		if err != nil {
			log.Printf("playerLineupUpdate UPDATE err : %+v\n", err)
			continue
		}

	}
}
