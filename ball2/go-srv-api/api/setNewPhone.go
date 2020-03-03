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

var setNewPhoneTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/setNewPhone`: `修改手機(新手機收驗證簡訊)`,
	},
	Input: jsObj{
		`apiKey`:      `changePhoneVerify 的回傳值`,
		`phone`:       `手機號碼`,
		`countryCode`: `國碼(不含加號)`,
	},
	Output: nil,
}

func setNewPhone(context *gin.Context) {

	exceptionPrefix := "API_SETNEWPHONE_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	apiKey, _ := ginEngine.GetAuthSession(context, "apiKey")
	if apiKey != context.PostForm("apiKey") {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"APIKEY_ERROR", userUuid, nil))
		return
	}

	phone := context.PostForm("phone")
	countryCode := context.PostForm("countryCode")

	row := database.QueryRow("SELECT `uuid` FROM `account` WHERE `phone` = ? AND `country` = ?",
		phone, countryCode,
	)
	err := row.Scan(&userUuid)
	if err == nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"ALREADY_REGISTERED", userUuid, nil))
		return
	}
	if err != database.ErrNoRows {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"SELECT_DB_ERROR", userUuid, err))
		return
	}

	result, err := rongcloud.SendCode(phone, countryCode)
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

	ginEngine.DeleteAuthSession(context, "apiKey")

	verifyType := "setNewPhone"
	verifyId := contentMap["sessionId"].(string)

	ginEngine.SetAuthSession(context, "phone", phone)
	ginEngine.SetAuthSession(context, "countryCode", countryCode)
	ginEngine.SetAuthSession(context, "verifyType", verifyType)
	ginEngine.SetAuthSession(context, "verifyId", verifyId)
	sendResultOk(context, nil)
	return
}
