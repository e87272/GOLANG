package api

import (
	"../commonData"
	"../commonData/game"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var subscribeListGameTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/subscribeList/game`: `取得比賽訂閱列表`,
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

func subscribeListGame(context *gin.Context) {

	exceptionPrefix := "API_SUBSCRIBELISTGAME_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	rows, err := database.Query("SELECT `target_uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ?",
		userUuid, "game",
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	gameList := []commonData.GameInfo{}
	for rows.Next() {
		var targetUuid string
		rows.Scan(&targetUuid)
		gameInfo, ok := game.GetGameInfo(targetUuid)
		if ok {
			gameList = append(gameList, gameInfo)
		}
	}

	sendResultOk(context, gameList)
	return
}
