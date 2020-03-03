package namiBackup

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"../commonFunc"
	"../data"
	"../external/namiDatabase"
)

func backupTeamList() {

	log.Printf("backupTeamList...\n")

	var namiTeamList data.NamiTeamList

	resp, err := commonFunc.GetNamiApi("/team/list", nil)
	if err != nil {
		log.Printf("backupTeamList err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &namiTeamList); err != nil {
		log.Printf("backupTeamList json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("namiTeamList : %+v\n", len(namiTeamList))
	for _, teamInfo := range namiTeamList {

		result, err := namiDatabase.Exec("UPDATE `team` SET `matchevent_id` = ? , `name_zh` = ? , `name_zht` = ? , `name_en` = ? , `short_name_zh` = ?  , `short_name_zht` = ? , `short_name_en` = ?  , `logo` = ? , `found` = ? , `website` = ? , `national` = ? , `country_logo` = ?  WHERE `id` = ?",
			teamInfo.Matchevent_id, teamInfo.Name_zh, teamInfo.Name_zht, teamInfo.Name_en, teamInfo.Short_name_zh, teamInfo.Short_name_zht, teamInfo.Short_name_en, teamInfo.Logo, teamInfo.Found, teamInfo.Website, teamInfo.National, teamInfo.Country_logo, teamInfo.Id)

		if err != nil {
			log.Printf("backupteamInfoList UPDATE DB err : %+v\n", err)
			continue
		}

		idAff, err := result.RowsAffected()
		if err != nil {
			log.Println("backupteamInfoList RowsAffected failed:", err)
			continue
		}
		if idAff == 0 {
			_, err := namiDatabase.Exec("INSERT INTO `team` (`id` , `matchevent_id`, `name_zh`, `name_zht`, `name_en`, `short_name_zh` , `short_name_zht`, `short_name_en` , `logo`, `found`, `website`, `national`, `country_logo`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				teamInfo.Id, teamInfo.Matchevent_id, teamInfo.Name_zh, teamInfo.Name_zht, teamInfo.Name_en, teamInfo.Short_name_zh, teamInfo.Short_name_zht, teamInfo.Short_name_en, teamInfo.Logo, teamInfo.Found, teamInfo.Website, teamInfo.National, teamInfo.Country_logo)
			if err != nil {
				log.Printf("backupteamInfoList INSERT DB err : %+v\n", err)
				continue
			}
		}
	}
	return
}
