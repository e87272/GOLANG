package syncApi

import (
	"../external/ginEngine"
)

func ApiRouter() {
	// 按字典序排列

	//同步公告
	ginEngine.GinEngine.POST("/sync/announcement", syncAnnouncement)
	appendTestCase(syncAnnouncementTestCase)

	//同步聯賽
	ginEngine.GinEngine.POST("/sync/league", syncLeague)
	appendTestCase(syncLeagueTestCase)

	//同步球員
	ginEngine.GinEngine.POST("/sync/player", syncPlayer)
	appendTestCase(syncPlayerTestCase)

	//同步排名
	ginEngine.GinEngine.POST("/sync/rank", syncRank)
	appendTestCase(syncRankTestCase)

	//同步賽程
	ginEngine.GinEngine.POST("/sync/schedule", syncSchedule)
	appendTestCase(syncScheduleTestCase)

	//同步隊伍
	ginEngine.GinEngine.POST("/sync/team", syncTeam)
	appendTestCase(syncTeamTestCase)

	createTestFile("syncTestSet.js")
}
