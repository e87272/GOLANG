package api

import (
	"../commonData"
	"../commonData/team"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var subscribeListTeamTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/subscribeList/team`: `取得球隊訂閱列表`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`teamCore`: jsObj{
				`teamUuid`: `球隊uuid`,
				`name`:     `球隊名稱`,
			},
			`manager`:    `教練`,
			`venue`:      `主場`,
			`found`:      `成立時間`,
			`playerList`: `球員列表`,
		},
	},
}

func subscribeListTeam(context *gin.Context) {

	exceptionPrefix := "API_SUBSCRIBELISTTEAM_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	rows, err := database.Query("SELECT `target_uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ?",
		userUuid, "team",
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	teamList := []commonData.TeamInfo{}
	for rows.Next() {
		var targetUuid string
		rows.Scan(&targetUuid)
		teamInfo, ok := team.GetTeamInfo(targetUuid)
		if ok {
			teamList = append(teamList, teamInfo)
		}
	}

	sendResultOk(context, teamList)
	return
}
