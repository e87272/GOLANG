package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var teamListTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/teamList`: `球隊列表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`teamCore`: jsObj{
				`teamUuid`: `球隊uuid`,
				`name`:     `球隊名稱`,
			},
			`playerCount`: `球員數量`,
		},
	},
}

func teamList(context *gin.Context) {

	clientTeamList, ok, exception := commonFunc.ClientTeamList()
	if !ok {
		sendResultErr(context, exception)
		return
	}

	sendResultOk(context, clientTeamList)
	return
}
