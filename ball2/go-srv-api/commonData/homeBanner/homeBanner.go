package homeBanner

import (
	"sync"

	commonData ".."
)

var mutexHomeBanner = new(sync.Mutex)

var homeBanner = []commonData.Announcement{}

func GetHomeBanner() []commonData.Announcement {

	mutexHomeBanner.Lock()
	defer mutexHomeBanner.Unlock()

	return homeBanner
}

func SetHomeBanner(homeBannerList []commonData.Announcement) {

	mutexHomeBanner.Lock()
	defer mutexHomeBanner.Unlock()

	homeBanner = homeBannerList
	return
}
