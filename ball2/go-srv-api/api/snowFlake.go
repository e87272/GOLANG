package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var snowFlakeTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/snowFlake`: `取得新的uuid`,
	},
	Input:  jsObj{},
	Output: `新的uuid`,
}

func snowFlake(context *gin.Context) {

	sendResultOk(context, commonFunc.GetUuid())
	return
}
