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
	"../data/leagueData"
	"../data/teamData"
	"../external/database"
)

func teamList() {

	fmt.Println(" teamList...")

	for {
		now := time.Now()
		// 计算下一天
		next := now.Add(time.Hour * 24)
		//隔天13:30的時間做資料更新
		next = time.Date(next.Year(), next.Month(), next.Day(), 13, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		teamListUpdate()
		teamList()
	}

}

func teamListUpdate() {

	var teamList data.NamiTeamList

	resp, err := commonFunc.GetNamiApi("/team/list", nil)
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

	if err := json.Unmarshal(body, &teamList); err != nil {
		log.Printf("teamList json.Unmarshal err : %+v\n", err)
		return
	}

	log.Printf("teamList len : %+v\n", len(teamList))

	teamData.SetNamiTeamListData(teamList)
	// for _, teamInfo := range teamList {
	// 	// log.Printf("teamInfo : %+v\n", teamInfo)
	// 	log.Printf("teamInfo id : %+v\n", teamInfo.Id)

	// 	_, err := teamData.SearchTeamUuidByNami(teamInfo.Id)

	// 	if err != nil {
	// 		go teamInfoListUpdate(teamInfo)
	// 		time.Sleep(time.Duration(500) * time.Millisecond)
	// 	}
	// }
}

func teamInfoListUpdate(teamListInfo data.NamiTeamListInfo) {

	var teamInfo data.NamiTeamInfo

	data := make(map[string]string)
	data["id"] = strconv.FormatInt(int64(teamListInfo.Id), 10)

	// log.Printf("teamInfo id : %+v\n", int64(teamListInfo.Id))

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

	teamUuid, err := teamData.SearchTeamUuidByNami(teamListInfo.Id)

	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("teamList select err : %+v\n", err)
			return
		}

		teamUuid = commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `team` (`uuid` ,`name` , `name_en` ,`manager` , `venue`, `found`) VALUES (?, ?, ?, ?, ?, ?)",
			teamUuid, teamInfo.Name_zh, teamInfo.Name_en, teamInfo.Manager.Name_zh, teamInfo.Venue.Name_zh, teamListInfo.Found,
		)
		if err != nil {

			log.Printf("teamInfoUpdate insert DB teamInfo : %+v\n", teamInfo)
			log.Printf("teamInfoUpdate insert DB err : %+v\n", err)
			return
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "team", teamUuid, "nami", teamListInfo.Id,
		)
		if err != nil {
			log.Printf("teamInfoUpdate update DB err : %+v\n", err)
			return
		}

		leagueUuid, err := leagueData.SearchLeagueUuidByNami(teamListInfo.Matchevent_id)

		if err != nil {
			log.Printf("namiLeagueId : %+v\n", teamListInfo.Matchevent_id)
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
			resp, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/team/"+teamInfo.Logo, teamUuid+".png", "/teamIcon/")

			body, err := ioutil.ReadAll(resp.Body)
			log.Printf("PostCdnUploadLink err : %+v\n", string(body))

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
		_, err = commonFunc.PostCdnUploadLink("http://cdn.sportnanoapi.com/football/team/"+teamInfo.Logo, teamUuid+".png", "/teamIcon/")
		if err != nil {
			log.Printf("PostCdnUploadLink err : %+v\n", err)
		}
	}
	_, err = database.Exec("UPDATE `team` SET `name` = ? ,`name_en` = ? ,`manager` = ? , `venue` = ? , `found` = ?  WHERE `uuid` = ?",
		teamInfo.Name_zh, teamInfo.Name_en, teamInfo.Manager.Name_zh, teamInfo.Venue.Name_zh, teamListInfo.Found, teamUuid)

	if err != nil {
		log.Printf("teamInfoUpdate update team err : %+v\n", err)
		return
	}

	leagueUuid, err := leagueData.SearchLeagueUuidByNami(teamListInfo.Matchevent_id)

	if err != nil {
		log.Printf("namiLeagueId : %+v\n", teamListInfo.Matchevent_id)
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
