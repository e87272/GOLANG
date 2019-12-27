package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "34.80.226.253", "http service address") //prod

// var addr = flag.String("addr", "127.0.0.1", "http service address")  //beta

// var addr = flag.String("addr", "34.80.25.124", "http service address") //alpha

// var addr = flag.String("addr", "127.0.0.1", "http service address") //alpha

var mutexSendMessage sync.Mutex

var mutexSetConnect sync.Mutex

var roomsAmount int64 = 3

type Conncore struct {
	Conn  *websocket.Conn
	Mutex *sync.Mutex
}

func (cc Conncore) Close() error {
	cc.Mutex.Lock()
	defer cc.Mutex.Unlock()
	if cc.Conn == nil {
		log.Println("Close cc.Conn : nil")
		return nil
	}
	err := cc.Conn.Close()
	cc.Conn = nil
	return err
}

func (cc Conncore) WriteMessage(data []byte) error {
	// log.Println("cc.Mutex.Lock()")
	// cc.Mutex.Lock()
	// defer log.Println("cc.Mutex.Unlock()")
	// defer cc.Mutex.Unlock()
	if cc.Conn == nil {
		log.Println("WriteMessage cc.Conn : nil")
		return nil
	}
	cc.Conn.WriteMessage(websocket.TextMessage, data)
	return nil
}

func (cc Conncore) ReadMessage() ([]byte, error) {
	cc.Mutex.Lock()
	defer cc.Mutex.Unlock()
	if cc.Conn == nil {
		log.Println("ReadMessage cc.Conn : nil")
		return []byte{}, nil
	}
	_, p, err := cc.Conn.ReadMessage()
	return p, err
}

func main() {
	rand.Seed(time.Now().Unix())

	var i int64
	for i = 0; i < 600; i++ {
		setConnect(i)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			log.Println("t:", t)
		case <-interrupt:
			log.Printf("interrupt")
			select {
			case t := <-time.After(time.Second):
				log.Printf("time: %+v", t)
			}
			return
		}
	}
}

func setConnect(count int64) {
	mutexSetConnect.Lock()
	defer mutexSetConnect.Unlock()

	select {
	case <-time.After(time.Millisecond * 100):
		log.Println("main:", count)
		go wsConnect(count)
	}
}

func wsConnect(count int64) {

	defer setConnect(count)
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	//log.Printf("connecting to %s", u.String())
	h := http.Header{}
	h.Add("X-Forwarded-For", "127.0.0.1")
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		log.Printf("dial: %+v", err)
		return
	}
	cc := Conncore{
		Conn:  conn,
		Mutex: new(sync.Mutex),
	}
	defer func() {
		t := time.Duration(int(time.Millisecond) * (5000 + rand.Intn(5000)))
		select {
		case <-time.After(t):
			cc.Close()
		}
	}()

	wsChan := make(chan bool)
	go func() {
		flag := true
		for flag {
			flag = receivePacketHandle(cc, count)
		}
		wsChan <- true
	}()

	// log.Printf("wsConnect count : %+v", count)
	tokenChange(cc, count)

	// time.Sleep(time.Duration(count*100) * time.Millisecond)
	ticker := time.NewTicker(time.Second * 10)
	// ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	defer log.Printf("wsConnect closed")

	for {
		select {
		case <-ticker.C:
			// log.Printf("ticker.C count : %+v", count)
			// sendChatMessage(cc, count)
		case <-wsChan:
			log.Printf("wsChan : %+v", count)
			return
		}
	}
}

