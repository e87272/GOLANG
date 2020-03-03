package syncApi

import (
	"../commonData/game"
	"../commonFunc"
	"../external/stamp"
	"github.com/gin-gonic/gin"
)

var syncScheduleTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/schedule`: `同步賽程表資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncSchedule(context *gin.Context) {
	game.ResetGameMap()
	today := stamp.Today()
	queryStartTime := today - stamp.Day
	queryEndTime := today + 7*stamp.Day
	commonFunc.QuerySchedule(queryStartTime, queryEndTime)

	sendResultOk(context, "ok")
	return
}
