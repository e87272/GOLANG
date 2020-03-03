package ginEngine

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var sessionStore sessions.Store
var GinEngine *gin.Engine

func GinInit() {

	// 初始化引擎
	gin.SetMode(gin.ReleaseMode)
	GinEngine = gin.Default()

	GinEngine.Use(cors.Default())

	// healthcheck
	GinEngine.StaticFile("/", "./client/healthCheck.html")
	GinEngine.Static("/unitTest", "./unitTest")

	// 初始化session
	InitSession()

	return
}
