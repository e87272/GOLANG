package api

import (
	"../commonData/game"
	"../commonData/league"
	"../commonData/player"
	"../commonData/team"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var cancelSubscribeTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/cancelSubscribe/game`:   `取消訂閱比賽`,
		`/cancelSubscribe/league`: `取消訂閱聯賽`,
		`/cancelSubscribe/player`: `取消訂閱球員`,
		`/cancelSubscribe/team`:   `取消訂閱球隊`,
	},
	Input: jsObj{
		`targetUuid`: `訂閱對象的uuid`,
	},
	Output: nil,
}

func cancelSubscribe(context *gin.Context) {

	exceptionPrefix := "API_CANCELSUBSCRIBE_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}
	targetType := context.Param("type")
	targetUuid := context.PostForm("targetUuid")

	switch targetType {
	case "game":
		_, ok = game.GetGameInfo(targetUuid)
	case "league":
		_, ok = league.GetLeagueInfo(targetUuid)
	case "player":
		_, ok = player.GetPlayerInfo(targetUuid)
	case "team":
		_, ok = team.GetTeamInfo(targetUuid)
	default:
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"TYPE_ERROR", userUuid, nil))
		return
	}
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"TARGET_UUID_ERROR", userUuid, nil))
		return
	}

	uuid := ""
	row := database.QueryRow("SELECT `uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ? AND `target_uuid` = ?",
		userUuid, targetType, targetUuid,
	)
	err := row.Scan(&uuid)
	if err == database.ErrNoRows {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_SUBSCRIB", userUuid, nil))
		return
	}
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	_, err = database.Exec("DELETE FROM `subscription` WHERE `uuid` = ?",
		uuid,
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"DELETE_DB_ERROR", userUuid, err))
		return
	}

	sendResultOk(context, nil)
	return
}
