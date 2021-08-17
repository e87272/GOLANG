package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"server/api"
	"server/common/ipData"
	"server/external/database"
	"server/external/ginEngine"
)

func main() {

	loadEnv()

	database.Linkdb()

	ipData.Init()

	ginEngine.GinInit()

	api.ApiRouter()

	// 绑定端口，然后启动应用
	ginEngine.GinEngine.Run(os.Getenv("serverPort"))

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
		os.Setenv(k, v)
	}
	log.Printf("serverPort : %+v", os.Getenv("serverPort"))
	log.Printf("version : %+v", os.Getenv("version"))
}
