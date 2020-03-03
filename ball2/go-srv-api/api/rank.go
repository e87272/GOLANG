package api

import (
	"../commonData/ranking"
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var rankTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/rank/fifa`: `FIFA排名`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`teamCore`: jsObj{
				`teamUuid`: `球隊uuid`,
				`name`:     `球隊名稱`,
			},
			`score`: `積分`,
		},
	},
}

func rank(context *gin.Context) {

	exceptionPrefix := "API_RANK_"

	rankType := context.Param("type")

	list, ok := ranking.GetRank(rankType)
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"TYPE_ERROR", "", nil))
		return
	}

	sendResultOk(context, list)
	return
}
