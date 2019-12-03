package command

import (
	"encoding/json"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/websocket"

	"../common"
	"../database"
	"../socket"
)

func Tokenchange(connect *websocket.Conn, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendTokenChange := socket.Cmd_r_token_change{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_TOKEN_CHANGE,
		Stamp: timeUnix,
	}}
	userPlatform, _ := common.Clientsuserplatformread(loginUuid)
	userUuid := userPlatform.Useruuid

	var packetToken socket.Cmd_c_token_change

	if err := json.Unmarshal([]byte(msg), &packetToken); err != nil {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_JSON_ERROR", userUuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connect, sendTokenChangeJson)
		return err
	}
	if packetToken.Payload.Platformuuid == "" || packetToken.Payload.Platform == "" {
		var user socket.User
		user.Userplatform.Useruuid = loginUuid
		user.Userplatform.Platformuuid = ""
		user.Userplatform.Platform = ""
		user.Globalrole = ""
		sendTokenChange.Payload = user
		sendTokenChange.Result = "ok"
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)

		common.Sendmessage(connect, sendTokenChangeJson)

		var client = common.Client{Room: make(map[string]socket.Roomcore), Conn: connect, Userplatform: user.Userplatform, Sidetext: make(map[string]common.Sidetextplatform)}
		common.Clientsinsert(loginUuid, client)

		var userConnect = make(map[string]*websocket.Conn)
		userConnect[loginUuid] = connect
		common.Clientsconnectinsert(user.Userplatform.Useruuid, userConnect)
		return nil
	}
	user, err := userCheck(packetToken.Payload.Platformuuid, packetToken.Payload.Platform, packetToken.Payload.Token)
	// log.Printf("Tokenchange user : %+v\n", user)
	if err != nil {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_DB_ERROR", userUuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connect, sendTokenChangeJson)
		return nil
	}
	// log.Printf("packetToken : %+v\n", packetToken)

	sideTextMap, err := common.Querysidetextmap(user.Userplatform.Useruuid)
	if err != nil {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_SIDETEXTMAP_ERROR", userUuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connect, sendTokenChangeJson)
		return nil
	}

	sendTokenChange.Result = "ok"
	sendTokenChange.Payload = user
	sendTokenChangeJson, _ := json.Marshal(sendTokenChange)

	common.Sendmessage(connect, sendTokenChangeJson)

	// log.Printf("Tokenchange sideTextMap : %+v\n", sideTextMap)
	var client = common.Client{Room: make(map[string]socket.Roomcore), Conn: connect, Userplatform: user.Userplatform, Sidetext: sideTextMap}

	// log.Printf("Tokenchange client : %+v\n", client)
	common.Clientsinsert(loginUuid, client)
	// log.Printf("Tokenchange Clients : %+v\n", common.Clients)
	common.Usersinfoinsert(user.Userplatform.Useruuid, user)
	common.Setredisuserinfo(user.Userplatform.Useruuid, user)

	_, ok := common.Clientsconnectread(user.Userplatform.Useruuid)
	if !ok {
		var userConnect = make(map[string]*websocket.Conn)
		userConnect[loginUuid] = connect
		common.Clientsconnectinsert(user.Userplatform.Useruuid, userConnect)
	} else {
		common.Clientsconnectloginuuidinsert(user.Userplatform.Useruuid, loginUuid, connect)
	}
	return nil
}

func userCheck(platformUuid string, platform string, token string) (socket.User, error) {

	user := socket.User{}

	isUser, err := common.Checkplatformuser(platform, platformUuid, token)
	if err != nil || !isUser {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_CHECKPLATFORMUSER_ERROR", platform+"-"+platformUuid, err)
		// return user, err
	}

	row := database.QueryRow("SELECT uuid,platformUuid,platform,globalRole FROM users WHERE platformUuid = ? AND platform = ?", platformUuid, platform)
	err = row.Scan(&user.Userplatform.Useruuid, &user.Userplatform.Platformuuid, &user.Userplatform.Platform, &user.Globalrole)
	// log.Printf("user : %+v\n", user)

	if err == database.ErrNoRows {
		uuid := common.Getid().Hexstring()
		_, err := database.Exec(
			"INSERT INTO users (uuid, platformUuid, platform, globalRole) VALUES (?, ?, ?, ?)",
			uuid,
			platformUuid,
			platform,
			"",
		)
		// log.Printf("INSERT INTO users (uuid, platformUuid, platform) VALUES (%s, %s ,%s)\n", uuid, platformUuid, platform)
		if err != nil {
			common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_USER_INSERT_ERROR", platform+"-"+platformUuid, err)
			return user, err
		}
		user.Userplatform.Useruuid = uuid
		user.Userplatform.Platformuuid = platformUuid
		user.Userplatform.Platform = platform
		user.Globalrole = ""
		return user, nil
	} else if err != nil {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_USER_SELECT_ERROR", platform+"-"+platformUuid, err)
		return user, err
	}

	rows, err := database.Query("select roomUuid from vipGroupUserList where userUuid = ?",
		user.Userplatform.Useruuid,
	)

	if err != nil {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_SELECT_VIPGROUPUSERLIST_ERROR", user.Userplatform.Useruuid, err)
		return user, err
	}

	for rows.Next() {
		var roomUuid string
		rows.Scan(&roomUuid)
		if user.Vipgroup == "" {
			user.Vipgroup = roomUuid
		} else {
			user.Vipgroup = user.Vipgroup + "," + roomUuid
		}
	}
	rows.Close()

	// log.Printf("Tokenchange user.Vipgroup : %+v\n", user.Vipgroup)

	rows, err = database.Query("select roomUuid from privateGroupUserList where userUuid = ?",
		user.Userplatform.Useruuid,
	)

	if err != nil {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_SELECT_PRIVATEGROUPUSERLIST_ERROR", user.Userplatform.Useruuid, err)
		return user, err
	}

	for rows.Next() {
		var roomUuid string
		rows.Scan(&roomUuid)
		if user.Privategroup == "" {
			user.Privategroup = roomUuid
		} else {
			user.Privategroup = user.Privategroup + "," + roomUuid
		}
	}
	rows.Close()

	// log.Printf("Tokenchange user.Privategroup : %+v\n", user.Privategroup)

	return user, nil
}
