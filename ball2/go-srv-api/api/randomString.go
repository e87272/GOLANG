package api

import (
	"net/http"

	"../commonFunc"

	"github.com/gin-gonic/gin"
)

var randomStringTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/randomString`: `隨機字串(產生密鑰用)`,
	},
	Input:  jsObj{},
	Output: nil,
}

func randomString(context *gin.Context) {

	randString := commonFunc.RandomString(16)
	context.String(http.StatusOK, randString)
	return
}
