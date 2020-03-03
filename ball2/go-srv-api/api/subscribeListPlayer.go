package api

import (
	"../commonData"
	"../commonData/player"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var subscribeListPlayerTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/subscribeList/player`: `取得球員訂閱列表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`playerUuid`:  `球員uuid`,
			`playerName`:  `球員名稱`,
			`countryUuid`: `所屬國家uuid`,
		},
	},
}

func subscribeListPlayer(context *gin.Context) {

	exceptionPrefix := "API_SUBSCRIBELISTPLAYER_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	rows, err := database.Query("SELECT `target_uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ?",
		userUuid, "player",
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	playerList := []commonData.PlayerInfo{}
	for rows.Next() {
		var targetUuid string
		rows.Scan(&targetUuid)
		playerInfo, ok := player.GetPlayerInfo(targetUuid)
		if ok {
			playerList = append(playerList, playerInfo)
		}
	}

	sendResultOk(context, playerList)
	return
}
