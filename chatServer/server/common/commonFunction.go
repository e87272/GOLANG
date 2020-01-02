package common

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	_ "net/http/pprof"

	"../database"
	"../socket"
)

func Myiplastdigit() string {

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	// log.Printf("localAddr : %+v\n", localAddr)

	idx := strings.LastIndex(localAddr, ":")

	myIpLastDigit := localAddr[0:idx]

	return myIpLastDigit
}

func Proclamationsearch(roomUuid string) map[string]socket.Proclamation {
	Mutexproclamationlist.Lock()
	defer Mutexproclamationlist.Unlock()
	proclamationlist, ok := Proclamationlist[roomUuid]
	// log.Printf("Proclamationsearch Proclamationlist : %+v\n", Proclamationlist)
	if !ok {
		proclamationlist = Queryproclamation(roomUuid)
		Proclamationlist[roomUuid] = proclamationlist
		// log.Printf("Proclamationsearch !ok Proclamationlist : %+v\n", Proclamationlist)
	}
	return proclamationlist
}

func Queryproclamation(roomUuid string) map[string]socket.Proclamation {

	var proclamationlist = map[string]socket.Proclamation{}
	rows, err := database.Query("select proclamationUuid,type,clientOrder,appType,title,content,style,url from proclamation where roomUuid = ?", roomUuid)
	if err != nil {
		Essyserrorlog("COMMON_QUERYPROCLAMATION_SELECT_DB_ERROR", "", err)
		return proclamationlist
	}

	var proclamationUuid string
	var proclamationType string
	var order string
	var appType string
	var title string
	var content string
	var style string
	var url string

	for rows.Next() {
		rows.Scan(&proclamationUuid, &proclamationType, &order, &appType, &title, &content, &style, &url)
		var proclamation = socket.Proclamation{}

		proclamation.Proclamationuuid = proclamationUuid
		proclamation.Roomuuid = roomUuid
		proclamation.Type = proclamationType
		proclamation.Order = order
		proclamation.Apptype = appType
		proclamation.Title = title
		proclamation.Content = content
		proclamation.Style = style
		proclamation.Url = url

		proclamationlist[proclamationUuid] = proclamation
	}
	rows.Close()
	packetStamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeUnix := strconv.FormatInt(packetStamp, 10)
	proclamation := socket.Cmd_b_proclamation{Base_B: socket.Base_B{Cmd: socket.CMD_B_PROCLAMATION, Stamp: timeUnix}, Payload: proclamationlist}
	proclamationJson, _ := json.Marshal(proclamation)
	Broadcast(roomUuid, proclamationJson, packetStamp)
	return proclamationlist
}
