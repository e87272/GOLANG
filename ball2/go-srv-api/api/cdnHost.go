package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

var cdnHostTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/cdnHost`: `取得cdn路徑`,
	},
	Input:  jsObj{},
	Output: nil,
}

func cdnHost(context *gin.Context) {

	sendResultOk(context, os.Getenv("cdnHost"))
	return
}
