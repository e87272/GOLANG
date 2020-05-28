package common

import (
	"encoding/json"

	"server/database"
	"server/socket"
)

func Sidetextsend(userUuid string, targetUuid string, msg []byte) {
	//全打預防多登入狀況

	// log.Printf("Clientsconnect : %+v\n", Clientsconnect)
	Mutexclientsconnect.Lock()
	// log.Printf("Mutexclientsconnect :Sidetextsend\n")
	defer func() {
		Mutexclientsconnect.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexclientsconnect :SidetextsendUNLock\n")
	}()

	// log.Printf("userUuid : %s\n", userUuid)
	user, ok := Clientsconnect[userUuid]
	// log.Printf("user : %+v\n", user)
	if ok {

		// log.Printf("Sidetextsend range user : %+v\n", user)
		for loginUuid, connect := range user {
			_, ok := Clientssidetextuserread(loginUuid, targetUuid)
			if !ok {
				sideTextUserPlatform, err := Querysidetextuuid(userUuid, targetUuid)
				if err != nil {
					Essyserrorlog("COMMON_SIDETEXTSEND_SIDETEXT_UUID_ERROR", userUuid, err)
					continue
				}
				Clientssidetextuserinsert(loginUuid, targetUuid, sideTextUserPlatform)
			}
			Sendmessage(connect, msg)
		}
	}

	// log.Printf("targetUuid : %s\n", targetUuid)
	targetUser, ok := Clientsconnect[targetUuid]
	// log.Printf("targetUser : %+v\n", targetUser)
	if ok {

		var targetSideText socket.Cmd_b_player_speak
		json.Unmarshal(msg, &targetSideText)
		targetSideText.Payload.Chattarget = userUuid
		targetSideTextJson, _ := json.Marshal(targetSideText)

		// log.Printf("Sidetextsend range targetUser : %+v\n", targetUser)
		for loginUuid, connect := range targetUser {
			_, ok := Clientssidetextuserread(loginUuid, userUuid)
			if !ok {
				sideTextUserPlatform, err := Querysidetextuuid(targetUuid, userUuid)
				if err != nil {
					Essyserrorlog("COMMON_SIDETEXTSEND_SIDETEXT_UUID_ERROR", userUuid, err)
					continue
				}
				Clientssidetextuserinsert(loginUuid, userUuid, sideTextUserPlatform)
			}
			Sendmessage(connect, targetSideTextJson)
		}
	}
}

func Sidetextdeletesend(userUuid string, targetUuid string, msg []byte) {
	//全打預防多登入狀況

	// log.Printf("Clientsconnect : %+v\n", Clientsconnect)
	Mutexclientsconnect.Lock()
	// log.Printf("Mutexclientsconnect :Sidetextsend\n")
	defer func() {
		Mutexclientsconnect.Unlock() // 完成後記得 解鎖 解鎖 解鎖
		// log.Printf("Mutexclientsconnect :SidetextsendUNLock\n")
	}()

	// log.Printf("userUuid : %s\n", userUuid)
	user := Clientsconnect[userUuid]
	// log.Printf("user : %+v\n", user)

	for loginUuid, connect := range user {
		Sendmessage(connect, msg)
		Clientssidetextuserdelete(loginUuid, targetUuid)
	}

	// log.Printf("targetUuid : %s\n", targetUuid)
	targetUser := Clientsconnect[targetUuid]
	// log.Printf("targetUser : %+v\n", targetUser)
	var targetSideText socket.Cmd_b_sidetext_delete
	json.Unmarshal([]byte(msg), &targetSideText)
	targetSideText.Payload = userUuid
	targetSideTextJson, _ := json.Marshal(targetSideText)

	for loginUuid, connect := range targetUser {
		Sendmessage(connect, targetSideTextJson)
		Clientssidetextuserdelete(loginUuid, userUuid)
	}
}

func Updatesidetextlist(sideTextData string) {

	var data Redispubsubsidetextdata

	err := json.Unmarshal([]byte(sideTextData), &data)
	if err != nil {
		panic(err)
	}

	Sidetextsend(data.Useruuid, data.Targetuuid, []byte(data.Datajson))

}

