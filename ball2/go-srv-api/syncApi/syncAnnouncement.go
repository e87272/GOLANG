package syncApi

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var syncAnnouncementTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/sync/announcement`: `同步公告資訊`,
	},
	Input:  jsObj{},
	Output: nil,
}

func syncAnnouncement(context *gin.Context) {

	commonFunc.QueryAnnouncement()

	sendResultOk(context, "ok")
	return
}
