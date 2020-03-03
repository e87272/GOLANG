package launchScreen

import (
	"sync"

	commonData ".."
)

var mutexLaunchScreen = new(sync.Mutex)

var launchScreen = []commonData.Announcement{}

func GetLaunchScreen() []commonData.Announcement {

	mutexLaunchScreen.Lock()
	defer mutexLaunchScreen.Unlock()

	return launchScreen
}

func SetLaunchScreen(launchScreenList []commonData.Announcement) {

	mutexLaunchScreen.Lock()
	defer mutexLaunchScreen.Unlock()

	launchScreen = launchScreenList
	return
}
