package api

import (
	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"github.com/gin-gonic/gin"
)

var passwordLoginTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/passwordLogin`: `透過密碼登入`,
	},
	Input: jsObj{
		`phone`:       `手機號碼`,
		`countryCode`: `國碼(不含加號)`,
		`password`:    `密碼`,
	},
	Output: nil,
}

func passwordLogin(context *gin.Context) {

	exceptionPrefix := "API_PASSWORDLOGIN_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"ALREADY_LOGGED_IN", userUuid, nil))
		return
	}
	phone := context.PostForm("phone")
	countryCode := context.PostForm("countryCode")
	password := context.PostForm("password")

	if password == "" {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"PASSWORD_IS_NULL", userUuid, nil))
		return
	}

	row := database.QueryRow("SELECT `uuid` FROM `account` WHERE `phone` = ? AND `country` = ? AND `password` = ?",
		phone, countryCode, password,
	)
	err := row.Scan(&userUuid)
	if err == database.ErrNoRows {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"PHONE_OR_PASSWORD_ERROR", userUuid, err))
		return
	}
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	ginEngine.SetAuthSession(context, "userUuid", userUuid)
	sendResultOk(context, nil)
	return
}
