package api

import (
	"../external/ginEngine"
)

func ApiRouter() {
	// 按字典序排列

	// 公告
	ginEngine.GinEngine.GET("/announcement/:type", announcement)
	appendTestCase(announcementTestCase)

	// 取消訂閱
	ginEngine.GinEngine.POST("/cancelSubscribe/:type", cancelSubscribe)
	appendTestCase(cancelSubscribeTestCase)

	// cdn路徑
	ginEngine.GinEngine.POST("/cdnHost", cdnHost)
	appendTestCase(cdnHostTestCase)

	// 修改暱稱
	ginEngine.GinEngine.POST("/changeNickname", changeNickname)
	appendTestCase(changeNicknameTestCase)

	// 修改密碼
	ginEngine.GinEngine.POST("/changePassword", changePassword)
	appendTestCase(changePasswordTestCase)

	// 修改手機(舊手機收驗證簡訊)
	ginEngine.GinEngine.POST("/changePhone", changePhone)
	appendTestCase(changePhoneTestCase)

	// 修改手機(舊手機輸入驗證碼)
	ginEngine.GinEngine.POST("/changePhoneVerify", changePhoneVerify)
	appendTestCase(changePhoneVerifyTestCase)

	// 聯賽資訊
	ginEngine.GinEngine.POST("/leagueInfo", leagueInfo)
	appendTestCase(leagueInfoTestCase)

	// 聯賽列表
	ginEngine.GinEngine.POST("/leagueList", leagueList)
	appendTestCase(leagueListTestCase)

	// 登入(收驗證簡訊)
	ginEngine.GinEngine.POST("/login", login)
	appendTestCase(loginTestCase)

	// 登入(輸入驗證碼)
	ginEngine.GinEngine.POST("/loginVerify", loginVerify)
	appendTestCase(loginVerifyTestCase)

	// 登出
	ginEngine.GinEngine.POST("/logout", logout)
	appendTestCase(logoutTestCase)

	// 使用者資訊
	ginEngine.GinEngine.POST("/myInfo", myInfo)
	appendTestCase(myInfoTestCase)

	// 用密碼登入
	ginEngine.GinEngine.POST("/passwordLogin", passwordLogin)
	appendTestCase(passwordLoginTestCase)

	// 隨機字串(產生密鑰用)
	ginEngine.GinEngine.POST("/randomString", randomString)
	appendTestCase(randomStringTestCase)

	// 排名
	ginEngine.GinEngine.POST("/rank/:type", rank)
	appendTestCase(rankTestCase)

	// 賽程表
	ginEngine.GinEngine.POST("/schedule", schedule)
	appendTestCase(scheduleTestCase)

	// 首頁賽程表
	ginEngine.GinEngine.POST("/schedule/homePage", scheduleHomePage)
	appendTestCase(scheduleHomePageTestCase)

	// 熱門賽程表
	ginEngine.GinEngine.POST("/schedule/hot", scheduleHot)
	appendTestCase(scheduleHotTestCase)

	// 修改手機(新手機收驗證簡訊)
	ginEngine.GinEngine.POST("/setNewPhone", setNewPhone)
	appendTestCase(setNewPhoneTestCase)

	// 修改手機(新手機輸入驗證碼)
	ginEngine.GinEngine.POST("/setNewPhoneVerify", setNewPhoneVerify)
	appendTestCase(setNewPhoneVerifyTestCase)

	// 取雪花
	ginEngine.GinEngine.POST("/snowFlake", snowFlake)
	appendTestCase(snowFlakeTestCase)

	// 訂閱
	ginEngine.GinEngine.POST("/subscribe/:type", subscribe)
	appendTestCase(subscribeTestCase)

	// 批次訂閱
	ginEngine.GinEngine.POST("/subscribeBatch/:type", subscribeBatch)
	appendTestCase(subscribeBatchTestCase)

	// 訂閱列表 比賽
	ginEngine.GinEngine.POST("/subscribeList/game", subscribeListGame)
	appendTestCase(subscribeListGameTestCase)

	// 訂閱列表 聯賽
	ginEngine.GinEngine.POST("/subscribeList/league", subscribeListLeague)
	appendTestCase(subscribeListLeagueTestCase)

	// 訂閱列表 球員
	ginEngine.GinEngine.POST("/subscribeList/player", subscribeListPlayer)
	appendTestCase(subscribeListPlayerTestCase)

	// 訂閱列表 球隊
	ginEngine.GinEngine.POST("/subscribeList/team", subscribeListTeam)
	appendTestCase(subscribeListTeamTestCase)

	// 隊伍資訊
	ginEngine.GinEngine.POST("/teamInfo", teamInfo)
	appendTestCase(teamInfoTestCase)

	// 隊伍列表
	ginEngine.GinEngine.POST("/teamList", teamList)
	appendTestCase(teamListTestCase)

	// 取 AES Token
	ginEngine.GinEngine.POST("/token", token)
	appendTestCase(tokenTestCase)

	ginEngine.GinEngine.GET("/xcopy", xcopy)

	createTestFile("apiTestSet.js")
}
