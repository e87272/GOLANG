package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiResult struct {
	Result  string      `json:"result"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

func sendResultOk(context *gin.Context, payload interface{}) {
	apiResult := apiResult{
		Result:  "ok",
		Payload: payload,
	}
	context.JSON(http.StatusOK, apiResult)
}

func sendResultErr(context *gin.Context, exception string) {
	apiResult := apiResult{
		Result: "err",
		Error:  exception,
	}
	context.JSON(http.StatusOK, apiResult)
}
