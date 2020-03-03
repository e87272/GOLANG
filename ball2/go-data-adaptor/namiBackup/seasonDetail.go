package namiBackup

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"../commonFunc"
	"../data"
)

func backupSeasonDetail() {

	log.Printf("backupSeasonDetail...\n")
	log.Printf("seasonIdList : %+v\n", len(seasonIdList))

	for _, seasonId := range seasonIdList {

		var seasonInfo data.NamiSeasonInfo
		data := make(map[string]string)
		data["id"] = strconv.FormatInt(int64(seasonId), 10)

		resp, err := commonFunc.GetNamiApi("/season/detail", data)
		if err != nil {
			log.Printf("backupSeasonDetail err : %+v\n", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err : %+v\n", err)
			return
		}
		// log.Printf("body : %+v\n", string(body))

		if err := json.Unmarshal(body, &seasonInfo); err != nil {
			log.Printf("backupSeasonDetail json.Unmarshal err : %+v\n", err)
			return
		}
	}

	return
}
