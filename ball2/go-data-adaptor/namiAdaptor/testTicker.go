package namiAdaptor

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../commonFunc"
)

func test() {

	fmt.Println(" test...")
	var testTicker = time.NewTicker(time.Second * 60)

	defer testTicker.Stop()

	for {
		<-testTicker.C

		data := make(map[string]string)

		data["id"] = "2787265"
		resp, err := commonFunc.GetNamiApi("https://open.sportnanoapi.com/api/sports/football/match/lineup", data)
		if err != nil {
			log.Printf("GetApiForm err : %+v\n", err)
			return
		}

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll err : %+v\n", err)
			return
		}
		// log.Printf("body test: %+v\n", string(body))

		fmt.Println("body test...")
	}

}
