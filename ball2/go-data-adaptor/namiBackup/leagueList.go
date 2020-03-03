package namiBackup

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"../commonFunc"
	"../data"
	"../external/namiDatabase"
)

func backupLeagueList() {

	log.Printf("backupLeagueList...\n")

	var namiLeagueList data.NamiLeagueList

	resp, err := commonFunc.GetNamiApi("/matchevent/list", nil)
	if err != nil {
		log.Printf("backupLeagueList err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &namiLeagueList); err != nil {
		log.Printf("backupLeagueList json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("Areas : %+v\n", len(namiLeagueList.Areas))
	go backupLeagueAreaList(namiLeagueList.Areas)

	log.Printf("Countrys : %+v\n", len(namiLeagueList.Countrys))
	go backupLeagueCountryList(namiLeagueList.Countrys)

	log.Printf("LeagueList : %+v\n", len(namiLeagueList.Matchevents))
	go backupLeagueInfoList(namiLeagueList.Matchevents)
}

func backupLeagueAreaList(areaList map[string]data.NamiAreaInfo) {

	tx, err := namiDatabase.Begin()
	if err != nil {
		log.Printf("backupLeagueAreaList Begin err : %+v\n", err)
		return
	}
	stmt, err := namiDatabase.Prepare(tx, "INSERT INTO area SET id=? , name_zh=? , name_zht=? , name_en=? "+
		"ON DUPLICATE KEY UPDATE name_zh=? , name_zht=? , name_en=? ")
	if err != nil {
		log.Printf("backupLeagueAreaList Prepare err : %+v\n", err)
		return
	}
	for _, areaInfo := range areaList {

		_, err = stmt.Exec(areaInfo.Id, areaInfo.Name_zh, areaInfo.Name_zht, areaInfo.Name_en, areaInfo.Name_zh, areaInfo.Name_zht, areaInfo.Name_en)
		if err != nil {
			log.Printf("backupLeagueAreaList UPDATE DB err : %+v\n", err)
			continue
		}
	}

	namiDatabase.Commit(tx)
	return
}

func backupLeagueCountryList(countryList map[string]data.NamiCountryInfo) {

	tx, err := namiDatabase.Begin()
	if err != nil {
		log.Printf("backupSeasonCountryList Begin err : %+v\n", err)
		return
	}
	stmt, err := namiDatabase.Prepare(tx, "INSERT INTO country SET id = ? , area_id = ? , name_zh = ? , name_zht = ? , name_en = ? , logo = ? "+
		"ON DUPLICATE KEY UPDATE area_id = ? , name_zh = ? , name_zht = ? , name_en = ? , logo = ? ")
	if err != nil {
		log.Printf("backupSeasonCountryList Prepare err : %+v\n", err)
		return
	}

	for _, countryInfo := range countryList {

		_, err = stmt.Exec(countryInfo.Id, countryInfo.Area_id, countryInfo.Name_zh, countryInfo.Name_zht, countryInfo.Name_en, countryInfo.Logo, countryInfo.Area_id, countryInfo.Name_zh, countryInfo.Name_zht, countryInfo.Name_en, countryInfo.Logo)
		if err != nil {
			log.Printf("backupLeagueCountryList UPDATE DB err : %+v\n", err)
			continue
		}
	}
	namiDatabase.Commit(tx)
	return
}

func backupLeagueInfoList(leagueList []data.NamiLeagueInfo) {

	tx, err := namiDatabase.Begin()
	if err != nil {
		log.Printf("backupLeagueInfoList Begin err : %+v\n", err)
		return
	}
	stmt, err := namiDatabase.Prepare(tx, "INSERT INTO league SET id = ? , area_id = ? , country_id = ? , type = ? , level = ? , name_zh = ?  , short_name_zh = ? , name_zht = ?  , short_name_zht = ? , name_en = ? , short_name_en = ? , logo = ? , seasonList = ? "+
		"ON DUPLICATE KEY UPDATE area_id = ? , country_id = ? , type = ? , level = ? , name_zh = ?  , short_name_zh = ? , name_zht = ?  , short_name_zht = ? , name_en = ? , short_name_en = ? , logo = ? ")
	if err != nil {
		log.Printf("backupLeagueInfoList Prepare err : %+v\n", err)
		return
	}

	for _, leagueInfo := range leagueList {

		_, err = stmt.Exec(leagueInfo.Id, leagueInfo.Area_id, leagueInfo.Country_id, leagueInfo.Type, leagueInfo.Level, leagueInfo.Name_zh, leagueInfo.Short_name_zh, leagueInfo.Name_zht, leagueInfo.Short_name_zht, leagueInfo.Name_en, leagueInfo.Short_name_en, leagueInfo.Logo, "{}",
			leagueInfo.Area_id, leagueInfo.Country_id, leagueInfo.Type, leagueInfo.Level, leagueInfo.Name_zh, leagueInfo.Short_name_zh, leagueInfo.Name_zht, leagueInfo.Short_name_zht, leagueInfo.Name_en, leagueInfo.Short_name_en, leagueInfo.Logo)
		if err != nil {

			log.Printf("backupLeagueInfoList UPDATE DB err : %+v\n", err)
			continue
		}
	}
	namiDatabase.Commit(tx)
	return
}
