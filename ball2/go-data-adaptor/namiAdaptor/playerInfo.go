package namiAdaptor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"../commonFunc"
	"../data"
	"../data/countryData"
	"../data/playerData"
	"../data/teamData"
	"../external/database"
)

func playerInfoUpdate(namiPlayerId int) (string, error) {

	var playerInfo data.NamiPlayerInfo

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(namiPlayerId), 10)

	// log.Printf("playerInfo id : %+v\n", int64(namiPlayerId))

	resp, err := commonFunc.GetNamiApi("/player/detail", data)
	if err != nil {
		log.Printf("/player/detail err : %+v\n", err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return "", err
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &playerInfo); err != nil {
		log.Printf("body : %+v\n", string(body))
		log.Printf("playerInfo json.Unmarshal err : %+v\n", err)
		return "", err
	}

	playerUuid, err := playerData.SearchPlayerUuidByNami(namiPlayerId)
	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("player select err : %+v\n", err)
			return "", err
		}

		countryUuid, err := countryData.SearchCountryUuidByNami(playerInfo.Country_id)

		if err != nil {
			log.Printf("playerInfoUpdate update countryUuid namiId : %+v\n", playerInfo.Country_id)
			log.Printf("playerInfoUpdate update countryUuid err : %+v\n", err)
			return "", err
		}

		playerUuid = commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `player` "+
			"(`uuid` ,`name` ,`name_en` , `weight`, `height`, `birthday`, `country_uuid`, `preferred_foot`,`contract_until`,`market_value`,`market_value_currency`,`position`)"+
			" VALUES (?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?)",
			playerUuid, playerInfo.Name_zh, playerInfo.Name_en, playerInfo.Weight, playerInfo.Height, playerInfo.Birthday, countryUuid, playerInfo.Preferred_foot, playerInfo.Contract_until, playerInfo.Market_value, playerInfo.Market_value_currency, playerInfo.Position,
		)
		if err != nil {

			log.Printf("playerInfoUpdate insert DB teamInfo : %+v\n", playerInfo)
			log.Printf("playerInfoUpdate insert DB err : %+v\n", err)
			return "", err
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "player", playerUuid, "nami", namiPlayerId,
		)
		if err != nil {
			log.Printf("playerInfoUpdate namiPlayerId : %+v\n", namiPlayerId)
			log.Printf("playerInfoUpdate INSERT DB err : %+v\n", err)
			return "", err
		}

		if playerInfo.Logo != "" {
			_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/player/"+playerInfo.Logo, playerUuid+".png", "/playerIcon/")
			if err != nil {
				log.Printf("PostCdnUploadLink err : %+v\n", err)
			}
		}

		teamUuid, err := teamData.SearchTeamUuidByNami(playerInfo.Team_id)
		if err != nil {
			log.Printf("playerInfoUpdate playerInfo.Team_id : %+v\n", playerInfo.Team_id)
			log.Printf("playerInfoUpdate SearchTeamUuidByNami err : %+v\n", err)
			return "", err
		}
		teamMemberListUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `team_member_list` (`uuid` ,`team_uuid` ,`player_uuid`) VALUES (?, ?, ?)",
			teamMemberListUuid, teamUuid, playerUuid,
		)

		if err != nil {
			log.Printf("playerInfoUpdate insert team_member_list err : %+v\n", err)
			return "", err
		}

		// log.Printf("playerInfo id : %+v\n", int64(namiPlayerId))
		// log.Printf("playerInfo playerUuid : %+v\n", playerUuid)
		return playerUuid, nil

	}

	countryUuid, err := countryData.SearchCountryUuidByNami(playerInfo.Country_id)

	if err != nil {
		log.Printf("playerInfoUpdate update countryUuid namiId : %+v\n", playerInfo.Country_id)
		log.Printf("playerInfoUpdate update countryUuid err : %+v\n", err)
		return "", err
	}

	_, err = database.Exec("UPDATE `player` SET "+
		"`name` = ? ,`name_en` = ? , `weight` = ? , `height` = ? , `birthday` = ? , `country_uuid` = ?, `preferred_foot` = ? , contract_until = ? , market_value = ? , market_value_currency = ? , position = ?"+
		" WHERE `uuid` = ?",
		playerInfo.Name_zh, playerInfo.Name_en, playerInfo.Weight, playerInfo.Height, playerInfo.Birthday, countryUuid, playerInfo.Preferred_foot, playerInfo.Contract_until, playerInfo.Market_value, playerInfo.Market_value_currency, playerInfo.Position,
		playerUuid,
	)

	if err != nil {
		log.Printf("playerInfoUpdate update player err : %+v\n", err)
		return "", err
	}

	if playerInfo.Logo != "" {
		_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/player/"+playerInfo.Logo, playerUuid+".png", "/playerIcon/")
		if err != nil {
			log.Printf("PostCdnUploadLink  err : %+v\n", err)
		}
	}

	teamUuid, err := teamData.SearchTeamUuidByNami(playerInfo.Team_id)
	if err != nil {
		log.Printf("playerInfoUpdate playerInfo.Team_id : %+v\n", playerInfo.Team_id)
		log.Printf("playerInfoUpdate SearchTeamUuidByNami err : %+v\n", err)
		return "", err
	}

	_, err = database.Exec("UPDATE `team_member_list` SET `team_uuid` = ?  WHERE `player_uuid` = ?",
		teamUuid, playerUuid)

	if err != nil {
		log.Printf("playerInfoUpdate update team_member_list err : %+v\n", err)
		return "", err
	}

	// log.Printf("playerInfo id : %+v\n", int64(namiPlayerId))
	// log.Printf("playerInfo playerUuid : %+v\n", playerUuid)

	return playerUuid, nil
}
