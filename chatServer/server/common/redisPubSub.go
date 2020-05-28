package common

import (
	"encoding/json"
	"context"
	"strconv"
	"time"

	"server/socket"
)

func Subscriberoom() {

	// log.Printf("Subscriberoom \n")

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"room")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {
		// log.Printf("Subscriberoom err : %+v\n", err)
		return
	}
	ch := Redispubsub.Channel()

	for msg := range ch {
		var data Redispubsubroomdata

		err = json.Unmarshal([]byte(msg.Payload), &data)
		if err != nil {
			panic(err)
		}

		var playerSpeak socket.Cmd_b_player_speak

		err := json.Unmarshal([]byte(data.Datajson), &playerSpeak)
		if err != nil {
			panic(err)
		}

		packetStamp, _ := strconv.ParseInt(playerSpeak.Stamp, 10, 64)

		switch playerSpeak.Payload.Chatmessage.Style {
		case "sys", "gift", "subscription":
		default:
			_, ok := Usersinforead(playerSpeak.Payload.Chatmessage.From.Useruuid)

			if ok {
				Isspeakcd(playerSpeak.Payload.Chatmessage.From.Useruuid, playerSpeak.Payload.Chatmessage.Stamp)
			}
		}

		// log.Printf("Subscriberoom : %+v\n", data.Datajson)
		// log.Printf("Subscriberoom time: %+v\n", time.Now())

		go Broadcast(data.RoomUuid, []byte(data.Datajson), packetStamp)

	}

}

func Subscribeuser() {

	// log.Printf("Subscriberoom \n")

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"user")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {
		// log.Printf("Subscribeuser err : %+v\n", err)
		return
	}
	ch := Redispubsub.Channel()

	for msg := range ch {
		var data Redispubsubuserdata

		err = json.Unmarshal([]byte(msg.Payload), &data)
		if err != nil {
			panic(err)
		}

		// log.Printf("Subscriberoom : %+v\n", data.Datajson)
		// log.Printf("Subscriberoom time: %+v\n", time.Now())

		user, ok := Clientsconnectread(data.Useruuid)

		if ok {
			for _, connect := range user {
				Sendmessage(connect, []byte(data.Datajson))

			}
		}
	}

}

func Subscribesidetext() {

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"sideText")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {

		return
	}
	ch := Redispubsub.Channel()
	for msg := range ch {
		var data Redispubsubsidetextdata

		err = json.Unmarshal([]byte(msg.Payload), &data)
		if err != nil {
			panic(err)
		}
		// log.Printf("Subscribesidetext data : %+v\n", data)
		Sidetextsend(data.Useruuid, data.Targetuuid, []byte(data.Datajson))
	}

}

func Subscriberoomsinfo() {

	// log.Printf("Subscriberoomsinfo \n")

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"roomsinfo")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {
		// log.Printf("Subscriberoomsinfo err : %+v\n", err)
		return
	}
	ch := Redispubsub.Channel()

	for msg := range ch {
		var data Redispubsubroomsinfo

		err = json.Unmarshal([]byte(msg.Payload), &data)
		if err != nil {
			panic(err)
		}

		var playerRoom socket.Cmd_b_player_room

		err := json.Unmarshal([]byte(data.Datajson), &playerRoom)
		if err != nil {
			panic(err)
		}
		packetStamp, _ := strconv.ParseInt(playerRoom.Stamp, 10, 64)

		// log.Printf("Subscriberoomsinfo : %+v\n", data)
		Roomspopulationcount(data)
		go Broadcast(data.Roomuuid, []byte(data.Datajson), packetStamp)

	}

}

func Subscribesync() {

	// log.Printf("Subscriberoomsinfo \n")

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"sync")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {
		// log.Printf("Subscribesync err : %+v\n", err)
		return
	}
	ch := Redispubsub.Channel()
	for msg := range ch {
		var syncData Syncdata

		err = json.Unmarshal([]byte(msg.Payload), &syncData)
		if err != nil {
			panic(err)
		}
		//更新列表並廣播
		switch syncData.Synctype {
		case "blockSync":
			go Queryblocklist()
			go Queryblockiplist()
		case "roomsInfoSync":
			go Syncroominfo(syncData.Data)
		case "apiKeySync":
			go Queryapikey()
		case "proclamationSync":
			go Queryproclamation(syncData.Data)
		case "sideTextSync":
			go Updatesidetextlist(syncData.Data)
		case "sideTextDeleteSync":
			go Deletesidetextlist(syncData.Data)
		case "dirtywordSync":
			go Querydirtyword()
		case "userInfoSync":
			go Queryuserinfo(syncData.Data)
		case "funcManagementSync":
			go Queryfunctionmanagement()
		case "userInfoSyncAndEmit":
			go Userinfosyncandemit(syncData.Data)
		case "friendUpdateState":
			go Updatefriendstate(syncData.Data)
		case "globalMessageSync":
			go Queryglobalmessage()
		case "clearUserMsg":
			go Sendclearusermsg(syncData.Data)
		case "memberCountSync":
			go Membercountbroadcast(syncData.Data)
		}
	}
}

func Subscribesudoresult() {

	// log.Printf("Subscribesudoresult \n")

	//参数1 频道名 字符串类型
	ctx := context.Background()
	Redispubsub := Redisclient.Subscribe(ctx,"sudoresult")

	defer Redispubsub.Close()

	_, err := Redispubsub.Receive(ctx)
	if err != nil {
		// log.Printf("Subscribesudoresult err : %+v\n", err)
		return
	}
	ch := Redispubsub.Channel()
	for msg := range ch {
		var sudoShellData struct {
			Sudoresult socket.Sudoresult
			Funcname   string
		}

		// log.Printf("sudoResult msg : %+v\n", msg)

		err = json.Unmarshal([]byte(msg.Payload), &sudoShellData)
		if err != nil {
			panic(err)
		}

		// log.Printf("sudoResult : %+v\n", sudoResult)
		timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		sudoShellDataPacket := socket.Cmd_b_admin_shell{Base_B: socket.Base_B{Cmd: socket.CMD_B_ADMIN_SHELL, Stamp: timeUnix}, Payload: sudoShellData.Sudoresult}
		sudoShellDataPacketJson, _ := json.Marshal(sudoShellDataPacket)
		BroadcastAdmin(sudoShellData.Sudoresult.Shelltarget, []byte(sudoShellDataPacketJson), sudoShellData.Funcname)
	}
}
