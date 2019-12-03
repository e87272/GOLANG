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
	proclamationlist, ok := Proclamationlist[roomUuid]
	Mutexproclamationlist.Unlock()
	if !ok {
		Queryproclamation(roomUuid)
		Mutexproclamationlist.Lock()
		proclamationlist = Proclamationlist[roomUuid]
		Mutexproclamationlist.Unlock()
	}
	return proclamationlist
}

func Queryproclamation(roomUuid string) {

	var proclamationlist = map[string]socket.Proclamation{}
	rows, err := database.Query("select proclamationUuid,type,clientOrder,appType,title,content,style,url from proclamation where roomUuid = ?", roomUuid)
	if err != nil {
		Essyserrorlog("COMMON_QUERYPROCLAMATION_SELECT_DB_ERROR", "", err)
		Mutexproclamationlist.Lock()
		Proclamationlist[roomUuid] = proclamationlist
		Mutexproclamationlist.Unlock()
		return
	}

	var proclamationUuid string
	var proclamationType string
	var order string
	var appType string
	var title string
	var content string
	var style string
	var url string

	Mutexproclamationlist.Lock()
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
	Proclamationlist[roomUuid] = proclamationlist

	Mutexproclamationlist.Unlock()

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	proclamation := socket.Cmd_b_proclamation{Base_B: socket.Base_B{Cmd: socket.CMD_B_PROCLAMATION, Stamp: timeUnix}, Payload: proclamationlist}
	proclamationJson, _ := json.Marshal(proclamation)
	Broadcast(roomUuid, proclamationJson)
	return
}
