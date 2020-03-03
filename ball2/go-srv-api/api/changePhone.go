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

var changePhoneTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/changePhone`: `修改手機(舊手機收驗證簡訊)`,
	},
	Input:  jsObj{},
	Output: nil,
}

func changePhone(context *gin.Context) {

	exceptionPrefix := "API_CHANGEPHONE_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	phone := ""
	countryCode := ""
	row := database.QueryRow("SELECT `phone`, `country` FROM `account` WHERE `uuid` = ?",
		userUuid,
	)
	err := row.Scan(&phone, &countryCode)
	if err != nil {
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

	verifyType := "changePhone"
	verifyId := contentMap["sessionId"].(string)

	ginEngine.SetAuthSession(context, "verifyType", verifyType)
	ginEngine.SetAuthSession(context, "verifyId", verifyId)
	sendResultOk(context, nil)
	return
}
