package ipData

import "sync"

var whiteList = make(map[string]string)
var mutexWhiteList = new(sync.Mutex)

func CheckWhiteIp(ip string) bool {
	mutexWhiteList.Lock()
	defer mutexWhiteList.Unlock()
	_, ok := whiteList[ip]
	return ok
}

func SetWhiteList(ip string) {
	mutexWhiteList.Lock()
	defer mutexWhiteList.Unlock()
	whiteList[ip] = ip
}
