package api

import (
	"strconv"

	"../commonData"
	"../commonFunc"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var scheduleTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/schedule`: `賽程表`,
	},
	Input: jsObj{
		`date`:       `日期(毫秒時間戳)，若為空則查近期比賽`,
		`leagueUuid`: `聯賽uuid，若為空則查所有聯賽`,
	},
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

func schedule(context *gin.Context) {

	exceptionPrefix := "API_SCHEDULE_"

	userUuid, _ := ginEngine.GetAuthSession(context, "userUuid")
	dateStr := context.PostForm("date")
	leagueUuid := context.PostForm("leagueUuid")

	gameList := []commonData.GameInfo{}
	if dateStr == "" {
		gameList = commonFunc.GetAllSchedule(leagueUuid)
	} else {
		date, err := strconv.ParseInt(dateStr, 10, 64)
		if err != nil {
			sendResultErr(context, commonFunc.Exception(exceptionPrefix+"DATE_ERROR", userUuid, err))
			return
		}
		gameList = commonFunc.GetSchedule(date, leagueUuid)
	}

	sendResultOk(context, gameList)
	return
}
