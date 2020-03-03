package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var leagueListTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/leagueList`: `取得聯賽列表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`leagueCore`: jsObj{
				`leagueUuid`: `聯賽uuid`,
				`leagueName`: `聯賽名稱`,
			},
			`sequence`:  `畫面排序`,
			`teamCount`: `參賽隊伍數量`,
		},
	},
}

func leagueList(context *gin.Context) {

	clientLeagueList, ok, exception := commonFunc.ClientLeagueList()
	if !ok {
		sendResultErr(context, exception)
		return
	}

	sendResultOk(context, clientLeagueList)
	return
}
