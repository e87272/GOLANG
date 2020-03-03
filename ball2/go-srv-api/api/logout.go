package api

import (
	"../commonFunc"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var logoutTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/logout`: `登出`,
	},
	Input:  jsObj{},
	Output: nil,
}

func logout(context *gin.Context) {

	exceptionPrefix := "API_LOGOUT_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	ginEngine.DeleteAuthSession(context, "userUuid")
	sendResultOk(context, nil)
	return
}
