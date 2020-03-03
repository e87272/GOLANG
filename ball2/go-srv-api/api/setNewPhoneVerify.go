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

var setNewPhoneVerifyTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/setNewPhoneVerify`: `修改手機(新手機輸入驗證碼)`,
	},
	Input: jsObj{
		`verifyCode`: `簡訊驗證碼`,
	},
	Output: nil,
}

func setNewPhoneVerify(context *gin.Context) {

	exceptionPrefix := "API_SETNEWPHONEVERIFY_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	verifyType, ok := ginEngine.GetAuthSession(context, "verifyType")
	if !ok || verifyType != "setNewPhone" {
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

	_, err = database.Exec("UPDATE `account` SET `phone` = ?, `country` = ? WHERE `uuid` = ?",
		phone, countryCode, userUuid,
	)
	if err != nil {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"UPDATE_DB_ERROR", userUuid, err))
		return
	}

	sendResultOk(context, nil)
	return
}
