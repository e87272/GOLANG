package ginEngine

import (
	"net/http"
	"os"
	"server/common"

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

	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	config.AllowCredentials = true
	config.AddAllowHeaders("token")
	GinEngine.Use(cors.New(config))

	GinEngine.Use(esLog)

	// healthcheck
	GinEngine.GET("/healthcheck", healthCheck)

	return
}

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "healthCheck ver : "+os.Getenv("version"))
}

func esLog(context *gin.Context) {

	if context.HandlerName() != "server/external/ginEngine.healthCheck" {
		context.PostForm("")
		common.SysLog(map[string]interface{}{
			"name":  context.HandlerName(),
			"form":  context.Request.Form,
			"token": context.Request.Header.Get("token"),
		})
	}

	// Pass on to the next-in-chain
	context.Next()
}
