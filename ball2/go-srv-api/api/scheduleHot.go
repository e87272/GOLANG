package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var scheduleHotTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/schedule/hot`: `熱門賽程表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`uuid`:          `比賽uuid`,
			`homeTeamUuid`:  `主隊uuid`,
			`guestTeamUuid`: `客隊uuid`,
			`startTime`:     `開始時間`,
			`leagueUuid`:    `所屬聯賽uuid`,
			`homeScore`:     `主隊得分`,
			`guestScore`:    `客隊得分`,
			`subtitle`:      `子標題`,
			`state`:         `比賽狀態`,
		},
	},
}

func scheduleHot(context *gin.Context) {

	sendResultOk(context, commonFunc.GetHotSchedule())
	return
}
