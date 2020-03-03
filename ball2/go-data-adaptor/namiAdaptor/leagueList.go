package namiAdaptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../commonFunc"
	"../data"
	"../data/continentData"
	"../data/countryData"
	"../data/leagueData"
	"../external/database"
)

var league_type = [4]string{"未知", "联赛", "杯赛", "友谊赛"} //赛事类型 0-未知 1-联赛 2-杯赛 3-友谊赛

func leagueList() {

	fmt.Println(" leagueList...")

	for {
		now := time.Now()
		// 计算下一天
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 13, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		leagueListUpdate()
		leagueList()
	}

}

func leagueListUpdate() {

	var namiLeagueList data.NamiLeagueList

	resp, err := commonFunc.GetNamiApi("/matchevent/list", nil)
	if err != nil {
		log.Printf("leagueListUpdate err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &namiLeagueList); err != nil {
		log.Printf("leagueListUpdate json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("Areas : %+v\n", len(namiLeagueList.Areas))
	areaInfoMap(namiLeagueList.Areas)

	log.Printf("Countrys : %+v\n", len(namiLeagueList.Countrys))
	countryInfoMap(namiLeagueList.Countrys)

	log.Printf("LeagueList : %+v\n", len(namiLeagueList.Matchevents))
	leagueInfoUpdate(namiLeagueList.Matchevents)
}

func leagueInfoUpdate(leagueList []data.NamiLeagueInfo) {

	continentUuid, err := continentData.SearchContinentUuidByNami(0)
	if err != nil {
		log.Printf("SearchContinentUuidByNami : %+v\n", 0)
		log.Printf("leagueInfoUpdate select continent err : %+v\n", err)
		return
	}

	countryUuid, err := countryData.SearchCountryUuidByNami(0)
	if err != nil {
		log.Printf("SearchCountryUuidByNami : %+v\n", 0)
		log.Printf("leagueInfoUpdate select country err : %+v\n", err)
		return
	}

	leagueUuid, err := leagueData.SearchLeagueUuidByNami(0)
	//撈不到資料預塞聯賽編號0
	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("leagueInfoUpdate select err : %+v\n", err)
			return
		}

		leagueUuid = commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `league` (`uuid` ,`name` ,`name_en` ,`short_name` ,`continent_uuid` ,`country_uuid` , `type` , `level` , `sequence`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			leagueUuid, "未知", "unknown", "未知", continentUuid, countryUuid, "未知", 10, 9999,
		)
		if err != nil {
			log.Printf("leagueInfoUpdate INSERT league err : %+v\n", err)
			return
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "league", leagueUuid, "nami", 0,
		)

		if err != nil {
			log.Printf("leagueInfoUpdate INSERT corporation_id err : %+v\n", err)
			return
		}

	}

	for _, leagueInfo := range leagueList {
		// log.Printf("leagueInfo : %+v\n", leagueInfo)
		// time.Sleep(time.Duration(1) * time.Second)
		continentUuid, err := continentData.SearchContinentUuidByNami(leagueInfo.Area_id)
		if err != nil {
			log.Printf("Area_id : %+v\n", leagueInfo.Area_id)
			log.Printf("countryInfoUpdate select continent err : %+v\n", err)
			continue
		}

		countryUuid, err := countryData.SearchCountryUuidByNami(leagueInfo.Country_id)
		if err != nil {
			log.Printf("leagueInfo : %+v\n", leagueInfo)
			log.Printf("leagueInfoUpdate select country err : %+v\n", err)
			continue
		}

		leagueUuid, err := leagueData.SearchLeagueUuidByNami(leagueInfo.Id)
		if err != nil {
			if err != database.ErrNoRows {
				log.Printf("leagueInfoUpdate select err : %+v\n", err)
				continue
			}
			leagueUuid = commonFunc.GetUuid()
			_, err := database.Exec("INSERT INTO `league` (`uuid` ,`name` ,`name_en` ,`short_name` ,`continent_uuid` ,`country_uuid` , `type` , `level` , `sequence`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
				leagueUuid, leagueInfo.Name_zh, leagueInfo.Name_en, leagueInfo.Short_name_zh, continentUuid, countryUuid, league_type[leagueInfo.Type], leagueInfo.Level, leagueInfo.Id,
			)
			if err != nil {
				log.Printf("leagueInfoUpdate INSERT league err : %+v\n", err)
				continue
			}

			corporationUuid := commonFunc.GetUuid()
			_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
				corporationUuid, "league", leagueUuid, "nami", leagueInfo.Id,
			)

			if err != nil {
				log.Printf("leagueInfoUpdate INSERT corporation_id err : %+v\n", err)
				continue
			}

			continue

		}

		if leagueInfo.Logo != "" {
			_, err = commonFunc.PostCdnUploadLink(leagueInfo.Logo, leagueUuid+".png", "/leagueIcon/")
			if err != nil {
				log.Printf("leagueInfoUpdate err : %+v\n", err)
			}
		}
		_, err = database.Exec("UPDATE `league` SET `name` = ? , `name_en` = ? , `short_name` = ? , `continent_uuid` = ? , `country_uuid` = ?, `type` = ? , `level` = ?  WHERE `uuid` = ?",
			leagueInfo.Name_zh, leagueInfo.Name_en, leagueInfo.Short_name_zh, continentUuid, countryUuid, league_type[leagueInfo.Type], leagueInfo.Level, leagueUuid)

		if err != nil {
			log.Printf("leagueInfoUpdate update DB err : %+v\n", err)
			continue
		}
	}
	return
}
