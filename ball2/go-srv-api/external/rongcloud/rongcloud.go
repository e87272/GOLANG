package rongcloud

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../stamp"
)

const generalVerifyTemplateId = "8lDDP4tMQAF9MeIrHZZxj2"

func api(name string, body map[string]string) (string, error) {

	apiUrl := "http://api.sms.ronghub.com/" + name + ".json"
	bodyStr := ""
	if len(body) > 0 {
		for k, v := range body {
			bodyStr += "&" + k + "=" + v
		}
		bodyStr = bodyStr[1:]
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(bodyStr))
	if err != nil {
		return "NEW_REQUEST_ERROR", err
	}

	appKey := os.Getenv("appKey")
	appSecret := os.Getenv("appSecret")
	nonce := strconv.FormatInt(rand.Int63(), 16)
	timestamp := strconv.FormatInt(stamp.Now(), 10)
	signatureByte := sha1.Sum([]byte(appSecret + nonce + timestamp))
	signature := hex.EncodeToString(signatureByte[:])
	contentType := "application/x-www-form-urlencoded"

	req.Header.Set("App-Key", appKey)
	req.Header.Set("Nonce", nonce)
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Signature", signature)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "DO_REQUEST_ERROR", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "READ_CONTENT_ERROR", err
	}

	return string(result), nil
}

func SendCode(phone string, countryCode string) (string, error) {
	result, err := api("sendCode", map[string]string{
		"mobile":     phone,
		"region":     countryCode,
		"templateId": generalVerifyTemplateId,
	})
	return result, err
}

func VerifyCode(verifyCode string, verifyId string) (string, error) {
	result, err := api("verifyCode", map[string]string{
		"code":      verifyCode,
		"sessionId": verifyId,
	})
	return result, err
}
