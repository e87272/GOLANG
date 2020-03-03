package namiAdaptor

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../commonFunc"
)

func gameDetail() {

	fmt.Println(" gameDetail...")

	// for _, teamInfo := range teamList {
	// 	teamInfoUpdate(teamInfo)
	// 	time.Sleep(time.Duration(1) * time.Second)
	// }

	// time.Sleep(time.Duration(1) * time.Second)

	// gameDetail()

	var gameLiveTicker = time.NewTicker(time.Second * 10)

	defer gameLiveTicker.Stop()

	for {
		<-gameLiveTicker.C

		data := make(map[string]string)

		data["id"] = "2787265"
		resp, err := commonFunc.GetNamiApi("/match/detail", nil)
		if err != nil {
			log.Printf("GetApiForm err : %+v\n", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err : %+v\n", err)
			return
		}
		log.Printf("body : %+v\n", string(body))

		// fmt.Println("body gameDetail...")
	}

}
