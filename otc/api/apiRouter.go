package api

import (
	"regexp"
	"server/common"
	"server/external/ginEngine"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var merchant string

func ApiRouter() {

	// 按功能排列

	// 创建基于cookie的存储引擎，参数是用于加密的密钥
	store := cookie.NewStore([]byte(common.GetUuid()))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	ginEngine.GinEngine.Use(sessions.Sessions("session", store))

	//ip檢查
	ginEngine.GinEngine.Use(authCheck)

	//backend api
	backendApiGroup := ginEngine.GinEngine.Group("/backend")

	backendApiGroup.Use(ipCheck)

	//client api
	clientApiGroup := ginEngine.GinEngine.Group("/client")

}

func authCheck(context *gin.Context) {

	session := sessions.Default(context)

	if strings.Contains(context.Request.URL.Path, "/login") {
		return
	}

	userUuid := session.Get("userUuid")
	sessionToken := session.Get("token")

	dbToken, ok := memberData.GetMember(userUuid)
	if !ok {
		name := "authCheck GetMember err"
		common.SysErrorLog(map[string]interface{}{
			"name":     name,
			"userUuid": userUuid,
			"token":    token,
		}, nil)
		sendResultErr(context, name)
		context.Abort()
		return
	}

	if sessionToken != dbToken {
		name := "authCheck token err"
		common.SysErrorLog(map[string]interface{}{
			"name":     name,
			"userUuid": userUuid,
			"token":    token,
		}, nil)
		sendResultErr(context, name)
		context.Abort()
		return
	}

}

func authCheck(context *gin.Context) {

	var ok bool
	var err error

	ok, err = regexp.MatchString(`^10\.\d+\.\d+\.\d+$`, context.ClientIP())
	if err != nil {
		name := "authCheck regexp MatchString err"
		common.SysErrorLog(map[string]interface{}{
			"name": name,
			"ip":   context.ClientIP(),
		}, err)
		sendResultErr(context, MERCHANT_IP_ERR)
		context.Abort()
		return
	}

	if ok {
		return
	}

	_, ok = merchantData.GetWhiteList(context.ClientIP())
	if !ok {
		name := "ip not in white list"
		common.SysErrorLog(map[string]interface{}{
			"name":     name,
			"merchant": merchant,
			"ip":       context.ClientIP(),
		}, nil)
		sendResultErr(context, MERCHANT_IP_ERR)
		context.Abort()
		return
	}

}
