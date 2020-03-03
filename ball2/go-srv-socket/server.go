package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"./commonData"
	"./commonFunc"
	"./external/database"
	"./external/elasticSearch"
	"./external/redis"
	"./socket"
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

	commonFunc.InitFunc()
	// api.Api()

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/echo", echoHandler)

	http.ListenAndServe(os.Getenv("goServerSocketPort"), nil)

}

func fileUnitTestHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("r.URL.Path :", r.URL.Path)
	http.ServeFile(w, r, ".."+r.URL.Path)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {

	loginUuid := commonFunc.GetUuid()

	log.Printf("loginUuid : %s \n", loginUuid)
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		commonFunc.EsSysErrorLog("MAIN_ECHOHANDLER_ERROR", loginUuid, err)
		if connect != nil {
			connect.Close()
		}
		return
	}
	//不安全可能被客端竄改還需修改
	//if len(r.Header["X-Forwarded-For"]) == 0 {
	//	commonFunc.EsSysErrorLog("MAIN_ECHOHANDLER_IP_ERROR", loginUuid, err)
	//	connect.Close()
	//	return
	//}
	//commonFunc.IpListInsert(loginUuid, r.Header["X-Forwarded-For"][0])

	defer func() {
		// log.Printf("defer : %+v\n", loginUuid)

		// log.Println("disconnect !!")
		client, _ := commonFunc.ClientsRead(loginUuid)
		room := client.Room

		// log.Printf("room : %+v - %s\n", room, loginUuid)

		for _, roomUuid := range room {

			// log.Printf("roominfo : %+v\n", roominfo)

			commonFunc.RoomsClientDelete(roomUuid, loginUuid)
			// log.Printf("delete : " + loginUuid)
			roomClient, ok := commonFunc.RoomsRead(roomUuid)
			if ok && len(roomClient) == 0 {
				commonFunc.RoomsDelete(roomUuid)
			}

			// 離開為單一不用通知

		}

		commonFunc.ClientsDelete(loginUuid)
		commonFunc.IpListDelete(loginUuid)
		clientConnection, _ := commonFunc.ClientsConnectionsRead(client.UserInfo.UserUuid)
		if len(clientConnection) == 1 {
			commonFunc.ClientsConnectionsDelete(client.UserInfo.UserUuid)
			commonFunc.UsersInfoDelete(client.UserInfo.UserUuid)
		} else {
			commonFunc.ClientsConnectionsLoginUuidDelete(client.UserInfo.UserUuid, loginUuid)
		}

		// connect.Close()
		// log.Printf("connect : close - %s \n", loginUuid)
	}()

	connCore := commonData.ConnCore{Conn: connect, Connmutex: new(sync.Mutex), LoginUuid: loginUuid}

	for {
		socket.SocketHandler(w, r, connCore)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL)
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
