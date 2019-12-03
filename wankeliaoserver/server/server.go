package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"./api"
	"./command"
	"./command/commandRoom"
	"./common"
	"./database"
	"./socket"
)

func main() {

	rand.Seed(time.Now().Unix())

	log.SetFlags(log.LstdFlags)

	// log.Printf("NumCPU : %v\n", runtime.NumCPU())

	// runtime.GOMAXPROCS(runtime.NumCPU())

	// log.Printf(common.Getid().Hexstring())

	loadEnv()

	loadLang("zh-CN")
	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	// log.Printf("main timeUnix : [%s] ", timeUnix)
	var err error
	common.Elasticclient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv("elasticSearchHost")))

	for err != nil {
		log.Printf("Elasticsearch err %+v\n", err)
		now := time.Now().UnixNano()
		for time.Now().UnixNano() <= now+1e9 {
		}
		common.Elasticclient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv("elasticSearchHost")))
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := common.Elasticclient.ElasticsearchVersion(os.Getenv("elasticSearchHost"))
	if err != nil {
		log.Printf("Elasticsearch err %+v\n", err)
	}
	log.Printf("Elasticsearch version %s\n", esversion)

	// common.Esdelete(os.Getenv("sysErrorLog"))
	common.Essyserrorinit(os.Getenv("sysLog"))

	// common.Esdelete(os.Getenv("sysErrorLog"))
	common.Essyserrorinit(os.Getenv("sysErrorLog"))

	// common.Esdelete(os.Getenv("sysRoomLog"))
	common.Eschatinit(os.Getenv("sysRoomLog"))

	// common.Esdelete(os.Getenv("sideText"))
	// common.Esdelete(os.Getenv("liveGroup"))
	// common.Esdelete(os.Getenv("vipGroup"))
	// common.Esdelete(os.Getenv("privateGroup"))

	common.Eschatinit(os.Getenv("sideText"))
	common.Eschatinit(os.Getenv("liveGroup"))
	common.Eschatinit(os.Getenv("vipGroup"))
	common.Eschatinit(os.Getenv("privateGroup"))

	// common.Esdelete("roomdirtywordhistory")
	// common.Esdelete("sidetextdirtywordhistory")
	common.Esdirtywordhistoryinit("roomdirtywordhistory")
	common.Essidetextdirtywordhistoryinit("sidetextdirtywordhistory")

	common.Redisclient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redisHost"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go common.Subscriberoom()
	go common.Subscribeuser()
	go common.Subscribesidetext()
	go common.Subscriberoomsinfo()
	go common.Subscribesync()
	go common.Subscribesudoresult()

	database.Linkdb()
	common.Queryblocklist()
	common.Querydirtyword()
	common.Queryfunctionmanagement()
	common.Queryglobalmessage()

	api.Api()

	go common.Alivecheck()
	go common.Servertick()

	http.HandleFunc("/", serveHome)
	if os.Getenv("environmentId") != "Online" {
		http.HandleFunc("/unitTest/", fileUnitTestHandler)
	}
	http.HandleFunc("/echo", echoHandler)

	http.ListenAndServe(os.Getenv("imServerPort"), nil)

}

func fileUnitTestHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("r.URL.Path :", r.URL.Path)
	http.ServeFile(w, r, ".."+r.URL.Path)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {

	loginUuid := common.Getid().Hexstring()

	// log.Printf("loginUuid : %s \n", loginUuid)
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Essyserrorlog("MAIN_ECHOHANDLER_ERROR", loginUuid, err)
		return
	}

	defer func() {
		// log.Printf("defer : %+v\n", loginUuid)

		// log.Println("disconnect !!")
		client, _ := common.Clientsread(loginUuid)
		userPlatform := client.Userplatform
		room := client.Room

		// log.Printf("room : %+v - %s\n", room, loginUuid)

		for _, roomInfo := range room {

			// log.Printf("roominfo : %+v\n", roominfo)

			common.Roomsclientdelete(roomInfo.Roomuuid, loginUuid)
			// log.Printf("delete : " + loginUuid)

			if len(common.Roomsread(roomInfo.Roomuuid)) == 0 {
				common.Roomsdelete(roomInfo.Roomuuid)
				common.Roomsinfodelete(roomInfo.Roomuuid)
			}

			// 離開為單一不用通知
			// historyUuid := common.Getid().Hexstring()
			// timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			// chatMessage := socket.Chatmessage{Historyuuid: historyUuid, From: userPlatform, Stamp: timeUnix, Message: "exit room", Style: "sys"}
			// roomBroadcast := socket.Cmd_b_player_room{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}}
			// roomBroadcast.Payload.Chatmessage = chatMessage
			// roomBroadcast.Payload.Chattarget = roomInfo.Roomuuid
			// roomBroadcastJson, _ := json.Marshal(roomBroadcast)

			// common.Redispubroomsinfo(roomInfo.Roomuuid, roomBroadcastJson)

		}

		common.Clientsdelete(loginUuid)
		clientConnect, _ := common.Clientsconnectread(userPlatform.Useruuid)
		if len(clientConnect) == 1 {
			common.Clientsconnectdelete(userPlatform.Useruuid)
			common.Usersinfodelete(userPlatform.Useruuid)
			common.Userfriendlistdelete(userPlatform.Useruuid)
		} else {
			common.Clientsconnectloginuuiddelete(userPlatform.Useruuid, loginUuid)
		}

		// connect.Close()
		// log.Printf("connect : close - %s \n", loginUuid)
	}()

	for {
		err = receivePacketHandle(connect, loginUuid)
		if err != nil {
			// log.Println("echoHandler write:", err)
			break
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../client/healthCheck.html")
}

func receivePacketHandle(connect *websocket.Conn, loginUuid string) error {

	// log.Printf("connect : %+v\n", common.Clientsread(connect))

	//ReadMessage只能讀一次 猜測是因為讀取指標的問題

	_, msg, err := connect.ReadMessage()
	if err != nil {
		// log.Printf("connect : %+v\n", err)
		return err
	}

	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	// log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
		//log.Println("receivePacketHandle Unmarshal:", err)
		return err
	}

	userPlatform, ok := common.Clientsuserplatformread(loginUuid)
	if !ok {
		common.Essyslog(string(msg), loginUuid, "")
	} else {
		common.Essyslog(string(msg), loginUuid, userPlatform.Useruuid)
	}

	switch mapResult["cmd"] {
	case socket.CMD_C_TOKEN_CHANGE:
		// log.Printf("CMD_C_TOKEN_CHANGE : " + uuid)
		err = command.Tokenchange(connect, msg, loginUuid)

		if err != nil {
			return err
		}
		break
	case socket.CMD_C_PLAYER_LOGOUT:

		err = command.Playerlogout(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_MEMBER_LIST:

		err = commandRoom.Getmemberlist(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_EXIT_ROOM:

		err = commandRoom.Playerexitroom(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_ENTER_ROOM:

		err = commandRoom.Playerenterroom(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_ROOM_INFO_EDIT:

		err = commandRoom.Roominfoedit(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_CHAT_HISTORY:

		err = commandRoom.Getroomhistory(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_FRIEND_INVITE:

		err = command.Friendinvite(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_FRIEND_LIST:

		err = command.Getfriendlist(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_CHATBLOCK:

		err = commandRoom.Blockroomchat(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PING:

		err = command.Healthcheck(connect, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PROCLAMATION:

		err = command.Proclamationsearch(connect, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_FRIEND_DELETE:

		err = command.Frienddelete(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_SIDETEXT_DELETE:

		err = command.Sidetextdelete(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_SIDETEXT_HISTORY:

		err = command.Getsidetexthistory(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_NEW_SIDETEXT:

		err = command.Getnewsidetext(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_FUNC_MANAGEMENT:

		err = command.Getfuncmanagement(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_TARGET_ADD_ROOM_BATCH:

		err = commandRoom.Targetaddroombatch(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_ENTER_ROOM_BATCH:

		err = commandRoom.Playerenterroombatch(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_MESSAGE_SEEN:

		err = command.Messageseen(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_KICK_ROOM_USER:

		err = commandRoom.Kickroomuser(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_CREATE_PRIVATE_ROOM:

		err = commandRoom.Createprivateroom(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_SEND_MSG:

		err = commandRoom.Playersendmsg(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_SIDETEXT:

		err = command.Sidetextsend(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_SEND_SHELL:

		err = command.Playersendshell(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_DIS_MISS_ROOM:

		err = commandRoom.Dismissroom(connect, msg, loginUuid)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_ROOM_ADMIN_ADD:

		err = commandRoom.Roomadminadd(connect, msg, loginUuid)

		if err != nil {
			return err
		}
		break
	case socket.CMD_C_ROOM_ADMIN_REMOVE:

		err = commandRoom.Roomadminremove(connect, msg, loginUuid)

		if err != nil {
			return err
		}
		break
	case socket.CMD_C_GET_LANG_LIST:

		err = command.Getlanglist(connect, msg, loginUuid)

		if err != nil {
			return err
		}
		break
	default:

	}

	return nil
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

func loadLang(lang string) {

	langFile, err := ioutil.ReadFile("config/" + lang + ".json")
	if err != nil {
		log.Fatal("找不到host.json")
	}

	langMap := make(map[string]string)

	err = json.Unmarshal(langFile, &langMap)
	if err != nil {
		// log.Printf("configHost err: %v\n", err)
		return
	}

	common.Mutexmutilangerrormsg.Lock()
	defer common.Mutexmutilangerrormsg.Unlock()
	common.Mutilangerrormsg[lang] = langMap

}
