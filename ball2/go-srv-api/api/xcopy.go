package api

import (
	"../commonFunc"
	"github.com/gin-gonic/gin"
)

func xcopy(context *gin.Context) {

	// 使用 c.Param(key) 获取 url 参数
	//announcementType := context.Param("type")

	var announcementInfo interface{}

	announcementInfo = commonFunc.GetLaunchScreen()

	sendResultOk(context, announcementInfo)
	return
}
