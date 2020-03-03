package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var scheduleHomePageTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/schedule/homePage`: `首頁賽程表`,
	},
	Input: jsObj{},
	Output: jsObj{
		`gameList`: []jsObj{
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
		`gameTotal`: `近期比賽總數`,
	},
}

func scheduleHomePage(context *gin.Context) {

	gameList := commonFunc.GetAllSchedule("")
	gameTotal := len(gameList)

	if len(gameList) > 4 {
		gameList = gameList[0:4]
	}

	schedule := make(map[string]interface{})
	schedule["gameTotal"] = gameTotal
	schedule["gameList"] = gameList

	sendResultOk(context, schedule)
	return
}
