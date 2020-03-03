package namiAdaptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../commonFunc"
	"../data"
	"../data/teamData"
	"../external/database"
)

var rankRegion = [6]string{"欧洲足联", "南美洲足联", "中北美洲及加勒比海足协", "非洲足联", "亚洲足联", "大洋洲足联"}

func rank() {

	fmt.Println(" rank...")

	for {
		now := time.Now()
		// 计算下一天
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 13, 30, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		rankFifaUpdate()
		rank()
	}

}

func rankFifaUpdate() {

	var rankList data.NamiRankListFifa

	resp, err := commonFunc.GetNamiApi("/ranking/fifa", nil)
	if err != nil {
		log.Printf("rankFifaUpdate err : %+v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err : %+v\n", err)
		return
	}
	// log.Printf("body : %+v\n", string(body))

	if err := json.Unmarshal(body, &rankList); err != nil {
		log.Printf("rankFifaUpdate json.Unmarshal err : %+v\n", err)
		return
	}

	for _, rankInfo := range rankList.Items {
		rankInfoUpdate(rankInfo)
		// time.Sleep(time.Duration(1) * time.Second)
	}
}

func rankInfoUpdate(rankInfo data.NamiRankInfoFifa) {

	teamUuid, err := teamData.SearchTeamUuidByNami(rankInfo.Team.Id)

	if err != nil {
		teamInfo, ok := teamData.GetNamiTeamListData(rankInfo.Team.Id)
		if !ok {
			log.Printf("rankInfoUpdate select err : %+v\n", err)
			return
		}
		teamInfoUpdate(teamInfo.Id, teamInfo.Matchevent_id)
	}

	var rankUuid string
	row := database.QueryRow("SELECT `uuid` FROM `rank` WHERE rank_type = 'fifa' AND `team_uuid` = ?",
		teamUuid,
	)
	err = row.Scan(&rankUuid)

	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("rankInfoUpdate select err : %+v\n", err)
			return
		}
		rankUuid = commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `rank` (`uuid` ,`team_uuid` ,`rank_type` , `rank`, `score`, `region`) VALUES (?, ?, ?, ?, ?, ?)",
			rankUuid, teamUuid, "fifa", rankInfo.Ranking, rankInfo.Points, rankRegion[rankInfo.Region_id],
		)

		if err != nil {
			log.Printf("rankInfoUpdate update DB err : %+v\n", err)
			return
		}
		return

	}

	_, err = database.Exec("UPDATE `rank` SET `rank` = ? , `score` = ? , `region` = ? WHERE `uuid` = ?",
		rankInfo.Ranking, rankInfo.Points, rankRegion[rankInfo.Region_id], rankUuid,
	)

	if err != nil {
		log.Printf("rankInfoUpdate update DB err : %+v\n", err)
		return
	}

	// log.Println("sucess")

}