var userAry = []string{
	"4RDDBn6z::5d4a380582e7c", //1
	"5CrWr1HY::5d4a380587c59", //2
	"zsTOCxKe::5d4a3805951e4", //3
	"Dq4kO7F8::5d4a380586ac8", //4
	"KBkRA0DT::5d4a380586f05", //5
	"LhPlxdUH::5d4a380587308", //8
	"7As6aO4U::5d4a38e321c23", //13
	"MgR5Kigg::5d4a380581167", //16
	"aJsx7QbB::5d4a38e3226ab", //17
	"624EBKXf::5d4a3805863d1",
	"5d75dd4c-4dc5-c075-654e-3e6a-183a7920",
	"QNdiYdYg::5d4a38058983e",
	"bl1qWP7h::5d4a3805868aa",
	"rbgsCR3G::5d4a38e322c3b",
	"dmylFFP3::5d4a3805a6b23",
	"5d96b255-39ab-959c-58f1-bc37-46346f2e",
	"5d5364b7-d974-eb65-8192-3ac6-04d5990b",
	"LOYcdVzI::5d4a3805b55ec",
	"5d773eb2-df6e-b0c1-a758-7592-e43f4080",
	"5d634e81-ef0d-0943-6305-2dd1-abe787cd",
	"5d9860b6-75da-da9d-dc63-43d9-8e211f88",
	"5d639374-b81c-0ae9-dbe0-c051-01e8dc14",
	"atZstmny::5d4a3805b7813",
	"5d9af1a8-f2ba-ec41-a742-8748-2345c264",
	"5da42dca-0120-eba1-61a1-1082-3523a585",
	"8c7loL28::5d4a380594415",
	"5d96c1cc-6b41-ca1c-7c23-03a0-6b928a6a",
	"vSru7gBP::5d4a38059a46f",
	"RGXGVlid::5d4a3805a67b5",
	"5d970f29-a755-42b7-2d2a-98cd-c674c0d6",
	"l5k3hf9s::5d4a3805a64be",
	"5d637698-1509-fa50-1a99-4476-07a5d69f",
	"4RDDBn6z::5d4a380582458",
	"4RDDBn6z::5d4a380582403",
	"j8ULJHSa::5d4a380585ef3",
	"5daff819-4890-5334-fc31-a1d5-2ae1be47",
	"5db00ea5-faee-4454-0908-0a45-daf2652c",
	"5d807ac8-ad0a-a592-7a73-b200-f676a8f7",
	"5daff631-9f4f-b229-3407-6f67-fdc108c5",
	"5d908644-76de-0e3d-06b3-1821-e892acb5",
	"5da53b3f-744a-5e63-fe6d-c055-3c308c2f",
	"5db015f4-afd7-b17e-6efd-d5e8-2fdb4992",
	"5db0162c-9c51-1de1-74e9-1279-246045ef",
	"5d64ad31-771e-cc23-5573-60f6-2b620e26",
	"5db01691-5625-4281-84e2-f74c-8aebdf2f",
	"5db01bbb-a838-bfbd-ba94-2bae-12947850",
	"5db2a062-7cb7-d16b-ab03-41e3-4e1165f9",
	"5db2aaae-eb81-5134-059e-328a-d4d2968e",
	"5db2c928-bfb4-8a7c-622f-855b-f15127be",
	"5db6b08c-9498-cef0-d12d-348d-80a0b8ed",
	"5dc0e08c-0da2-e5a3-12c2-67ed-644ccd4f",
	"5dc27083-b7f8-33c0-7ced-166f-9bf2836b",
	"4RDDBn6z::5d4a380582492",
	"5d9d6e21-2fa9-7e1d-ad42-7fce-acf01878",
	"5d919dc6-58bf-4b1d-bc49-1066-84443087",
	"KNx3wi1Y::5d4a38059fffb",
	"5d4baad6-2bdf-a14f-e352-88cd-c6ffecfe",
	"5d68bd7a-9787-d2c2-28f7-2962-dbf61b41",
	"4RDDBn6z::5d4a380582406",
	"5dd346a6-ede4-3120-3ff0-ea3b-08952ef3",
	"qdn3L1KN::5d4a3805a03a8",
	"5d9bdfad-ad53-a5b3-c312-0d8b-c3f7b5fe",
	"4RDDBn6z::5d4a380582489",
	"swH90LpZ::5d4a380597b6f",
	"5d9de360-8d47-41bd-a4eb-df9d-5d89819f",
	"4RDDBn6z::5d4a380582418",
	"5d919dd1-cac0-0fe0-505d-e4c9-f55f8913",
	"W3LqGtZT::5d4a3805a5b6a",
	"5d83263f-95c5-15dd-4357-cf2b-a3a675a1",
	"5d9ebfa5-4db0-1dbe-2940-973d-f4f9c9ff",
	"4RDDBn6z::5d4a380582408",
	"5d9eb192-b53a-f93c-4296-984e-ce6430e4",
	"5d919de2-c783-a8e6-7af6-16c3-a435a505",
	"5d97bd03-18ed-bb6b-f4a6-2f4d-2196e8ba",
	"5d935f94-6b7a-d6be-8841-8ba4-e0000c34",
	"5d9eb640-a400-90c0-e798-c56d-e0d55d38",
	"5d9ec8ff-890d-30d8-8456-ec61-2429a9f8",
	"5d9eb640-6264-ef98-c6a6-6004-78f9fbac",
	"5d9ecdb3-c231-8019-d3e4-dce3-561d75a5",
	"5d9ed258-7c0b-f813-7200-b64f-0e8590e4",
	"4RDDBn6z::5d4a380582415",
	"4RDDBn6z::5d4a380582495",
	"5d9eb187-01fc-8ce1-44a2-0498-e070a407",
	"5d9eb640-61a3-d483-109a-f998-b4633dd3",
	"JzxM6CCS::5d4a380594d3c",
	"4RDDBn6z::5d4a380582411",
	"4RDDBn6z::5d4a380582413",
	"3cB3DvMl::5d4a3805a62ce",
	"5dd4b6c9-d61a-2782-0100-a477-16a1d98d",
	"5dd4d9eb-b149-f400-daeb-58a3-44b5b3f8",
	"5d9f5a80-ca2f-d057-d461-4e41-3e856050",
	"tUgp95PV::5d4a3805bba73",
	"5d9f5a83-4955-f859-e084-74f9-d2c35ce8",
	"5d9f5a82-f816-c5c4-b820-ca01-67617fe6",
	"5d9f5a56-4e66-a4e6-d38f-bc71-3ca82c40",
	"5d8c54ac-e30e-a43d-bb9e-604b-2a50a179",
	"5d9de377-0622-a5be-dde1-492d-5ec98bc5",
	"5d9e041c-5213-ae39-17d4-dfa0-371de0a0",
	"5d9f5a74-2625-d751-101d-84e2-7d91b01c",
	"aJsx7QbB::5d4a38e3226ab",
	"5d9de39b-35a2-cfea-f115-9521-b6ca0b67",
	"9udWJI5Q::5d4a39396b603",
	"5da933f8-6943-5fa6-a370-60e4-0ca8d05a",
	"5d9f5a88-b9d5-58d5-1da2-6a9e-6589d2ad",
	"5d9f5a66-1502-dcff-861e-dfe8-8cb0c740",
}

