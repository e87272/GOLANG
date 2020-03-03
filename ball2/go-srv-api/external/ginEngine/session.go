package ginEngine

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func InitSession() {
	// 初始化基于redis的存储引擎
	// 参数说明：
	//    第1个参数 - redis最大的空闲连接数
	//    第2个参数 - 数通信协议tcp或者udp
	//    第3个参数 - redis地址, 格式，host:port
	//    第4个参数 - redis密码
	//    第5个参数 - session加密密钥
	SessionStore, err := redis.NewStore(10, "tcp", os.Getenv("redisHost"), "", []byte(os.Getenv("sessionKey")))
	if err != nil {
		log.Printf("SessionStore err : %+v ", err)
	}

	GinEngine.Use(sessions.Sessions("authSession", SessionStore))
}

func SetAuthSession(c *gin.Context, key string, value string) {

	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
	log.Printf("SetAuthSession session : %+v ", session)
}

func GetAuthSession(c *gin.Context, key string) (string, bool) {

	session := sessions.Default(c)
	val := session.Get(key)
	if val == nil {
		return "", false
	}
	log.Printf("GetAuthSession key : %+v ", key)
	log.Printf("GetAuthSession val : %+v ", val)
	return val.(string), true
}

func DeleteAuthSession(c *gin.Context, key string) {

	session := sessions.Default(c)
	session.Delete(key)
	// 保存session数据
	session.Save()
	log.Printf("DeleteAuthSession session : %+v ", session)
}

func ClearAuthSession(c *gin.Context) {

	session := sessions.Default(c)
	session.Clear()
	// 保存session数据
	session.Save()
	log.Printf("ClearAuthSession session : %+v ", session)
}
