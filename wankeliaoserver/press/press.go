package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "34.80.226.253", "http service address") //prod
// var addr = flag.String("addr", "34.80.25.124", "http service address") //prod

// var addr = flag.String("addr", "127.0.0.1", "http service address")  //beta
var mutex sync.Mutex

func main() {
	for i := 0; i < 2000; i++ {
		log.Println("main:", i)
		now := time.Now().UnixNano()
		go wsConnect(i)
		for time.Now().UnixNano() <= now+1e6*10 {
		}
	}

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			log.Println("t:", t)
		}
	}
}
func wsConnect(count int) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			receivePacketHandle(c)
		}
	}()

	tokenChange(c)

	time.Sleep(time.Duration(count*90) * time.Millisecond)
	ticker := time.NewTicker(time.Second * 200)
	// ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()
	// sendChatMessage(c)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// log.Println("t:", t)
			//log.Printf("recv roomInfo: %+v", roomInfo)
			sendChatMessage(c)
		case <-interrupt:
			log.Println("interrupt")
			select {
			case <-done:
			case c := <-time.After(time.Second):
				log.Println("c", c)
			}
			return
		}
	}
}

func tokenChange(c *websocket.Conn) {

	msg := map[string]interface{}{}
	msg["cmd"] = "2"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["platform"] = "MM"
	payload["platformUuid"] = "5daff819-4890-5334-fc31-a1d5-2ae1be47" //prod
	// payload["platformUuid"] = "5db00be4-b7fa-6e69-46e7-866a-42f161c0" //beta

	msg["payload"] = payload

	packetMsg, _ := json.Marshal(msg)

	// log.Printf("tokenChange : %s", packetMsg)

	c.WriteMessage(websocket.TextMessage, packetMsg)
}

func enterRoom(c *websocket.Conn) {

	msg := map[string]interface{}{}
	msg["cmd"] = "10"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["roomType"] = "liveGroup"
	payload["roomUuid"] = "001d9921a5a58000" //prod
	// payload["roomUuid"] = "000572bf8d4a5001" //beta
	payload["roomName"] = ""
	payload["adminSet"] = ""
	msg["payload"] = payload

	packetMsg, _ := json.Marshal(msg)

	c.WriteMessage(websocket.TextMessage, packetMsg)
}

func sendChatMessage(c *websocket.Conn) {
	mutex.Lock()
	// log.Println("Lock")

	msg := map[string]interface{}{}
	msg["cmd"] = "80"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["chatTarget"] = "001d9921a5a58000" //prod
	// payload["chatTarget"] = "000572bf8d4a5001" //beta
	payload["message"] = timeUnix
	payload["style"] = "1"
	msg["payload"] = payload

	packetMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("sendChatMessage json err:", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, packetMsg)
	if err != nil {
		log.Println("sendChatMessage WriteMessage err:", err)
		return
	}

	// log.Println("Unlock")
	mutex.Unlock()
}
func receivePacketHandle(connect *websocket.Conn) {

	_, msg, err := connect.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	//log.Printf("recv msg: %s", msg)

	if err != nil {
		log.Println("receivePacketHandle Readpacket:", err)
	}

	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	////log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
		log.Println("receivePacketHandle Unmarshal:", err)
	}

	// log.Printf("mapResult : %+v\n", mapResult)

	switch mapResult["cmd"] {
	case "3":
		enterRoom(connect)
		break
	case "11":
		break
	case "81":
		// log.Printf("mapResult : %+v\n", mapResult)
		break
	}
}
