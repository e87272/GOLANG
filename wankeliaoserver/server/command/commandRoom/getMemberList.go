package commandRoom

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../../common"
	"../../database"
	"../../socket"
)

func Getmemberlist(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendMemberList := socket.Cmd_r_get_member_list{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_GET_MEMBER_LIST,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetGetMemberList socket.Cmd_c_get_member_list
	err := json.Unmarshal([]byte(msg), &packetGetMemberList)
	if err != nil {
		sendMemberList.Base_R.Result = "err"
		sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETMEMBERLIST_JSON_ERROR", userUuid, err)
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return err
	}
	sendMemberList.Base_R.Idem = packetGetMemberList.Base_C.Idem

	if loginUuid == userUuid {
		//block處理
		sendMemberList.Base_R.Result = "err"
		sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETFRIENDLIST_GUEST", userUuid, nil)
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return nil
	}

	roomCore, ok := common.Clientsroomread(loginUuid, packetGetMemberList.Payload.Roomuuid)
	if !ok {
		//block處理
		sendMemberList.Base_R.Result = "err"
		sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETMEMBERLIST_ROOM_UUID_ERROR", userUuid, nil)
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return nil
	}
	if roomCore.Roomtype != packetGetMemberList.Payload.Roomtype {
		sendMemberList.Base_R.Result = "err"
		sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETMEMBERLIST_ROOM_TYPE_ERROR", userUuid, nil)
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return nil
	}

	if roomCore.Roomtype == "liveGroup" {
		sendMemberList.Base_R.Result = "err"
		sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETMEMBERLIST_NOT_LIVEGROUP_LIST", userUuid, nil)
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return nil
		// return getLiveGroupMemberList(connect, packetGetMemberList, sendMemberList)
	} else {
		memberList,err := getGroupMemberList(connect, packetGetMemberList, sendMemberList)
		if err != nil{
			sendMemberList.Base_R.Result = "err"
			sendMemberList.Base_R.Exp = common.Exception("COMMAND_GETMEMBERLIST_ERROR", userUuid, err)
			sendMemberListJson, _ := json.Marshal(sendMemberList)
			common.Sendmessage(connect, sendMemberListJson)
			return nil
		}

		sendMemberList.Base_R.Result = "ok"
		sendMemberList.Payload = memberList
		sendMemberListJson, _ := json.Marshal(sendMemberList)
		common.Sendmessage(connect, sendMemberListJson)
		return nil
	}

}

func getGroupMemberList(connect *websocket.Conn, packetGetMemberList socket.Cmd_c_get_member_list, sendMemberList socket.Cmd_r_get_member_list) ([]socket.Userplatform,error) {

	var memberList []socket.Userplatform

	userListName := packetGetMemberList.Payload.Roomtype + "UserList"
	rows, err := database.Query("SELECT users.uuid, users.platform, users.platformUuid FROM users RIGHT JOIN "+userListName+" ON users.uuid="+userListName+".userUuid WHERE "+userListName+".roomUuid = ?",
		packetGetMemberList.Payload.Roomuuid,
	)

	// log.Printf("SELECT users.uuid, users.platform, users.platformUuid FROM users RIGHT JOIN " + userListName + " ON users.uuid=" + userListName + ".userUuid WHERE " + userListName + ".roomUuid = " + packetGetMemberList.Payload.Roomuuid + "\n")

	if err != nil {
		return memberList,err
	}

	for rows.Next() {
		var userUuid string
		var platform string
		var platformUuid string
		rows.Scan(&userUuid, &platform, &platformUuid)
		member := socket.Userplatform{}
		member.Useruuid = userUuid
		member.Platform = platform
		member.Platformuuid = platformUuid
		memberList = append(memberList, member)
		// log.Printf("getGroupMemberList member : %+v\n", member)
	}
	rows.Close()

	common.Setredismembercount(packetGetMemberList.Payload.Roomuuid, len(memberList))
	return memberList,nil
}

func getLiveGroupMemberList(connect *websocket.Conn, packetGetMemberList socket.Cmd_c_get_member_list, sendMemberList socket.Cmd_r_get_member_list) error {

	resp, err := common.PostApiForm("http://127.0.0.1"+os.Getenv("imServerPort")+"/emit/roomUserSearch", url.Values{"roomUuid": {packetGetMemberList.Payload.Roomuuid}})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var memberListJson string
	err = json.Unmarshal(body, &memberListJson)
	if err != nil {
		return err
	}

	var memberList []socket.Userplatform
	err = json.Unmarshal([]byte(memberListJson), &memberList)
	if err != nil {
		return err
	}

	sendMemberList.Base_R.Result = "ok"
	sendMemberList.Payload = memberList
	sendMemberListJson, _ := json.Marshal(sendMemberList)
	common.Sendmessage(connect, sendMemberListJson)

	return nil
}
