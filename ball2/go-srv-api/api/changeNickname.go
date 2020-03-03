package api

import (
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var changeNicknameTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/changeNickname`: `更換暱稱`,
	},
	Input: jsObj{
		`nickname`: `新的暱稱`,
	},
	Output: nil,
}

func changeNickname(context *gin.Context) {

	exceptionPrefix := "API_CHANGENICKNAME_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}
	nickname := context.PostForm("nickname")

	_, err := database.Exec("UPDATE `account` SET `nickname` = ? WHERE `uuid` = ?", nickname, userUuid)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"UPDATE_DB_ERROR", userUuid, err))
		return
	}

	sendResultOk(context, nil)
	return
}