var roomsUuid = []string{}

func tokenChange(cc Conncore, i int64) {
	msg := map[string]interface{}{}
	msg["cmd"] = "2"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["platform"] = "MM"
	payload["platformUuid"] = userAry[i%int64(len(userAry))] //prod

	// log.Printf("tokenChange userAry i : %+v", i)
	// log.Printf("tokenChange userAry[i] : %+v", userAry[i])
	// log.Printf("tokenChange platformUuid : %+v", payload["platformUuid"])
	// payload["platformUuid"] = "5db00be4-b7fa-6e69-46e7-866a-42f161c0" //beta

	payload["token"] = "solar"
	msg["payload"] = payload

	packetMsg, _ := json.Marshal(msg)

	// log.Printf("tokenChange : %s", packetMsg)

	cc.Conn.WriteMessage(websocket.TextMessage, packetMsg)
}

func enterRoom(cc Conncore, count int64) {

	msg := map[string]interface{}{}
	msg["cmd"] = "10"
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["station"] = "smzb2"
	payload["ownerPlatform"] = "MM"
	if roomsAmount > int64(len(userAry)) {
		roomsAmount = int64(len(userAry))
	}
	payload["Ownerplatformuuid"] = userAry[count%roomsAmount]

	// beta
	// payload["station"] = "team2"
	// payload["ownerPlatform"] = "team2"
	// payload["Ownerplatformuuid"] = "team2"

	payload["roomName"] = ""
	payload["adminSet"] = ""
	msg["payload"] = payload

	packetMsg, _ := json.Marshal(msg)

	cc.WriteMessage(packetMsg)
}

func sendChatMessage(cc Conncore, count int64) {
	msg := map[string]interface{}{}
	msg["cmd"] = "80"
	// msg["cmd"] = strconv.FormatInt(rand.Intn(150), 10)
	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	msg["idem"] = timeUnix
	payload := map[string]interface{}{}
	payload["chatTarget"] = roomsUuid[count%roomsAmount] //prod
	// payload["chatTarget"] = "000572bf8d4a5001" //beta
	payload["message"] = timeUnix
	payload["style"] = "string"
	msg["payload"] = payload

	packetMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("sendChatMessage json err:", err)
		return
	}

	err = cc.WriteMessage(packetMsg)
	if err != nil {
		log.Println("sendChatMessage WriteMessage err:", err)
		return
	}

	// log.Println("Unlock")
}

func receivePacketHandle(cc Conncore, count int64) bool {

	msg, err := cc.ReadMessage()
	if err != nil {
		return false
	}
	//log.Printf("recv msg: %s", msg)

	//timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	////log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	err = json.Unmarshal([]byte(msg), &mapResult)
	if err != nil {
		log.Println("receivePacketHandle Unmarshal:", err)
		return true
	}

	// log.Printf("mapResult : %+v\n", mapResult)

	switch mapResult["cmd"] {
	case "3":
		log.Printf("wsConnect connected")
		enterRoom(cc, count)
		break
	case "11":
		var packetPlayerEnterRoom struct {
			Cmd     string `json:"cmd"`
			Idem    string `json:"idem"`
			Payload struct {
				Roomcore struct {
					Roomuuid string `json:"roomUuid"`
					Roomtype string `json:"roomType"`
				} `json:"roomCore"`
			} `json:"payload"`
		}
		if err := json.Unmarshal([]byte(msg), &packetPlayerEnterRoom); err != nil {
			log.Println("packetPlayerEnterRoom Unmarshal:", err)
		}
		roomsUuid = append(roomsUuid, packetPlayerEnterRoom.Payload.Roomcore.Roomuuid)
		break
	case "81":
		// log.Printf("mapResult : %+v\n", mapResult)
		break
	}
	return true
}
