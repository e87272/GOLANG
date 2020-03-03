package api

import (
	"../commonFunc"
	"../external/ginEngine"
	"../external/stamp"
	"github.com/gin-gonic/gin"
)

var tokenTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/token`: `取得token(登入websocket用)`,
	},
	Input:  jsObj{},
	Output: `token(登入websocket用)`,
}

func token(context *gin.Context) {

	userUuid, _ := ginEngine.GetAuthSession(context, "userUuid")

	token := commonFunc.Token{
		Timestamp: stamp.Now(),
		Content:   userUuid,
	}
	tokenStr := commonFunc.AesTokenEncrypt(token)

	sendResultOk(context, tokenStr)
	return
}
