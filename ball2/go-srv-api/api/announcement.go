package api

import (
	"../commonData"
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

var announcementTestCase = testCase{
	Method: `GET`,
	Title: jsObj{
		`/announcement/HB`: `首頁橫幅`,
		`/announcement/LS`: `起始頁`,
	},
	Input: jsObj{},
	Output: []jsObj{
		{
			`uuid`:     `公告uuid`,
			`type`:     `"HB"或"LS"`,
			`sequence`: `順序`,
			`content`:  `公告內容`,
			`url`:      `連結網址`,
		},
	},
}

func announcement(context *gin.Context) {

	exceptionPrefix := "API_ANNOUNCEMENT_"

	// 使用 c.Param(key) 获取 url 参数
	announcementType := context.Param("type")
	var announcementInfo = []commonData.Announcement{}

	switch announcementType {
	case "HB":
		announcementInfo = commonFunc.GetHomeBanner()
	case "LS":
		announcementInfo = commonFunc.GetLaunchScreen()
	default: //default case
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"TYPE_ERROR", "", nil))
		return
	}

	sendResultOk(context, announcementInfo)
	return
}
