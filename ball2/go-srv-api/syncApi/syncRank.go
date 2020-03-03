package syncApi

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var syncRankTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/rank`: `同步排名資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncRank(context *gin.Context) {

	commonFunc.QueryRank()

	sendResultOk(context, "ok")
	return
}
