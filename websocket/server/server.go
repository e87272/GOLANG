package main

import (
	"log"
	"flag"
	"net/http"
	"time"
	"strconv"
	"encoding/json"

	"github.com/gorilla/websocket"
	
	"./socket"
)

type Client struct {
	room map[string]socket.RoomInfo

	// The websocket connection.
	conn *websocket.Conn

	user socket.User
}

var addr = flag.String("addr", ":8800", "http service address")
var rooms = make(map[string]map[*websocket.Conn]Client)
var roomsInfo = make(map[string]socket.RoomInfo)
var clients = make(map[*websocket.Conn]Client)

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
	http.ServeFile(w, r, "../client/index.html")
}

func main() {
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

		connect, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		defer func() {
			log.Println("disconnect !!")
			
			for _ , roominfo := range clients[connect].room {
				delete(rooms[roominfo.Roomname] , connect)
				if(len(rooms[roominfo.Roomname]) == 0){
					delete(rooms,roominfo.Roomname)
					delete(roomsInfo,roominfo.Roomname)
				}
			}

			delete(clients , connect)

			connect.Close()
		}()

		for {
			err = receivePacketHandle(connect)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/websocket.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../client/websocket.js")
	})
	
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	
	log.Println("server start at :8800")
	log.Fatal(http.ListenAndServe(":8800", nil))

}


