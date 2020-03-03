package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"../commonFunc"
	"../external/database"
	"../external/ginEngine"
	"../external/rongcloud"
	"github.com/gin-gonic/gin"
)

var loginVerifyTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/loginVerify`: `登入(輸入驗證碼)`,
	},
	Input: jsObj{
		`verifyCode`: `簡訊驗證碼`,
	},
	Output: nil,
}

func loginVerify(context *gin.Context) {

	exceptionPrefix := "API_LOGINVERIFY_"

	userUuid := ""
	verifyType, ok := ginEngine.GetAuthSession(context, "verifyType")
	if !ok || verifyType != "login" {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_VERIFY", userUuid, nil))
		return
	}

	verifyId, _ := ginEngine.GetAuthSession(context, "verifyId")
	verifyCode := context.PostForm("verifyCode")

	result, err := rongcloud.VerifyCode(verifyCode, verifyId)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+result, userUuid, err))
		return
	}

	contentMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(result), &contentMap)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"JSON_ERROR", userUuid, err))
		return
	}

	httpCode := contentMap["code"].(float64)
	if httpCode != 200 {
		err = errors.New(strconv.FormatFloat(httpCode, 'f', -1, 64))
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"FAIL", userUuid, err))
		return
	}

	ok = contentMap["success"].(bool)
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"FAIL", userUuid, nil))
		return
	}

	phone, _ := ginEngine.GetAuthSession(context, "phone")
	countryCode, _ := ginEngine.GetAuthSession(context, "countryCode")

	ginEngine.DeleteAuthSession(context, "phone")
	ginEngine.DeleteAuthSession(context, "countryCode")
	ginEngine.DeleteAuthSession(context, "verifyType")
	ginEngine.DeleteAuthSession(context, "verifyId")

	row := database.QueryRow("SELECT `uuid` FROM `account` WHERE `phone` = ? AND `country` = ?",
		phone, countryCode,
	)
	err = row.Scan(&userUuid)
	if err != nil {
		if err != database.ErrNoRows {
			sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
			return
		}
		userUuid = commonFunc.GetUuid()
		nickname := "user" + userUuid
		userIcon := "userIcon/default.png"
		_, err = database.Exec("INSERT INTO `account` (`uuid`, `phone`, `country`, `password`, `salt`, `nickname`, `icon`) VALUES (?, ?, ?, ?, ?, ?, ?)",
			userUuid, phone, countryCode, "", "", nickname, userIcon,
		)
		if err != nil {
			sendResultErr(context, commonFunc.Exception(exceptionPrefix+"INSERT_DB_ERROR", userUuid, err))
			return
		}
	}

	ginEngine.SetAuthSession(context, "userUuid", userUuid)
	sendResultOk(context, nil)
	return
}
