package api

import (
	"../commonData"
	"../commonData/league"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var subscribeListLeagueTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/subscribeList/league`: `取得聯賽訂閱列表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`leagueUuid`: `聯賽uuid`,
			`leagueName`: `聯賽名稱`,
			`teamList`:   `參賽球隊列表`,
		},
	},
}

func subscribeListLeague(context *gin.Context) {

	exceptionPrefix := "API_SUBSCRIBELISTLEAGUE_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	rows, err := database.Query("SELECT `target_uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ?",
		userUuid, "league",
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	leagueList := []commonData.LeagueInfo{}
	for rows.Next() {
		var targetUuid string
		rows.Scan(&targetUuid)
		leagueInfo, ok := league.GetLeagueInfo(targetUuid)
		if ok {
			leagueList = append(leagueList, leagueInfo)
		}
	}

	sendResultOk(context, leagueList)
	return
}
