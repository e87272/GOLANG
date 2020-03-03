package syncApi

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var syncTeamTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/team`: `同步球隊資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncTeam(context *gin.Context) {

	commonFunc.QueryTeam()

	sendResultOk(context, "ok")
	return
}
