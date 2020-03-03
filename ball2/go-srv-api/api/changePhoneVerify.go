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

var changePhoneVerifyTestCase = testCase{
	Method: `POST`,
	Title: jsObj{
		`/changePhoneVerify`: `修改手機(舊手機輸入驗證碼)`,
	},
	Input: jsObj{
		`verifyCode`: `驗證碼`,
	},
	Output: `apiKey，在 setNewPhone 時要帶入`,
}

func changePhoneVerify(context *gin.Context) {

	exceptionPrefix := "API_CHANGEPHONEVERIFY_"

	userUuid, ok := ginEngine.GetAuthSession(context, "userUuid")
	if !ok {
		sendResultErr(context, commonFunc.Exception(exceptionPrefix+"NO_LOGIN", userUuid, nil))
		return
	}

	verifyType, ok := ginEngine.GetAuthSession(context, "verifyType")
	if !ok || verifyType != "changePhone" {
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

	ginEngine.DeleteAuthSession(context, "verifyType")
	ginEngine.DeleteAuthSession(context, "verifyId")

	apiKey := commonFunc.RandomString(16)

	ginEngine.SetAuthSession(context, "apiKey", apiKey)
	sendResultOk(context, apiKey)
	return
}
