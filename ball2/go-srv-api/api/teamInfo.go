package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var teamInfoTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/teamInfo`: `球隊資訊`,
	},
	Input: jsObj{
		`teamUuid`: `球隊uuid`,
	},
	Output: jsObj{
		`teamCore`: jsObj{
			`teamUuid`: `球隊uuid`,
			`name`:     `球隊名稱`,
		},
		`manager`:    `教練`,
		`venue`:      `主場`,
		`found`:      `成立時間`,
		`playerList`: `球員列表`,
	},
}

func teamInfo(context *gin.Context) {

	teamUuid := context.PostForm("teamUuid")

	clientTeamInfo, ok, exception := commonFunc.ClientTeamInfoSearch(teamUuid)
	if !ok {
		sendResultErr(context, exception)
		return
	}

	sendResultOk(context, clientTeamInfo)
	return
}
