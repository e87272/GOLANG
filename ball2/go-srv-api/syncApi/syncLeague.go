package syncApi

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var syncLeagueTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/league`: `同步聯賽資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncLeague(context *gin.Context) {

	commonFunc.QueryLeague()

	sendResultOk(context, "ok")
	return
}