func Deletesidetextlist(sideTextData string) {

	var data Redispubsubsidetextdata

	err := json.Unmarshal([]byte(sideTextData), &data)
	if err != nil {
		panic(err)
	}

	Sidetextdeletesend(data.Useruuid, data.Targetuuid, []byte(data.Datajson))

}

func Querysidetextmap(userUuid string) (map[string]Sidetextplatform, error) {

	// Essyslog("Querysidetextlist userUuid : "+userUuid, "", "")

	var sideTextMap = make(map[string]Sidetextplatform)

	sideTextUser := Sidetextplatform{}

	rows, err := database.Query("select sideTextUuid,forwardUuid,backwardUuid,forwardPlatformUuid,forwardPlatform,backwardPlatformUuid,backwardPlatform from sideText where forwardUuid = ? or backwardUuid = ?", userUuid, userUuid)

	// Essyslog("select sideTextUuid,forwardUuid,backwardUuid,forwardPlatformUuid,forwardPlatform,backwardPlatformUuid,backwardPlatform from sideText where forwardUuid = "+userUuid+" or backwardUuid = "+userUuid, "", "")
	if err != nil {
		Essyserrorlog("COMMON_QUERYSIDETEXTLIST_ERROR", userUuid, err)
		return sideTextMap, err
	}

	var sideTextUuid string
	var forwardUuid string
	var backwardUuid string
	var forwardPlatformUuid string
	var forwardPlatform string
	var backwardPlatformUuid string
	var backwardPlatform string
	for rows.Next() {
		rows.Scan(&sideTextUuid, &forwardUuid, &backwardUuid, &forwardPlatformUuid, &forwardPlatform, &backwardPlatformUuid, &backwardPlatform)

		sideTextUser.Sidetextuuid = sideTextUuid
		if forwardUuid == userUuid {
			sideTextUser.Userplatform = socket.Userplatform{Useruuid: backwardUuid, Platformuuid: backwardPlatformUuid, Platform: backwardPlatform}
			sideTextMap[backwardUuid] = sideTextUser
		} else {
			sideTextUser.Userplatform = socket.Userplatform{Useruuid: forwardUuid, Platformuuid: forwardPlatformUuid, Platform: forwardPlatform}
			sideTextMap[forwardUuid] = sideTextUser
		}
		// sideTextUserJson, _ := json.Marshal(sideTextUser)
		// Essyslog("Clientssidetextuserinsert rows : "+string(sideTextUserJson), "", userUuid)

	}
	rows.Close()

	return sideTextMap, err
}

func Querysidetextuuid(uuid string, targetUuid string) (Sidetextplatform, error) {

	// log.Printf("Querysidetextuuid uuid: %+v\n", uuid)
	// log.Printf("Querysidetextuuid targetUuid: %+v\n", targetUuid)

	var sideTextPlatform Sidetextplatform

	var sideTextUuid string
	var forwardUuid string
	var backwardUuid string
	var forwardPlatformUuid string
	var forwardPlatform string
	var backwardPlatformUuid string
	var backwardPlatform string

	if uuid < targetUuid {
		forwardUuid = uuid
		backwardUuid = targetUuid
	} else {
		forwardUuid = targetUuid
		backwardUuid = uuid
	}

	row := database.QueryRow("select sideTextUuid,forwardUuid,backwardUuid,forwardPlatformUuid,forwardPlatform,backwardPlatformUuid,backwardPlatform from sideText where forwardUuid = ? and backwardUuid = ?", forwardUuid, backwardUuid)

	err := row.Scan(&sideTextUuid, &forwardUuid, &backwardUuid, &forwardPlatformUuid, &forwardPlatform, &backwardPlatformUuid, &backwardPlatform)
	sideTextPlatform.Sidetextuuid = sideTextUuid
	if forwardUuid == targetUuid {
		sideTextPlatform.Userplatform.Useruuid = forwardUuid
		sideTextPlatform.Userplatform.Platformuuid = forwardPlatformUuid
		sideTextPlatform.Userplatform.Platform = forwardPlatform
	} else {
		sideTextPlatform.Userplatform.Useruuid = backwardUuid
		sideTextPlatform.Userplatform.Platformuuid = backwardPlatformUuid
		sideTextPlatform.Userplatform.Platform = backwardPlatform
	}

	return sideTextPlatform, err
}
