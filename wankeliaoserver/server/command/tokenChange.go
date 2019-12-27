package command

import (
	"encoding/json"

	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"../common"
	"../database"
	"../socket"
)

func Tokenchange(connCore common.Conncore, msg []byte, loginUuid string) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sendTokenChange := socket.Cmd_r_token_change{Base_R: socket.Base_R{
		Cmd:   socket.CMD_R_TOKEN_CHANGE,
		Stamp: timeUnix,
	}}

	var packetToken socket.Cmd_c_token_change

	if err := json.Unmarshal([]byte(msg), &packetToken); err != nil {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_JSON_ERROR", "loginUuid : "+loginUuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connCore, sendTokenChangeJson)
		return err
	}
	if packetToken.Payload.Platformuuid == "" || packetToken.Payload.Platform == "" || packetToken.Payload.Token == "" {
		var user socket.User
		user.Userplatform.Useruuid = loginUuid
		user.Userplatform.Platformuuid = ""
		user.Userplatform.Platform = ""
		user.Globalrole = ""
		sendTokenChange.Payload = user
		sendTokenChange.Result = "ok"
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)

		common.Sendmessage(connCore, sendTokenChangeJson)
		var client = common.Client{Room: make(map[string]socket.Roomcore), Conncore: connCore, Userplatform: user.Userplatform, Sidetext: make(map[string]common.Sidetextplatform)}
		common.Clientsinsert(loginUuid, client)

		var userConnect = make(map[string]common.Conncore)
		userConnect[loginUuid] = connCore
		common.Clientsconnectinsert(user.Userplatform.Useruuid, userConnect)
		return nil
	}
	user, ok, err := userCheck(packetToken.Payload.Platformuuid, packetToken.Payload.Platform, packetToken.Payload.Token)
	// log.Printf("Tokenchange user : %+v\n", user)
	// log.Printf("Tokenchange ok : %+v\n", ok)
	// log.Printf("Tokenchange err : %+v\n", err)
	if !ok {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_USERCHECK_ERROR", "loginUuid : "+loginUuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connCore, sendTokenChangeJson)
		connCore.Conn.Close()
		return nil
	}
	// log.Printf("packetToken : %+v\n", packetToken)

	sideTextMap, err := common.Querysidetextmap(user.Userplatform.Useruuid)
	if err != nil {
		sendTokenChange.Base_R.Result = "err"
		sendTokenChange.Base_R.Exp = common.Exception("COMMAND_TOKENCHANGE_SIDETEXTMAP_ERROR", user.Userplatform.Useruuid, err)
		sendTokenChangeJson, _ := json.Marshal(sendTokenChange)
		common.Sendmessage(connCore, sendTokenChangeJson)
		connCore.Conn.Close()
		return nil
	}

	sendTokenChange.Result = "ok"
	sendTokenChange.Payload = user
	sendTokenChangeJson, _ := json.Marshal(sendTokenChange)

	common.Sendmessage(connCore, sendTokenChangeJson)

	// log.Printf("Tokenchange sideTextMap : %+v\n", sideTextMap)
	var client = common.Client{Room: make(map[string]socket.Roomcore), Conncore: connCore, Userplatform: user.Userplatform, Sidetext: sideTextMap}

	// log.Printf("Tokenchange client : %+v\n", client)
	common.Clientsinsert(loginUuid, client)
	// log.Printf("Tokenchange Clients : %+v\n", common.Clients)
	common.Usersinfoinsert(user.Userplatform.Useruuid, user)
	common.Setredisuserinfo(user.Userplatform.Useruuid, user)

	_, ok = common.Clientsconnectread(user.Userplatform.Useruuid)
	if !ok {
		var userConnect = make(map[string]common.Conncore)
		userConnect[loginUuid] = connCore
		common.Clientsconnectinsert(user.Userplatform.Useruuid, userConnect)
	} else {
		common.Clientsconnectloginuuidinsert(user.Userplatform.Useruuid, loginUuid, connCore)
	}
	return nil
}

func userCheck(platformUuid string, platform string, token string) (socket.User, bool, error) {

	user := socket.User{}
	if token != os.Getenv("tokenByTest") {
		isUser, err := common.Checkplatformuser(platform, platformUuid, token)
		if err != nil || !isUser {
			common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_CHECKPLATFORMUSER_ERROR", platform+"-"+platformUuid, err)
			return user, false, err
		}
	}
	var toBeCharge string //待轉化為時間戳的字串 注意 這裡的小時和分鐘還要秒必須寫 因為是跟著模板走的 修改模板的話也可以不寫
	var createTime int64
	row := database.QueryRow("SELECT uuid,platformUuid,platform,globalRole,created_at FROM users WHERE platformUuid = ? AND platform = ?", platformUuid, platform)
	err := row.Scan(&user.Userplatform.Useruuid, &user.Userplatform.Platformuuid, &user.Userplatform.Platform, &user.Globalrole, &toBeCharge)

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
			return user, false, err
		}
		user.Userplatform.Useruuid = uuid
		user.Userplatform.Platformuuid = platformUuid
		user.Userplatform.Platform = platform
		user.Globalrole = ""
	} else if err != nil {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_USER_SELECT_ERROR", platform+"-"+platformUuid, err)
		return user, false, err
	}

	if toBeCharge != "" {
		timeLayout := "2006-01-02 15:04:05"                             //轉化所需模板
		loc, _ := time.LoadLocation("Local")                            //重要：獲取時區
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在對應時區轉化為time.time型別
		createTime = theTime.UnixNano() / int64(time.Millisecond)       //轉化為時間戳 型別是int64
	} else {
		createTime = time.Now().UnixNano() / int64(time.Millisecond)
	}

	// log.Printf("userCheck Now : %+v\n", time.Now().UnixNano()/int64(time.Millisecond))
	// log.Printf("userCheck createTime : %+v\n", createTime)
	// log.Printf("userCheck common.Newusercdtime : %+v\n", common.Newusercdtime)
	if time.Now().UnixNano()/int64(time.Millisecond)-createTime < common.Newusercdtime {
		common.Blocknewuserlistinsert(user.Userplatform.Useruuid, createTime+common.Newusercdtime)
	}

	rows, err := database.Query("select roomUuid from vipGroupUserList where userUuid = ?",
		user.Userplatform.Useruuid,
	)

	if err != nil {
		common.Essyserrorlog("COMMAND_TOKENCHANGE_USERCHECK_SELECT_VIPGROUPUSERLIST_ERROR", user.Userplatform.Useruuid, err)
		return user, false, err
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
		return user, false, err
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

	return user, true, nil
}
