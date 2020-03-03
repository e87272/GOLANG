package namiBackup

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"../commonFunc"
	"../data"
	"../external/namiDatabase"
)

func backupSeasonList() {

	log.Printf("backupSeasonList...\n")

	var namiSeasonList data.NamiSeasonList

	resp, err := commonFunc.GetNamiApi("/season/list", nil)
	if err != nil {
		log.Printf("backupSeasonList err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &namiSeasonList); err != nil {
		log.Printf("backupSeasonList json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("Areas : %+v\n", len(namiSeasonList.Areas))
	backupSeasonAreaList(namiSeasonList.Areas)

	log.Printf("Countrys : %+v\n", len(namiSeasonList.Countries))
	backupSeasonCountryList(namiSeasonList.Countries)

	log.Printf("SeasonList : %+v\n", len(namiSeasonList.Competitions))
	backupSeasonInfoList(namiSeasonList.Competitions)

	return
}

func backupSeasonAreaList(areaList []data.NamiAreaInfo) {

	tx, err := namiDatabase.Begin()
	if err != nil {
		log.Printf("backupSeasonAreaList Begin err : %+v\n", err)
		return
	}
	stmt, err := namiDatabase.Prepare(tx, "INSERT INTO area SET id=? , name_zh=? , name_zht=? , name_en=? "+
		"ON DUPLICATE KEY UPDATE name_zh=? , name_zht=? , name_en=? ")
	if err != nil {
		log.Printf("backupSeasonAreaList Prepare err : %+v\n", err)
		return
	}
	for _, areaInfo := range areaList {

		_, err = stmt.Exec(areaInfo.Id, areaInfo.Name_zh, areaInfo.Name_zht, areaInfo.Name_en, areaInfo.Name_zh, areaInfo.Name_zht, areaInfo.Name_en)
		if err != nil {
			log.Printf("backupSeasonAreaList UPDATE DB err : %+v\n", err)
			continue
		}
	}

	namiDatabase.Commit(tx)
	return
}

func backupSeasonCountryList(countryList []data.NamiCountryInfo) {

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
			log.Printf("backupSeasonCountryList UPDATE DB err : %+v\n", err)
			continue
		}
	}
	namiDatabase.Commit(tx)
	return
}

func backupSeasonInfoList(seasonList []data.NamiSeasonListInfo) {

	tx, err := namiDatabase.Begin()
	if err != nil {
		log.Printf("backupSeasonInfoList Begin err : %+v\n", err)
		return
	}
	stmt, err := namiDatabase.Prepare(tx, "INSERT INTO league SET id = ? , area_id = ? , country_id = ? , type = ? , level = ? , name_zh = ?  , short_name_zh = ? , name_zht = ?  , short_name_zht = ? , name_en = ? , short_name_en = ? , logo = ? , seasonList = ? "+
		"ON DUPLICATE KEY UPDATE seasonList = ? ")
	if err != nil {
		log.Printf("backupSeasonInfoList Prepare err : %+v\n", err)
		return
	}
	for _, seasonInfo := range seasonList {

		for _, seasonId := range seasonInfo.Seasons {
			seasonIdList = append(seasonIdList, seasonId.Id)
		}
		seasonJson, _ := json.Marshal(seasonInfo.Seasons)

		_, err = stmt.Exec(seasonInfo.Id, 0, 0, 0, 0, "", "", "", "", "", "", "", string(seasonJson), string(seasonJson))
		if err != nil {

			log.Printf("backupSeasonInfoListseasonInfo.Id : %+v\n", seasonInfo.Id)
			log.Printf("backupSeasonInfoList UPDATE DB err : %+v\n", err)
			continue
		}
	}
	namiDatabase.Commit(tx)
	return
}
