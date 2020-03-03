package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"math/rand"
	"net/http"
	"os"
	"time"

	"./commonFunc"
	"./external/database"
	"./external/elasticSearch"
	"./external/namiDatabase"
	"./external/redis"
	"./namiAdaptor"
	"./namiBackup"
)

func main() {

	rand.Seed(time.Now().Unix())

	// log.Printf("NumCPU : %v\n", runtime.NumCPU())

	// runtime.GOMAXPROCS(runtime.NumCPU())

	// log.Printf(common.Getid().Hexstring())

	loadEnv()

	elasticSearch.EsInit()

	redis.RedisInit()

	database.Linkdb()
	namiDatabase.Linkdb()

	commonFunc.InitFunc()

	namiAdaptor.InitAdaptorTick()
	namiBackup.InitNamiBackup()

	http.HandleFunc("/", serveHome)

	http.ListenAndServe(os.Getenv("goServerAdaptorPort"), nil)

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./client/healthCheck.html")
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
