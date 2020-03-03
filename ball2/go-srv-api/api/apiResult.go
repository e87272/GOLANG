package api

import (
	"net/http"

	"../commonData"
	"github.com/gin-gonic/gin"
)

type apiResult struct {
	Result  string               `json:"result"`
	Error   commonData.Exception `json:"error"`
	Payload interface{}          `json:"payload"`
}

func sendResultOk(context *gin.Context, payload interface{}) {
	apiResult := apiResult{
		Result:  "ok",
		Payload: payload,
	}
	context.JSON(http.StatusOK, apiResult)
}

func sendResultErr(context *gin.Context, exception commonData.Exception) {
	apiResult := apiResult{
		Result: "err",
		Error:  exception,
	}
	context.JSON(http.StatusOK, apiResult)
}
