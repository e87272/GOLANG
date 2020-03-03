package api

import (
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var changePasswordTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/changePassword`: `更換密碼`,
	},
	Input: jsObj{
		`oldPassword`: `更換前的密碼`,
		`newPassword`: `更換後的密碼`,
	},
	Output: nil,
}

func changePassword(context *gin.Context) {

	exceptionPrefix := "API_CHANGEPASSWORD_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}
	oldPassword := context.PostForm("oldPassword")
	newPassword := context.PostForm("newPassword")

	if len(newPassword) == 0 {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NEW_PASSWORD_ERROR", userUuid, nil))
		return
	}

	res, err := database.Exec("UPDATE `account` SET `password` = ? WHERE `uuid` = ? AND `password` = ?",
		newPassword, userUuid, oldPassword,
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"UPDATE_DB_ERROR", userUuid, err))
		return
	}

	count, err := res.RowsAffected()
	if err == nil && count == 0 {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"PASSWORD_ERROR", userUuid, err))
		return
	}

	sendResultOk(context, nil)
	return
}
