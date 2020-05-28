package common

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"
	"net/url"

	"server/database"
)

var apiKeyList = map[string]string{}
var mutexApiKeyList = new(sync.Mutex)

func ResponseWithJson(w http.ResponseWriter, code int, data map[string]interface{}) {
	jsonStr, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonStr)
}

func Queryapikey() {

	rows, err := database.Query("select apiKey, clientName from apiKeyList")
	if err != nil {
		// log.Printf("Queryapikey err : %+v\n", err)
		return
	}

	var list = map[string]string{}
	var apiKey string
	var clientName string
	for rows.Next() {
		rows.Scan(&apiKey, &clientName)
		// log.Printf("ip : %+v\n", ip)
		list[apiKey] = clientName
		// log.Printf("stamp : %+v\n", stamp)
		// log.Printf("BlockchatList : %+v\n", BlockchatList)
	}
	rows.Close()

	mutexApiKeyList.Lock()
	defer mutexApiKeyList.Unlock()
	apiKeyList = list
}

func Checkapikey(apiKey string, clientName string) bool {
	mutexApiKeyList.Lock()
	defer mutexApiKeyList.Unlock()
	value, ok := apiKeyList[apiKey]
	return (ok && value == clientName)
}

func Apicheck(w http.ResponseWriter, r *http.Request, path string) bool {

	if r.URL.Path != path {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return false
	}

	if r.Method != "POST" {
		http.Error(w, "not post.", http.StatusNotFound)
		return false
	}

	apiKeyData, apiKeyOk := r.Header["Api-Key"]
	clientNameData, clientNameOk := r.Header["Client-Name"]
	if !apiKeyOk {
		// log.Printf("no Api-Key")
		http.Error(w, "no Api-Key.", http.StatusNotFound)
		return false
	} else if !clientNameOk {
		// log.Printf("no Client-Name")
		http.Error(w, "no Client-Name.", http.StatusNotFound)
		return false
	} else {
		apiKey := apiKeyData[0]
		clientName := clientNameData[0]
		if !Checkapikey(apiKey, clientName) {
			// log.Printf("Api-Key err")
			http.Error(w, "Api-Key err", http.StatusNotFound)
			return false
		}
	}

	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm() err", http.StatusNotFound)
		return false
	}
	// log.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
	return true
}

func PostApiForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Api-Key", "JTEuXtr47S8hI8SmsnmtuVPl4roHNCX-OSWCKTu6KQLnOEwM2mwj_2vL6UD89cOUOW89VXyBdMtPxA9fpeiMRwGYILD7LO6mvNWn3hL8UfR38zY7QQjogdzNuz12g0mo")
	req.Header.Set("Client-Name", "MM-Server")
	return http.DefaultClient.Do(req)
}

func Checkplatformuser(platform string, platformUuid string, token string) (bool, error) {
	var err error

	redisToken, ok := Getredisusertoken(platform + "_" + platformUuid)
	if ok && redisToken == token {
		return true, err
	}

	urlLink := os.Getenv("userCheckUrl"+platform) + "/auth/user/tokenCheck"
	data := url.Values{"token": {token}}
	req, err := http.NewRequest("POST", urlLink, strings.NewReader(data.Encode()))
	if err != nil {
		Essyserrorlog("COMMON_CHECKPLATFORMUSER_NEW_REQUEST_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Essyserrorlog("COMMON_CHECKPLATFORMUSER_REQUEST_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Essyserrorlog("COMMON_CHECKPLATFORMUSER_BODY_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	var platformResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &platformResult)
	if err != nil {
		Essyserrorlog("COMMON_CHECKPLATFORMUSER_JSON_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	if platformResult.Code != 10000 {
		Essyserrorlog("COMMON_CHECKPLATFORMUSER_PLATFORM_UUID_ERROR", platform+"-"+platformUuid, nil)
		return false, nil
	}
	Setredisusertoken(platform+"_"+platformUuid, token)
	return true, err
}

func Getplatformuser(platform string, platformUuid string) (bool, error) {
	urlLink := os.Getenv("userCheckUrl"+platform) + "/users"
	data := url.Values{"uuid": {platformUuid}}
	// log.Printf("Getplatformuser urlLink : %+v\n", urlLink)
	// log.Printf("Getplatformuser data : %+v\n", data)
	req, err := http.NewRequest("POST", urlLink, strings.NewReader(data.Encode()))
	if err != nil {
		Essyserrorlog("COMMON_GETPLATFORMUSER_NEW_REQUEST_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Essyserrorlog("COMMON_GETPLATFORMUSER_REQUEST_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Essyserrorlog("COMMON_GETPLATFORMUSER_BODY_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	// log.Printf("Getplatformuser body : %+v\n", string(body))
	var platformResult struct {
		Code    int                      `json:"code"`
		Message string                   `json:"message"`
		Data    []map[string]interface{} `json:"data"`
	}
	err = json.Unmarshal(body, &platformResult)
	if err != nil {
		Essyserrorlog("COMMON_GETPLATFORMUSER_JSON_ERROR", platform+"-"+platformUuid, err)
		return false, err
	}

	// log.Printf("Getplatformuser platformResult : %+v\n", platformResult)
	if platformResult.Code != 10000 {
		Essyserrorlog("COMMON_GETPLATFORMUSER_PLATFORM_UUID_ERROR", platform+"-"+platformUuid, nil)
		return false, nil
	}

	return true, err
}
func Alivecheck() {

	defer log.Println("Alivecheck end")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// log.Println("t:", t)
			// log.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
			for ip, _ := range Roomspopulation {
				_, err := PostApiForm("http://"+ip+"/emit/getSnowFlake", url.Values{})
				if err != nil {
					Essyserrorlog("COMMON_ALIVECHECK_ERROR", "myIp : "+Myiplastdigit()+"  targetIp : "+ip, err)
					delete(Roomspopulation, ip)
					Essyslog("delete Roomspopulation ip : "+ip, Myiplastdigit(), "")
				}
			}
		case <-interrupt:
			// log.Println("interrupt")
			return
		}
	}

}
