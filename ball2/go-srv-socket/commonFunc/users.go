package commonFunc

import (
	_ "net/http/pprof"

	"../commonData"
)

func UsersInfoInsert(userUuid string, userInfo commonData.UserInfo) {
	commonData.MutexUsersInfo.Lock()
	defer commonData.MutexUsersInfo.Unlock()

	commonData.UsersInfo[userUuid] = userInfo

	return
}

func UsersInfoRead(userUuid string) (commonData.UserInfo, bool) {
	commonData.MutexUsersInfo.Lock()
	defer commonData.MutexUsersInfo.Unlock()

	userInfo, ok := commonData.UsersInfo[userUuid]

	if !ok {
		return commonData.UserInfo{}, false
	}

	return userInfo, true
}

func UsersInfoDelete(userUuid string) {
	commonData.MutexUsersInfo.Lock()
	defer commonData.MutexUsersInfo.Unlock()

	delete(commonData.UsersInfo, userUuid)
}
