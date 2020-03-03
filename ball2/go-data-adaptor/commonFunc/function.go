package commonFunc

import (
	"net"
	"os"
	"strings"

	"../data"
	"../external/redis"

	"net/http"
	_ "net/http/pprof"
)

func InitFunc() {

	go redis.Subscribe("roomMessage", SubRedisRoomMessage)

}

func MyIp() string {

	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	// log.Printf("localAddr : %+v\n", localAddr)

	idx := strings.LastIndex(localAddr, ":")

	myIp := localAddr[0:idx]

	return myIp
}

func Exception(msg string, userUuid string, err error) data.Exception {

	if msg == "" {
		return data.Exception{}
	}

	var code = EsSysErrorLog(msg, userUuid, err)

	return data.Exception{Code: code, Message: msg}
}

func GetNamiApi(url string, data map[string]string) (resp *http.Response, err error) {

	req, err := http.NewRequest("GET", os.Getenv("namiApi")+url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	q.Add("user", "xcqq")
	q.Add("secret", "f706c0ca1548287fc471q")

	for key, val := range data {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return http.DefaultClient.Do(req)
}

func PostCdnUploadLink(fileLink string, fileName string, uploadPath string) (*http.Response, error) {

	apiUrl := os.Getenv("cdnApi") + "/uploadLink"

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("fileLink", fileLink)
	q.Add("fileName", fileName)
	q.Add("uploadPath", uploadPath)

	req.URL.RawQuery = q.Encode()

	return http.DefaultClient.Do(req)
}
