package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var userInfoTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/userInfo`: `使用者資訊`,
	},
	Input: jsObj{
		`userUuid`: `使用者uuid`,
	},
	Output: jsObj{
		`nickname`: `使用者暱稱`,
		`icon`:     `使用者頭像`,
	},
}

func userInfo(context *gin.Context) {

	userUuid := context.PostForm("userUuid")

	userInfo, ok, exception := commonFunc.UserInfoSearch(userUuid)
	if !ok {
		sendResultErr(context, exception)
		return
	}

	var clientGetUserInfo struct {
		Nickname string `json:"nickname"`
		Icon     string `json:"icon"`
	}

	clientGetUserInfo.Nickname = userInfo.Nickname
	clientGetUserInfo.Icon = userInfo.Icon
	sendResultOk(context, clientGetUserInfo)
	return
}
