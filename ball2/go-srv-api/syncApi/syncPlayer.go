package syncApi

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var syncPlayerTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/player`: `同步球員資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncPlayer(context *gin.Context) {

	commonFunc.QueryPlayer()

	sendResultOk(context, "ok")
	return
}
