package commonFunc

import (
	"../commonData"
	"../commonData/game"
	"../commonData/homeBanner"
	"../commonData/launchScreen"
	"../commonData/league"
	"../commonData/player"
	"../commonData/ranking"
	"../commonData/team"
	"../external/database"
)

func QueryLeague() {

	rows, _ := database.Query("SELECT `uuid`, `name`, `sequence` FROM `league`")

	for rows.Next() {
		var leagueUuid string
		var leagueName string
		var sequence int
		rows.Scan(&leagueUuid, &leagueName, &sequence)
		leagueInfo := commonData.LeagueInfo{}
		leagueInfo.LeagueCore.LeagueUuid = leagueUuid
		leagueInfo.LeagueCore.LeagueName = leagueName
		leagueInfo.Sequence = sequence
		leagueInfo.TeamList = make(map[string]commonData.TeamCore)
		league.SetLeagueInfo(leagueInfo)

	}
	rows.Close()

	rows, _ = database.Query("SELECT `league_uuid`, `team_uuid` FROM `league_team_list`")

	for rows.Next() {
		var leagueUuid string
		var teamUuid string

		rows.Scan(&leagueUuid, &teamUuid)

		// log.Printf("QueryLeague league_team_list: %+v\n", leagueUuid+"-"+teamUuid)

		teamInfo, ok := team.GetTeamInfo(teamUuid)

		if ok {
			league.SetLeagueInfoTeam(leagueUuid, teamInfo.TeamCore)
		}

	}
	rows.Close()

	return
}

func QueryTeam() {

	rows, _ := database.Query("SELECT `uuid`, `name`, `manager`, `venue`, `found` FROM `team`")

	for rows.Next() {

		teamInfo := commonData.TeamInfo{}

		rows.Scan(&teamInfo.TeamCore.TeamUuid, &teamInfo.TeamCore.Name, &teamInfo.Manager, &teamInfo.Venue, &teamInfo.Found)
		teamInfo.PlayerList = make(map[string]commonData.PlayerInfo)

		team.SetTeamInfo(teamInfo.TeamCore.TeamUuid, teamInfo)

	}
	rows.Close()

	rows, _ = database.Query("SELECT `team_uuid`, `player_uuid` FROM `team_member_list`")

	for rows.Next() {
		var teamUuid string
		var playerUuid string

		rows.Scan(&teamUuid, &playerUuid)

		playerInfo, ok := player.GetPlayerInfo(playerUuid)

		if ok {
			team.SetTeamInfoPlayer(teamUuid, playerUuid, playerInfo)
		}

	}
	rows.Close()

	return
}

func QueryPlayer() {

	rows, _ := database.Query("SELECT `uuid`, `name`, `country_uuid` FROM `player`")

	for rows.Next() {
		var playerUuid string
		var playerName string
		var countryUuid string

		rows.Scan(&playerUuid, &playerName, &countryUuid)
		playerInfo := commonData.PlayerInfo{}
		playerInfo.PlayerUuid = playerUuid
		playerInfo.PlayerName = playerName
		playerInfo.CountryUuid = countryUuid

		player.SetPlayerInfo(playerInfo)

	}
	rows.Close()

	return
}

func QueryAnnouncement() {

	rows, _ := database.Query("SELECT `uuid`, `type`, `sequence`, `content`, `url` FROM `announcement` WHERE `type` = ?", "LS")
	lanchScreen := []commonData.Announcement{}
	for rows.Next() {
		var announcement commonData.Announcement
		rows.Scan(&announcement.Uuid, &announcement.Type, &announcement.Sequence, &announcement.Content, &announcement.Url)

		lanchScreen = append(lanchScreen, announcement)

	}
	rows.Close()
	launchScreen.SetLaunchScreen(lanchScreen)

	rows, _ = database.Query("SELECT `uuid`, `type`, `sequence`, `content`, `url` FROM `announcement` WHERE `type` = ?", "HB")
	homeBannerInfo := []commonData.Announcement{}

	for rows.Next() {
		var announcement commonData.Announcement
		rows.Scan(&announcement.Uuid, &announcement.Type, &announcement.Sequence, &announcement.Content, &announcement.Url)

		homeBannerInfo = append(homeBannerInfo, announcement)

	}
	rows.Close()
	homeBanner.SetHomeBanner(homeBannerInfo)

	return
}

func QueryRank() {

	rows, _ := database.Query("SELECT `team_uuid`, `rank_type`, `score` FROM `rank` ORDER BY `rank`")
	rankMap := map[string][]commonData.RankInfo{}

	for rows.Next() {
		var teamUuid string
		var rankType string
		var score int
		rows.Scan(&teamUuid, &rankType, &score)

		teamInfo, ok := team.GetTeamInfo(teamUuid)
		if !ok {
			continue
		}

		list, ok := rankMap[rankType]
		if !ok {
			list = []commonData.RankInfo{}
		}

		rankInfo := commonData.RankInfo{
			TeamCore: teamInfo.TeamCore,
			Score:    score,
		}
		list = append(list, rankInfo)
		rankMap[rankType] = list
	}
	rows.Close()

	ranking.SetRankMap(rankMap)

	return
}

func QuerySchedule(queryStartTime int64, queryEndTime int64) {

	rows, err := database.Query(
		"SELECT `uuid`, `home_team`, `guest_team`, `start_time`, `league_uuid`, `home_score`, `guest_score`, `subtitle` FROM `schedule`"+
			" WHERE `start_time` >= ?"+
			" AND `start_time` < ?"+
			" ORDER BY `start_time`",
		queryStartTime, queryEndTime,
	)
	if err != nil {
		EsSysErrorLog("COMMONFUNC_QUERYSCHEDULE_SELECT_DB_ERROR", "", err)
		return
	}

	gameList := []commonData.GameInfo{}
	for rows.Next() {
		gameInfo := commonData.GameInfo{}
		rows.Scan(
			&gameInfo.Uuid,
			&gameInfo.HomeTeamUuid,
			&gameInfo.GuestTeamUuid,
			&gameInfo.StartTime,
			&gameInfo.LeagueUuid,
			&gameInfo.HomeScore,
			&gameInfo.GuestScore,
			&gameInfo.Subtitle,
		)
		gameList = append(gameList, gameInfo)
	}
	rows.Close()

	game.SetGameInfoBatch(gameList)
	return
}
