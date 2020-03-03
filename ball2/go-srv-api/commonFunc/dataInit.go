package commonFunc

import "../external/stamp"

func DataInit() {

	// 撈取球員資訊
	QueryPlayer()

	// 撈取隊伍資訊
	QueryTeam()

	// 撈取分類資訊
	QueryLeague()

	// 撈取起始頁資訊
	QueryAnnouncement()

	// 撈取排名
	QueryRank()

	// 撈取賽程
	today := stamp.Today()
	queryStartTime := today - stamp.Day
	queryEndTime := today + 7*stamp.Day

	QuerySchedule(queryStartTime, queryEndTime)

	return
}
