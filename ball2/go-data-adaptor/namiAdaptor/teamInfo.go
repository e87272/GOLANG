package namiAdaptor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"../commonFunc"
	"../data"
	"../data/countryData"
	"../data/leagueData"
	"../data/playerData"
	"../data/teamData"
	"../external/database"
)

func teamInfoUpdate(namiTeamId int, namiLeagueId int) {

	var teamInfo data.NamiTeamInfo

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(namiTeamId), 10)

	// log.Printf("teamInfo id : %+v\n", int64(namiTeamId))

	resp, err := commonFunc.GetNamiApi("/team/detail", data)
	if err != nil {
		log.Printf("teamListTicker err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &teamInfo); err != nil {
		log.Printf("body : %+v\n", string(body))
		log.Printf("teamInfo json.Unmarshal err : %+v\n", err)
		return
	}

	teamUuid, err := teamData.SearchTeamUuidByNami(namiTeamId)

	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("teamList select err : %+v\n", err)
			return
		}

		teamUuid = commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `team` (`uuid` ,`name` , `name_en` ,`manager` , `venue`) VALUES (?, ?, ?, ?, ?)",
			teamUuid, teamInfo.Name_zh, teamInfo.Name_en, teamInfo.Manager.Name_zh, teamInfo.Venue.Name_zh,
		)
		if err != nil {

			log.Printf("teamInfoUpdate insert DB teamInfo : %+v\n", teamInfo)
			log.Printf("teamInfoUpdate insert DB err : %+v\n", err)
			return
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "team", teamUuid, "nami", namiTeamId,
		)
		if err != nil {
			log.Printf("teamInfoUpdate update DB err : %+v\n", err)
			return
		}

		leagueUuid, err := leagueData.SearchLeagueUuidByNami(namiLeagueId)

		if err != nil {
			log.Printf("namiLeagueId : %+v\n", namiLeagueId)
			log.Printf("teamInfoUpdate leagueUuid select corporation_id err : %+v\n", err)
			return
		}

		leagueTeamListUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `league_team_list` (`uuid` ,`league_uuid` ,`team_uuid`) VALUES (?, ?, ?)",
			leagueTeamListUuid, leagueUuid, teamUuid,
		)

		if err != nil {
			log.Printf("teamInfoUpdate insert league_team_list err : %+v\n", err)
			return
		}

		if teamInfo.Logo != "" {
			_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/team/"+teamInfo.Logo, teamUuid+".png", "/teamIcon/")
			if err != nil {
				log.Printf("PostCdnUploadLink err : %+v\n", err)
			}
		}

		for _, playerInfo := range teamInfo.Players {
			// log.Printf("playerInfo : %+v\n", playerInfo.Player.Id)
			teamPlayerInfoUpdate(teamUuid, playerInfo)
		}

		return

	}

	if teamInfo.Logo != "" {
		_, err := commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/team/"+teamInfo.Logo, teamUuid+".png", "/teamIcon/")
		if err != nil {
			log.Printf("PostCdnUploadLink err : %+v\n", err)
		}
	}
	_, err = database.Exec("UPDATE `team` SET `name` = ? ,`name_en` = ? ,`manager` = ? , `venue` = ?   WHERE `uuid` = ?",
		teamInfo.Name_zh, teamInfo.Name_en, teamInfo.Manager.Name_zh, teamInfo.Venue.Name_zh, teamUuid)

	if err != nil {
		log.Printf("teamInfoUpdate update team err : %+v\n", err)
		return
	}

	leagueUuid, err := leagueData.SearchLeagueUuidByNami(namiLeagueId)

	if err != nil {
		log.Printf("namiLeagueId : %+v\n", namiLeagueId)
		log.Printf("teamInfoUpdate leagueUuid select corporation_id err : %+v\n", err)
		return
	}

	_, err = database.Exec("UPDATE `league_team_list` SET `league_uuid` = ?  WHERE `team_uuid` = ?",
		leagueUuid, teamUuid)

	if err != nil {
		log.Printf("teamInfoUpdate update DB err : %+v\n", err)
		return
	}

	for _, playerInfo := range teamInfo.Players {
		// log.Printf("playerInfo : %+v\n", playerInfo.Player.Id)
		teamPlayerInfoUpdate(teamUuid, playerInfo)
	}
}

var mutexPlayer = new(sync.Mutex)

