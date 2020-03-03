var testSet = [
	{
		"Method": "GET",
		"Title": {
			"/announcement/HB": "首頁橫幅",
			"/announcement/LS": "起始頁"
		},
		"Input": {},
		"Output": [
			{
				"content": "公告內容",
				"sequence": "順序",
				"type": "\"HB\"或\"LS\"",
				"url": "連結網址",
				"uuid": "公告uuid"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/cancelSubscribe/game": "取消訂閱比賽",
			"/cancelSubscribe/league": "取消訂閱聯賽",
			"/cancelSubscribe/player": "取消訂閱球員",
			"/cancelSubscribe/team": "取消訂閱球隊"
		},
		"Input": {
			"targetUuid": "訂閱對象的uuid"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/cdnHost": "取得cdn路徑"
		},
		"Input": {},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/changeNickname": "更換暱稱"
		},
		"Input": {
			"nickname": "新的暱稱"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/changePassword": "更換密碼"
		},
		"Input": {
			"newPassword": "更換後的密碼",
			"oldPassword": "更換前的密碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/changePhone": "修改手機(舊手機收驗證簡訊)"
		},
		"Input": {},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/changePhoneVerify": "修改手機(舊手機輸入驗證碼)"
		},
		"Input": {
			"verifyCode": "驗證碼"
		},
		"Output": "apiKey，在 setNewPhone 時要帶入"
	},
	{
		"Method": "POST",
		"Title": {
			"/leagueInfo": "取得聯賽資訊"
		},
		"Input": {
			"leagueUuid": "聯賽uuid"
		},
		"Output": {
			"leagueCore": {
				"leagueName": "聯賽名稱",
				"leagueUuid": "聯賽uuid"
			},
			"teamList": "參賽球隊列表"
		}
	},
	{
		"Method": "POST",
		"Title": {
			"/leagueList": "取得聯賽列表"
		},
		"Input": {},
		"Output": [
			{
				"leagueCore": {
					"leagueName": "聯賽名稱",
					"leagueUuid": "聯賽uuid"
				},
				"sequence": "畫面排序",
				"teamCount": "參賽隊伍數量"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/login": "登入(收驗證簡訊)"
		},
		"Input": {
			"countryCode": "國碼(不含加號)",
			"phone": "手機號碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/loginVerify": "登入(輸入驗證碼)"
		},
		"Input": {
			"verifyCode": "簡訊驗證碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/logout": "登出"
		},
		"Input": {},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/myInfo": "自己的使用者資訊"
		},
		"Input": {
			"icon": "使用者頭像",
			"nickname": "使用者暱稱",
			"phone": "使用者手機號碼",
			"userUuid": "使用者uuid"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/passwordLogin": "透過密碼登入"
		},
		"Input": {
			"countryCode": "國碼(不含加號)",
			"password": "密碼",
			"phone": "手機號碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/randomString": "隨機字串(產生密鑰用)"
		},
		"Input": {},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/rank/fifa": "FIFA排名"
		},
		"Input": {},
		"Output": [
			{
				"score": "積分",
				"teamCore": {
					"name": "球隊名稱",
					"teamUuid": "球隊uuid"
				}
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/schedule": "賽程表"
		},
		"Input": {
			"date": "日期(毫秒時間戳)，若為空則查近期比賽",
			"leagueUuid": "聯賽uuid，若為空則查所有聯賽"
		},
		"Output": [
			{
				"guestScore": "客隊得分",
				"guestTeamUuid": "客隊uuid",
				"homeScore": "主隊得分",
				"homeTeamUuid": "主隊uuid",
				"leagueUuid": "所屬聯賽uuid",
				"startTime": "開始時間",
				"state": "比賽狀態",
				"subtitle": "子標題",
				"uuid": "比賽uuid"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/schedule/homePage": "首頁賽程表"
		},
		"Input": {},
		"Output": {
			"gameList": [
				{
					"guestScore": "客隊得分",
					"guestTeamUuid": "客隊uuid",
					"homeScore": "主隊得分",
					"homeTeamUuid": "主隊uuid",
					"leagueUuid": "所屬聯賽uuid",
					"startTime": "開始時間",
					"state": "比賽狀態",
					"subtitle": "子標題",
					"uuid": "比賽uuid"
				}
			],
			"gameTotal": "近期比賽總數"
		}
	},
	{
		"Method": "POST",
		"Title": {
			"/schedule/hot": "熱門賽程表"
		},
		"Input": {},
		"Output": [
			{
				"guestScore": "客隊得分",
				"guestTeamUuid": "客隊uuid",
				"homeScore": "主隊得分",
				"homeTeamUuid": "主隊uuid",
				"leagueUuid": "所屬聯賽uuid",
				"startTime": "開始時間",
				"state": "比賽狀態",
				"subtitle": "子標題",
				"uuid": "比賽uuid"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/setNewPhone": "修改手機(新手機收驗證簡訊)"
		},
		"Input": {
			"apiKey": "changePhoneVerify 的回傳值",
			"countryCode": "國碼(不含加號)",
			"phone": "手機號碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/setNewPhoneVerify": "修改手機(新手機輸入驗證碼)"
		},
		"Input": {
			"verifyCode": "簡訊驗證碼"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/snowFlake": "取得新的uuid"
		},
		"Input": {},
		"Output": "新的uuid"
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribe/game": "訂閱比賽",
			"/subscribe/league": "訂閱聯賽",
			"/subscribe/player": "訂閱球員",
			"/subscribe/team": "訂閱球隊"
		},
		"Input": {
			"targetUuid": "訂閱對象的uuid"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribeBatch/game": "批次訂閱比賽(重設比賽訂閱列表)",
			"/subscribeBatch/league": "批次訂閱聯賽(重設聯賽訂閱列表)",
			"/subscribeBatch/player": "批次訂閱球員(重設球員訂閱列表)",
			"/subscribeBatch/team": "批次訂閱球隊(重設球隊訂閱列表)"
		},
		"Input": {
			"targetUuid": "新的訂閱列表，逗號分隔"
		},
		"Output": null
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribeList/game": "取得比賽訂閱列表"
		},
		"Input": {},
		"Output": [
			{
				"guestScore": "客隊得分",
				"guestTeamUuid": "客隊uuid",
				"homeScore": "主隊得分",
				"homeTeamUuid": "主隊uuid",
				"leagueUuid": "所屬聯賽uuid",
				"startTime": "開始時間",
				"state": "比賽狀態",
				"subtitle": "子標題",
				"uuid": "比賽uuid"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribeList/league": "取得聯賽訂閱列表"
		},
		"Input": {},
		"Output": [
			{
				"leagueName": "聯賽名稱",
				"leagueUuid": "聯賽uuid",
				"teamList": "參賽球隊列表"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribeList/player": "取得球員訂閱列表"
		},
		"Input": {},
		"Output": [
			{
				"countryUuid": "所屬國家uuid",
				"playerName": "球員名稱",
				"playerUuid": "球員uuid"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/subscribeList/team": "取得球隊訂閱列表"
		},
		"Input": {},
		"Output": [
			{
				"found": "成立時間",
				"manager": "教練",
				"playerList": "球員列表",
				"teamCore": {
					"name": "球隊名稱",
					"teamUuid": "球隊uuid"
				},
				"venue": "主場"
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/teamInfo": "球隊資訊"
		},
		"Input": {
			"teamUuid": "球隊uuid"
		},
		"Output": {
			"found": "成立時間",
			"manager": "教練",
			"playerList": "球員列表",
			"teamCore": {
				"name": "球隊名稱",
				"teamUuid": "球隊uuid"
			},
			"venue": "主場"
		}
	},
	{
		"Method": "POST",
		"Title": {
			"/teamList": "球隊列表"
		},
		"Input": {},
		"Output": [
			{
				"playerCount": "球員數量",
				"teamCore": {
					"name": "球隊名稱",
					"teamUuid": "球隊uuid"
				}
			}
		]
	},
	{
		"Method": "POST",
		"Title": {
			"/token": "取得token(登入websocket用)"
		},
		"Input": {},
		"Output": "token(登入websocket用)"
	}
]