func receivePacketHandle(connect *websocket.Conn) error{

	mtype, msg, err := connect.ReadMessage()
	if err != nil {
		log.Println("write:", err)
		return err
	}


	timeUnix := strconv.FormatInt(time.Now().Unix(),10)

	// log.Printf("timeUnix : [%s] ", timeUnix)

    var mapResult map[string]interface{}
    //使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
    if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
        panic(err)
	}
	
	// log.Printf("mapResult : %+v\n", mapResult)

	switch mapResult["cmd"] {
		case socket.CMD_C_TOKEN_CHANGE:
			var packetToken socket.Cmd_c_token_change_struct

			if err := json.Unmarshal([]byte(msg), &packetToken); err != nil {
				panic(err)
			}

			roomMap := make(map[string]socket.RoomInfo)
			roomMap[packetToken.Payload.Roomname] = socket.RoomInfo { Roomname : packetToken.Payload.Roomname }
			userinfo := socket.User {Id : packetToken.Payload.Nickname , Nickname : packetToken.Payload.Nickname , Icon : "" , Role : "" , Status : ""}
			client := Client{room: roomMap , conn: connect , user : userinfo}
			// log.Printf("client : %+v\n", client)
			clients[connect] = client
			//log.Printf("clients : %+v\n", clients)
			if(rooms[packetToken.Payload.Roomname] == nil){
				var adduser = make(map[*websocket.Conn]Client)
				adduser[connect] = client
				rooms[packetToken.Payload.Roomname] = adduser
				roomsInfo[packetToken.Payload.Roomname] = socket.RoomInfo { Roomname : packetToken.Payload.Roomname }
			}else{
				rooms[packetToken.Payload.Roomname][connect] = client
			}
			// log.Printf("rooms : %+v\n", rooms)

			tokenChange := socket.Cmd_r_token_change_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_TOKEN_CHANGE , Idem : packetToken.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}}}
			tokenChangeJson, _ := json.Marshal(tokenChange)
			
			connect.WriteMessage(mtype ,  tokenChangeJson)
			
			roominfo := socket.RoomInfo { Roomname : packetToken.Payload.Roomname }
			chatMessage := socket.ChatMessage {From : clients[connect].user  , Stamp : timeUnix , Text : "enter room" , Style : "1" , Roominfo : roominfo}
			logoutBrocast := socket.Cmd_b_player_enter_room_struct { Base_B : socket.Base_B {Cmd : socket.CMD_B_PLAYER_ENTER_ROOM , Stamp : timeUnix} , Payload : chatMessage}
			logoutBrocastJson, err := json.Marshal(logoutBrocast)
			if err != nil {
				log.Println("json err:", err)
			}

			broadcast(packetToken.Payload.Roomname, mtype , logoutBrocastJson)

			break;
		case socket.CMD_C_PLAYER_LOGOUT:
			var packetLogout socket.Cmd_c_player_logout_struct

			if err := json.Unmarshal([]byte(msg), &packetLogout); err != nil {
				panic(err)
			}

			if !checkinroom( packetLogout.Payload.Roomname , clients[connect] ){
				logout := socket.Cmd_r_player_logout_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_LOGOUT , Idem : packetLogout.Idem , Stamp : timeUnix, Result : "err", Exp :  socket.Exception{Code : "101", Message : "NOT_IN_ROOM"}}}
				logoutJson, _ := json.Marshal(logout)
				connect.WriteMessage(mtype , logoutJson)
				break;
			}

			logout := socket.Cmd_r_player_logout_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_LOGOUT , Idem : packetLogout.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}}}
			logoutJson, _ := json.Marshal(logout)
			
			connect.WriteMessage(mtype , logoutJson)

			userinfo := clients[connect].user
			

			for _ , roominfo := range clients[connect].room {
				delete(rooms[roominfo.Roomname] , connect)
				if(len(rooms[roominfo.Roomname]) == 0){
					delete(rooms,roominfo.Roomname)
					delete(roomsInfo,roominfo.Roomname)
				}
			}

			delete(clients , connect)

			connect.Close()

			roominfo := socket.RoomInfo { Roomname : packetLogout.Payload.Roomname }
			chatMessage := socket.ChatMessage {From : userinfo , Stamp : timeUnix , Text : "exit room" , Style : "1" , Roominfo : roominfo}
			logoutBrocast := socket.Cmd_b_player_exit_room_struct { Base_B : socket.Base_B {Cmd : socket.CMD_B_PLAYER_EXIT_ROOM , Stamp : timeUnix} , Payload : chatMessage}
			logoutBrocastJson, _ := json.Marshal(logoutBrocast)

			broadcast(packetLogout.Payload.Roomname, mtype , logoutBrocastJson)

			break;
		case socket.CMD_C_PLAYER_EXIT_ROOM:
			
			var packetExitRoom socket.Cmd_c_player_exit_room_struct

			if err := json.Unmarshal([]byte(msg), &packetExitRoom); err != nil {
				panic(err)
			}

			if checkinroom(packetExitRoom.Payload.Roomname , clients[connect] ){
				
				roominfo := socket.RoomInfo { Roomname : packetExitRoom.Payload.Roomname }
				chatMessage := socket.ChatMessage {From : clients[connect].user , Stamp : timeUnix , Text : "exit room" , Style : "1" , Roominfo : roominfo}
				exitroomBrocast := socket.Cmd_b_player_exit_room_struct { Base_B : socket.Base_B {Cmd : socket.CMD_B_PLAYER_EXIT_ROOM , Stamp : timeUnix} , Payload : chatMessage}
				exitroomBrocastJson, err := json.Marshal(exitroomBrocast)
				if err != nil {
					log.Println("json err:", err)
				}
	
				broadcast(packetExitRoom.Payload.Roomname, mtype , exitroomBrocastJson)

				delete(clients[connect].room,packetExitRoom.Payload.Roomname)
				delete(rooms[packetExitRoom.Payload.Roomname], connect)
				exitroom := socket.Cmd_r_player_exit_room_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_EXIT_ROOM , Idem : packetExitRoom.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}}}
				exitroomJson, _ := json.Marshal(exitroom)
				connect.WriteMessage(mtype , exitroomJson)

				if(len(rooms[packetExitRoom.Payload.Roomname]) == 0){
					delete(rooms,packetExitRoom.Payload.Roomname)
					delete(roomsInfo,packetExitRoom.Payload.Roomname)
				}
				
				log.Printf("roomsInfo : %+v\n", roomsInfo)
				
				log.Printf("rooms : %+v\n", rooms)
				
				log.Printf("clients : %+v\n", clients)

			}else{
				exitroom := socket.Cmd_r_player_exit_room_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_EXIT_ROOM , Idem : packetExitRoom.Idem , Stamp : timeUnix, Result : "err", Exp :  socket.Exception{Code : "101", Message : "NOT_IN_ROOM"}}}
				exitroomJson, _ := json.Marshal(exitroom)
				connect.WriteMessage(mtype , exitroomJson)
				break;
			}

			break;
		case socket.CMD_C_PLAYER_ENTER_ROOM:

			var packetEnterRoom socket.Cmd_c_player_enter_room_struct

			if err := json.Unmarshal([]byte(msg), &packetEnterRoom); err != nil {
				panic(err)
			}

			if checkinroom(packetEnterRoom.Payload.Roomname , clients[connect] ){
				enterroom := socket.Cmd_r_player_enter_room_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_ENTER_ROOM , Idem : packetEnterRoom.Idem , Stamp : timeUnix, Result : "err", Exp :  socket.Exception{Code : "101", Message : "IN_ROOM"}}}
				enterroomJson, _ := json.Marshal(enterroom)
				connect.WriteMessage(mtype , enterroomJson)
				break;
			}
			client := clients[connect]
			client.room[packetEnterRoom.Payload.Roomname] = socket.RoomInfo { Roomname : packetEnterRoom.Payload.Roomname }
			if(rooms[packetEnterRoom.Payload.Roomname] == nil){
				var adduser = make(map[*websocket.Conn]Client)
				adduser[connect] = client
				rooms[packetEnterRoom.Payload.Roomname] = adduser
				roomsInfo[packetEnterRoom.Payload.Roomname] = socket.RoomInfo { Roomname : packetEnterRoom.Payload.Roomname }
			}else{
				rooms[packetEnterRoom.Payload.Roomname][connect] = client
			}
			

			enterroom := socket.Cmd_r_player_enter_room_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_ENTER_ROOM , Idem : packetEnterRoom.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}}}
			enterroomJson, _ := json.Marshal(enterroom)
			connect.WriteMessage(mtype , enterroomJson)

			roominfo := socket.RoomInfo { Roomname : packetEnterRoom.Payload.Roomname }
			chatMessage := socket.ChatMessage {From : clients[connect].user , Stamp : timeUnix , Text : "enter room" , Style : "1" , Roominfo : roominfo}
			enterroomBrocast := socket.Cmd_b_player_enter_room_struct { Base_B : socket.Base_B {Cmd : socket.CMD_B_PLAYER_ENTER_ROOM , Stamp : timeUnix} , Payload : chatMessage}
			enterroomBrocastJson, _ := json.Marshal(enterroomBrocast)

			broadcast(packetEnterRoom.Payload.Roomname, mtype , enterroomBrocastJson)

			break;
		case socket.CMD_C_PLAYER_SEND_MSG:
			var packetSendMsg socket.Cmd_c_player_send_msg_struct

			if err := json.Unmarshal([]byte(msg), &packetSendMsg); err != nil {
				panic(err)
			}
			
			if !checkinroom(packetSendMsg.Payload.Roominfo.Roomname , clients[connect] ){
				SendMsg := socket.Cmd_r_player_send_msg_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_SEND_MSG , Idem : packetSendMsg.Idem , Stamp : timeUnix, Result : "err", Exp :  socket.Exception{Code : "101", Message : "NOT_IN_ROOM"}}}
				SendMsgJson, _ := json.Marshal(SendMsg)
				connect.WriteMessage(mtype , SendMsgJson)
				break;
			}
			// log.Printf("packetSendMsg : %+v\n", packetSendMsg)

			SendMsg := socket.Cmd_r_player_send_msg_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_PLAYER_SEND_MSG , Idem : packetSendMsg.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}}}
			SendMsgJson, _ := json.Marshal(SendMsg)

			connect.WriteMessage(mtype , SendMsgJson)

			roominfo := socket.RoomInfo { Roomname : packetSendMsg.Payload.Roominfo.Roomname }
			chatMessage := socket.ChatMessage {From : clients[connect].user , Stamp : timeUnix , Text : packetSendMsg.Payload.Text , Style : packetSendMsg.Payload.Style , Roominfo : roominfo}
			sendMsgBrocast := socket.Cmd_b_player_speak_struct { Base_B : socket.Base_B {Cmd : socket.CMD_B_PLAYER_SPEAK , Stamp : timeUnix} , Payload : chatMessage}
			sendMsgBrocastJson, _ := json.Marshal(sendMsgBrocast)
			
			broadcast(packetSendMsg.Payload.Roominfo.Roomname, mtype , sendMsgBrocastJson)
			break;
		case socket.CMD_C_GET_MEMBER_LIST:
			
			var packetRoomList socket.Cmd_c_get_member_list
			
			if err := json.Unmarshal([]byte(msg), &packetRoomList); err != nil {
				panic(err)
			}
			
			memberList := make([]socket.User, 0, len(rooms[packetRoomList.Payload.Roomname]))
			for _ , v := range rooms[packetRoomList.Payload.Roomname] {
				memberList = append(memberList, v.user)
			}
			
			SendMemberList := socket.Cmd_r_get_member_list_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_GET_MEMBER_LIST , Idem : packetRoomList.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}} , Payload : memberList}
			SendMemberListJson, _ := json.Marshal(SendMemberList)

			connect.WriteMessage(mtype , SendMemberListJson)

			break;
		case socket.CMD_C_GET_ROOM_LIST:
			var packetRoomList socket.Cmd_c_get_room_list
			
			if err := json.Unmarshal([]byte(msg), &packetRoomList); err != nil {
				panic(err)
			}
			
			roomList := make([]socket.RoomInfo, 0, len(roomsInfo))
			for _ , roominfo := range roomsInfo {
				roomList = append(roomList, roominfo)
			}
			// log.Printf("roomList : %+v\n", roomList)
			
			SendRoomList := socket.Cmd_r_get_room_list_struct { Base_R : socket.Base_R {Cmd : socket.CMD_R_GET_ROOM_LIST , Idem : packetRoomList.Idem , Stamp : timeUnix, Result : "ok", Exp :  socket.Exception{Code : "", Message : ""}} , Payload : roomList}
			SendRoomListJson, _ := json.Marshal(SendRoomList)
			connect.WriteMessage(mtype , SendRoomListJson)
			break;
		default:


	}
	// Write message back to the client message owner with:
	//connect.WriteMessage(mtype , msg.Body)
	// Write message to all except this client with:
	//broadcast(hubName , mtype , msg)

 return nil
}

func broadcast(roomName string , mtype int , msg []byte) {

	//log.Printf("roomName : %s\n", roomName)
	targetroom := rooms[roomName]
	//log.Printf("targetroom : %+v\n", targetroom)
	
    for clientConnect := range targetroom {
		clientConnect.WriteMessage(mtype, msg)
	}
}

func checkinroom(roomName string , client Client) bool{
	for _, clientRoom := range client.room {
        if roomName == clientRoom.Roomname {
            return true
        }
    }
	return false
}