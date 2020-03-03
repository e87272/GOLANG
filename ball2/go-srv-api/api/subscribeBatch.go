package api

import (
	"strings"

	"../commonData/game"
	"../commonData/league"
	"../commonData/player"
	"../commonData/team"
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var subscribeBatchTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/subscribeBatch/game`:   `批次訂閱比賽(重設比賽訂閱列表)`,
		`/subscribeBatch/league`: `批次訂閱聯賽(重設聯賽訂閱列表)`,
		`/subscribeBatch/player`: `批次訂閱球員(重設球員訂閱列表)`,
		`/subscribeBatch/team`:   `批次訂閱球隊(重設球隊訂閱列表)`,
	},
	Input: jsObj{
		`targetUuid`: `新的訂閱列表，逗號分隔`,
	},
	Output: nil,
}

func subscribeBatch(context *gin.Context) {

	exceptionPrefix := "API_SUBSCRIBEBATCH_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}
	targetType := context.Param("type")
	targetUuidListStr := context.PostForm("targetUuid")

	targetUuidList := []string{}
	if len(targetUuidListStr) > 0 {
		targetUuidList = strings.Split(targetUuidListStr, ",")
	}

	newSubscriptionAry := []string{}
	switch targetType {
	case "game":
		for _, targetUuid := range targetUuidList {
			_, ok = game.GetGameInfo(targetUuid)
			if ok {
				newSubscriptionAry = append(newSubscriptionAry, targetUuid)
			}
		}
	case "league":
		for _, targetUuid := range targetUuidList {
			_, ok = league.GetLeagueInfo(targetUuid)
			if ok {
				newSubscriptionAry = append(newSubscriptionAry, targetUuid)
			}
		}
	case "player":
		for _, targetUuid := range targetUuidList {
			_, ok = player.GetPlayerInfo(targetUuid)
			if ok {
				newSubscriptionAry = append(newSubscriptionAry, targetUuid)
			}
		}
	case "team":
		for _, targetUuid := range targetUuidList {
			_, ok = team.GetTeamInfo(targetUuid)
			if ok {
				newSubscriptionAry = append(newSubscriptionAry, targetUuid)
			}
		}
	default:
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"TYPE_ERROR", userUuid, nil))
		return
	}

	oldSubscriptionMap := map[string]string{}
	rows, err := database.Query("SELECT `uuid`, `target_uuid` FROM `subscription` WHERE `user_uuid` = ? AND `type` = ?",
		userUuid, targetType,
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}
	for rows.Next() {
		var uuid string
		var targetUuid string
		rows.Scan(&uuid, &targetUuid)
		oldSubscriptionMap[targetUuid] = uuid
	}

	count := 0
	insertAry := []interface{}{}
	for _, targetUuid := range newSubscriptionAry {

		_, ok := oldSubscriptionMap[targetUuid]
		if ok {
			oldSubscriptionMap[targetUuid] = ""
			continue
		}

		uuid := commonFunc.GetUuid()
		insertAry = append(insertAry, uuid, userUuid, targetType, targetUuid)
		count++
	}
	if count > 0 {

		valueStr := strings.Repeat(",(?,?,?,?)", count)[1:]
		_, err = database.Exec("INSERT INTO `subscription` (`uuid`, `user_uuid`, `type`, `target_uuid`) VALUES "+valueStr,
			insertAry...,
		)
		if err != nil {
			sendResultErr(context, commonFunc.Exception(exceptionPrefix+"INSERT_DB_ERROR", userUuid, err))
			return
		}
	}

	count = 0
	deleteAry := []interface{}{}
	for _, uuid := range oldSubscriptionMap {
		if uuid != "" {
			deleteAry = append(deleteAry, uuid)
			count++
		}
	}
	if count > 0 {

		valueStr := strings.Repeat(",?", count)[1:]
		_, err = database.Exec("DELETE FROM `subscription` WHERE `uuid` IN ("+valueStr+")",
			deleteAry...,
		)
		if err != nil {
			sendResultErr(context, commonFunc.Exception(exceptionPrefix+"DELETE_DB_ERROR", userUuid, err))
			return
		}
	}

	sendResultOk(context, nil)
	return
}
