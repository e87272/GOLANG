package api

import (
	"../commonFunc"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var myInfoTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/myInfo`: `自己的使用者資訊`,
	},
	Input: jsObj{
		`userUuid`: `使用者uuid`,
		`nickname`: `使用者暱稱`,
		`phone`:    `使用者手機號碼`,
		`icon`:     `使用者頭像`,
	},
	Output: nil,
}

func myInfo(context *gin.Context) {

	exceptionPrefix := "API_MYINFO_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	userInfo, ok, exception := commonFunc.UserInfoSearch(userUuid)
	if !ok {
		sendResultErr(context, exception)
		return
	}

	sendResultOk(context, userInfo)
	return
}
