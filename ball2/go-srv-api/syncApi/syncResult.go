package syncApi

import (
	"encoding/json"
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
	apiResultJson, _ := json.Marshal(apiResult)
	context.String(http.StatusOK, string(apiResultJson))
}

func sendResultErr(context *gin.Context, exception commonData.Exception) {
	apiResult := apiResult{
		Result: "err",
		Error:  exception,
	}
	apiResultJson, _ := json.Marshal(apiResult)
	context.String(http.StatusOK, string(apiResultJson))
}
