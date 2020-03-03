package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"math/rand"
	"os"
	"time"

	"./api"
	"./commonFunc"
	"./external/database"
	"./external/elasticSearch"
	"./external/ginEngine"
	"./external/redis"
	"./syncApi"
)

func main() {

	rand.Seed(time.Now().Unix())

	// log.Printf("NumCPU : %v\n", runtime.NumCPU())

	// runtime.GOMAXPROCS(runtime.NumCPU())

	// log.Printf(common.Getid().Hexstring())

	loadEnv()

	// loadLang("zh-CN")
	timeUnix := time.Now().UnixNano() / int64(time.Millisecond)

	log.Printf("main timeUnix : [%d] ", timeUnix)

	elasticSearch.EsInit()

	redis.RedisInit()

	database.Linkdb()

	ginEngine.GinInit()

	commonFunc.DataInit()
	commonFunc.FuncInit()

	api.ApiRouter()
	syncApi.ApiRouter()

	// 绑定端口，然后启动应用
	log.Printf("goServerApiPort : %v\n", os.Getenv("goServerApiPort"))
	ginEngine.GinEngine.Run(os.Getenv("goServerApiPort"))

}

func loadEnv() {

	config, err := ioutil.ReadFile("config/host.json")
	if err != nil {
		log.Fatal("找不到host.json")
	}

	configHost := make(map[string]string)
	// log.Printf("configHost : %v\n", configHost)

	err = json.Unmarshal(config, &configHost)
	if err != nil {
		// log.Printf("configHost err: %v\n", err)
		return
	}

	// log.Printf("configHost : %v\n", configHost)
	for k, v := range configHost {
		// log.Printf("%s : %s\n", k, v)
		_ = os.Setenv(k, v)
	}
}
