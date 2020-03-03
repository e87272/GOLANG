package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"../commonFunc"
	"../external/ginEngine"
	"../external/rongcloud"
	"github.com/gin-gonic/gin"
)

var loginTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/login`: `登入(收驗證簡訊)`,
	},
	Input: jsObj{
		`phone`:       `手機號碼`,
		`countryCode`: `國碼(不含加號)`,
	},
	Output: nil,
}

func login(context *gin.Context) {

	exceptionPrefix := "API_LOGIN_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"ALREADY_LOGGED_IN", userUuid, nil))
		return
	}

	phone := context.PostForm("phone")
	countryCode := context.PostForm("countryCode")

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

	verifyType := "login"
	verifyId := contentMap["sessionId"].(string)

	ginEngine.SetAuthSession(context, "phone", phone)
	ginEngine.SetAuthSession(context, "countryCode", countryCode)
	ginEngine.SetAuthSession(context, "verifyType", verifyType)
	ginEngine.SetAuthSession(context, "verifyId", verifyId)
	sendResultOk(context, nil)
	return
}