func teamPlayerInfoUpdate(teamUuid string, playerInfo data.NamiTeamPlayerInfo) {

	//預防多對相同球員出錯
	mutexPlayer.Lock()
	defer mutexPlayer.Unlock()

	playerUuid, err := playerData.SearchPlayerUuidByNami(playerInfo.Player.Id)
	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("player select err : %+v\n", err)
			return
		}

		countryUuid, err := countryData.SearchCountryUuidByNami(playerInfo.Player.Country_id)

		if err != nil {
			log.Printf("playerInfoUpdate update countryUuid namiId : %+v\n", playerInfo.Player.Country_id)
			log.Printf("playerInfoUpdate update countryUuid err : %+v\n", err)
			return
		}

		playerUuid = commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `player` "+
			"(`uuid` ,`name` ,`name_en` , `weight`, `height`, `birthday`, `country_uuid`, `preferred_foot`,`contract_until`,`market_value`,`market_value_currency`,`shirt_number`,`position`)"+
			" VALUES (?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?)",
			playerUuid, playerInfo.Player.Name_zh, playerInfo.Player.Name_en, playerInfo.Player.Weight, playerInfo.Player.Height, playerInfo.Player.Birthday, countryUuid, playerInfo.Player.Preferred_foot, playerInfo.Player.Contract_until, playerInfo.Player.Market_value, playerInfo.Player.Market_value_currency, playerInfo.Shirt_number, playerInfo.Position,
		)
		if err != nil {

			log.Printf("playerInfoUpdate insert DB teamInfo : %+v\n", playerInfo)
			log.Printf("playerInfoUpdate insert DB err : %+v\n", err)
			return
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "player", playerUuid, "nami", playerInfo.Player.Id,
		)
		if err != nil {
			log.Printf("playerInfoUpdate playerInfo.Player.Id : %+v\n", playerInfo.Player.Id)
			log.Printf("playerInfoUpdate INSERT DB err : %+v\n", err)
			return
		}

		teamMemberListUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `team_member_list` (`uuid` ,`team_uuid` ,`player_uuid`) VALUES (?, ?, ?)",
			teamMemberListUuid, teamUuid, playerUuid,
		)

		if err != nil {
			log.Printf("playerInfoUpdate insert team_member_list err : %+v\n", err)
			return
		}

		if playerInfo.Player.Logo != "" {
			_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/player/"+playerInfo.Player.Logo, playerUuid+".png", "/playerIcon/")
			if err != nil {
				log.Printf("PostCdnUploadLink err : %+v\n", err)
			}
		}

		return

	}

	countryUuid, err := countryData.SearchCountryUuidByNami(playerInfo.Player.Country_id)

	if err != nil {
		log.Printf("playerInfoUpdate update countryUuid namiId : %+v\n", playerInfo.Player.Country_id)
		log.Printf("playerInfoUpdate update countryUuid err : %+v\n", err)
		return
	}

	if playerInfo.Player.Logo != "" {
		_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/player/"+playerInfo.Player.Logo, playerUuid+".png", "/playerIcon/")
		if err != nil {
			log.Printf("PostCdnUploadLink  err : %+v\n", err)
		}
	}

	_, err = database.Exec("UPDATE `player` SET "+
		"`name` = ? ,`name_en` = ? , `weight` = ? , `height` = ? , `birthday` = ? , `country_uuid` = ?, `preferred_foot` = ? , contract_until = ? , market_value = ? , market_value_currency = ? , shirt_number = ? , position = ?"+
		" WHERE `uuid` = ?",
		playerInfo.Player.Name_zh, playerInfo.Player.Name_en, playerInfo.Player.Weight, playerInfo.Player.Height, playerInfo.Player.Birthday, countryUuid, playerInfo.Player.Preferred_foot, playerInfo.Player.Contract_until, playerInfo.Player.Market_value, playerInfo.Player.Market_value_currency, playerInfo.Shirt_number, playerInfo.Position,
		playerUuid,
	)

	if err != nil {
		log.Printf("playerInfoUpdate update player err : %+v\n", err)
		return
	}

	_, err = database.Exec("UPDATE `team_member_list` SET `team_uuid` = ?  WHERE `player_uuid` = ?",
		teamUuid, playerUuid)

	if err != nil {
		log.Printf("playerInfoUpdate update team_member_list err : %+v\n", err)
	}

}
