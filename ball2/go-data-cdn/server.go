package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"math/rand"
	"net/http"
	"os"
	"time"

	"./commonFunc"
	"./data"
	"./external/elasticSearch"
)

func main() {

	rand.Seed(time.Now().Unix())

	// log.Printf("NumCPU : %v\n", runtime.NumCPU())

	// runtime.GOMAXPROCS(runtime.NumCPU())

	// log.Printf(common.Getid().Hexstring())

	loadEnv()

	elasticSearch.EsInit()

	commonFunc.InitFunc()

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/uploadLink", uploadLink)

	http.HandleFunc("/uploadBase64", uploadBase64)

	http.ListenAndServe(os.Getenv("goCdnPort"), nil)

}

func fileUnitTestHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("r.URL.Path :", r.URL.Path)
	http.ServeFile(w, r, ".."+r.URL.Path)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./client/healthCheck.html")
}

func uploadLink(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadPath := r.FormValue("uploadPath")
	fileName := r.FormValue("fileName")
	fileLink := r.FormValue("fileLink")

	// Get the data
	resp, err := http.Get(fileLink)
	if err != nil {
		exception := commonFunc.Exception("MAIN_UPLOADLINK_LINK_ERROR", "", err)
		sendResultErr(w, exception)

		return
	}
	defer resp.Body.Close()

	if _, err := os.Stat(os.Getenv("cdnLocalPath") + "/" + uploadPath); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("cdnLocalPath")+"/"+uploadPath, os.ModePerm)
	}

	// 创建一个文件用于保存
	out, err := os.Create(os.Getenv("cdnLocalPath") + "/" + uploadPath + "/" + fileName)
	if err != nil {
		exception := commonFunc.Exception("MAIN_UPLOADLINK_CREATE_FILE_ERROR", "", err)
		sendResultErr(w, exception)

		return
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		exception := commonFunc.Exception("MAIN_UPLOADLINK_COPY_FILE_ERROR", "", err)
		sendResultErr(w, exception)

		return
	}

	sendResultOk(w, os.Getenv("cdnHost")+"/"+uploadPath+"/"+fileName)
}

func uploadBase64(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadPath := r.FormValue("uploadPath")
	fileName := r.FormValue("fileName")
	file := r.FormValue("file")

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(file)

	if err != nil {

		exception := commonFunc.Exception("MAIN_UPLOADBASE64_ERROR", "", err)
		sendResultErr(w, exception)

		return
	}

	if _, err := os.Stat(os.Getenv("cdnLocalPath") + "/" + uploadPath); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("cdnLocalPath")+"/"+uploadPath, os.ModePerm)
	}

	cdnFileLink := uploadPath + "/" + fileName
	ioutil.WriteFile(os.Getenv("cdnLocalPath")+"/"+cdnFileLink, decodeBytes, os.ModePerm)

	sendResultOk(w, os.Getenv("cdnHost")+"/"+cdnFileLink)

	return
}

type apiResult struct {
	Result  string         `json:"result"`
	Error   data.Exception `json:"error"`
	Payload interface{}    `json:"payload"`
}

func sendResultOk(w http.ResponseWriter, payload interface{}) {
	apiResult := apiResult{
		Result:  "ok",
		Payload: payload,
	}

	ResponseWithJson(w, http.StatusOK, apiResult)

}

func sendResultErr(w http.ResponseWriter, exception data.Exception) {
	apiResult := apiResult{
		Result: "err",
		Error:  exception,
	}

	ResponseWithJson(w, http.StatusOK, apiResult)

}

func ResponseWithJson(w http.ResponseWriter, code int, data interface{}) {
	jsonStr, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonStr)
}

func loadEnv() {

	config, err := ioutil.ReadFile("config/host.json")
	if err != nil {
		log.Fatal("找不到host.json")
	}

	configHost := make(map[string]string)
	// log.Printf("configHost : %v\n", configHost)

	err = json.Unmarshal(config, &configHost)
	if err != nil {
		// log.Printf("configHost err: %v\n", err)
		return
	}

	// log.Printf("configHost : %v\n", configHost)
	for k, v := range configHost {
		// log.Printf("%s : %s\n", k, v)
		_ = os.Setenv(k, v)
	}
}
