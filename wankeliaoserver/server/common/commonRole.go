package common

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../database"
	"../socket"
)

func Queryblocklist() {

	rows, err := database.Query("select blockUuid,blockUserUuid,blocktarget,timeStamp from chatBlock where timeStamp >= ?", time.Now().UnixNano()/int64(time.Millisecond))

	// log.Printf("select blockUuid,blockUserUuid,blocktarget,timeStamp from chatBlock where timeStamp >= %+v\n", time.Now().UnixNano()/int64(time.Millisecond))

	if err != nil {
		Essyserrorlog("COMMON_QUERYBLOCKLIST_ERROR", "", err)
		return
	}
	var blockUuid string
	var blockUserUuid string
	var blocktarget string
	var timeStamp int64

	Mutexblockchatlist.Lock()
	BlockchatList = make(map[string]map[string]int64)
	for rows.Next() {
		rows.Scan(&blockUuid, &blockUserUuid, &blocktarget, &timeStamp)
		stampMap, ok := BlockchatList[blockUserUuid]
		if !ok {
			stampMap = make(map[string]int64)
			BlockchatList[blockUserUuid] = stampMap
		}
		stampMap[blocktarget] = timeStamp
		// log.Printf("blockUuid : %+v\n", blockUuid)
		// log.Printf("blockUserUuid : %+v\n", blockUserUuid)
		// log.Printf("blocktarget : %+v\n", blocktarget)
		// log.Printf("timeStamp : %+v\n", timeStamp)
		// log.Printf("stamp : %+v\n", stampMap)
		// log.Printf("BlockchatList : %+v\n", BlockchatList)
	}
	Mutexblockchatlist.Unlock()
	rows.Close()
}

func Checkblock(roomUuid string, uuid string) bool {

	Mutexblockchatlist.Lock()
	defer Mutexblockchatlist.Unlock()

	// log.Printf("BlockchatList : %+v\n", BlockchatList)
	blockChatListUuid, ok := BlockchatList[uuid]
	if ok {
		stamp, ok := blockChatListUuid[roomUuid]
		if ok {
			if stamp > time.Now().UnixNano()/int64(time.Millisecond) {
				return true
			}
		}

		stamp, ok = blockChatListUuid[uuid]
		if ok {
			if stamp > time.Now().UnixNano()/int64(time.Millisecond) {
				return true
			}
		}
	}
	return false
}

func Checkadmin(roomUuid string, userUuid string, functionName string) bool {
	role := Functionmanagementread(functionName)
	roleAry := strings.Split(role, ",")
	user, ok := Usersinforead(userUuid)
	if !ok {
		return false
	}
	globalUserRoleAry := strings.Split(user.Globalrole, ",")
	roleMap := map[string]bool{}
	for _, val := range roleAry {
		roleMap[val] = true
	}
	for _, val := range globalUserRoleAry {
		_, ok = roleMap[val]
		if ok {
			return true
		}
	}

	roomInfo, ok := Roomsinforead(roomUuid)
	if !ok {
		return false
	}
	var adminSet map[string]string
	err := json.Unmarshal([]byte(roomInfo.Adminset), &adminSet)
	if err != nil {
		return false
	}
	userRole, ok := adminSet[userUuid]
	if !ok {
		return false
	}
	userRoleAry := strings.Split(userRole, ",")
	roleMap = map[string]bool{}
	for _, val := range roleAry {
		roleMap[val] = true
	}
	for _, val := range userRoleAry {
		_, ok = roleMap[val]
		if ok {
			return true
		}
	}
	return false
}

func Checktargetadmin(roomUuid string, targetUuid string, targetGlobalRole string, functionName string) bool {
	role := Functionmanagementread(functionName)
	roleAry := strings.Split(role, ",")

	globalUserRoleAry := strings.Split(targetGlobalRole, ",")
	roleMap := map[string]bool{}
	for _, val := range roleAry {
		roleMap[val] = true
	}
	for _, val := range globalUserRoleAry {
		_, ok := roleMap[val]
		if ok {
			return true
		}
	}

	roomInfo, _ := Roomsinforead(roomUuid)
	var adminSet map[string]string
	err := json.Unmarshal([]byte(roomInfo.Adminset), &adminSet)
	if err != nil {
		return false
	}
	userRole, ok := adminSet[targetUuid]
	if !ok {
		return false
	}
	userRoleAry := strings.Split(userRole, ",")
	roleMap = map[string]bool{}
	for _, val := range roleAry {
		roleMap[val] = true
	}
	for _, val := range userRoleAry {
		_, ok = roleMap[val]
		if ok {
			return true
		}
	}
	return false
}

func Queryfunctionmanagement() {
	Mutexfunctionmanagement.Lock()
	defer Mutexfunctionmanagement.Unlock()
	rows, err := database.Query(
		"select functionName,functionRule from functionManagement",
	)

	if err != nil {
		Essyserrorlog("COMMON_QUERYFUNCTIONMANAGEMENT_ERROR", "", err)
	}
	var functionName string
	var functionRule string

	Functionmanagement = make(map[string]string)

	for rows.Next() {

		rows.Scan(&functionName, &functionRule)
		Functionmanagement[functionName] = functionRule

	}
	rows.Close()

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendFunctionManagement := socket.Cmd_b_func_management{Base_B: socket.Base_B{Cmd: socket.CMD_B_FUNC_MANAGEMENT, Stamp: timeUnix}, Payload: Functionmanagement}
	sendFunctionManagementJson, _ := json.Marshal(sendFunctionManagement)

	Mutexclients.Lock()
	defer Mutexclients.Unlock()
	for _, client := range Clients {
		Sendmessage(client.Conn, sendFunctionManagementJson)
	}
}

func Functionmanagementread(functionName string) string {
	Mutexfunctionmanagement.Lock()
	defer Mutexfunctionmanagement.Unlock()

	return Functionmanagement[functionName]
}

func Pubsudoresult(shellTarget string, userPlatform socket.Userplatform, Cmd []string, targetUserPlatform socket.Userplatform, functionName string) {

	sudoShellData := struct {
		Sudoresult socket.Sudoresult
		FuncName   string
	}{}
	sudoShellData.Sudoresult = socket.Sudoresult{Shelltarget: shellTarget, Userplatform: userPlatform, Cmd: Cmd, Targetuserplatform: targetUserPlatform}
	sudoShellData.FuncName = functionName
	sudoShellDataJson, _ := json.Marshal(sudoShellData)
	Redispubdata("sudoresult", string(sudoShellDataJson))
}

func Checkword(s string) bool {
	return regexp.MustCompile(`^\w+$`).MatchString(s)
}
