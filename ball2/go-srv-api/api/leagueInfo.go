package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var leagueInfoTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/leagueInfo`: `取得聯賽資訊`,
	},
	Input: jsObj{
		`leagueUuid`: `聯賽uuid`,
	},
	Output: jsObj{
		`leagueCore`: jsObj{
			`leagueUuid`: `聯賽uuid`,
			`leagueName`: `聯賽名稱`,
		},
		`teamList`: `參賽球隊列表`,
	},
}

func leagueInfo(context *gin.Context) {

	leagueUuid := context.PostForm("leagueUuid")

	clientLeagueInfo, ok, exception := commonFunc.ClientLeagueInfoSearch(leagueUuid)
	if !ok {
		sendResultErr(context, exception)
		return
	}

	sendResultOk(context, clientLeagueInfo)
	return
}
