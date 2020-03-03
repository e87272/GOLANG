package data

type NamiGameList struct {
	Teams map[string]struct {
		Id            int    //球队id
		Name_zh       string //中文名称
		Name_zht      string //粤语名称
		Name_en       string //英文名称
		Logo          string //logo,url前缀:http://cdn.sportnanoapi.com/football/team/
		Matchevent_id int    //赛事id
	}
	Events  map[string]NamiLeague
	Stages  map[string]NamiStage
	Matches [][]interface{}
}
type NamiLeague struct {
	Id             int    //id
	Name_zh        string //中文名称
	Short_name_zh  string //中文简称
	Name_zht       string //粤语名称
	Short_name_zht string //粤语简称
	Name_en        string //英文名称
	Short_name_en  string //英文简称
	Logo           string //logo,url前缀:http://cdn.sportnanoapi.com/football/competition/
}
type NamiStage struct {
	Id          int //id
	Mode        int //比赛方式, 1-积分 2-淘汰
	Group_count int //总分组数,0表示没有分组
	Round_count int //总轮数,0表示没有轮次
	Name_zh     string
	Name_zht    string
	Name_en     string
}

type NamiLeagueList struct {
	Areas       map[string]NamiAreaInfo
	Countrys    map[string]NamiCountryInfo
	Matchevents []NamiLeagueInfo
}

type NamiLeagueDetail struct {
	Areas       []NamiAreaInfo
	Countrys    []NamiCountryInfo
	Matchevents []NamiLeagueDetailInfo
}

type NamiAreaInfo struct {
	Id       int    //地区id
	Name_zh  string //中文
	Name_zht string //粤语
	Name_en  string //英文
}

type NamiCountryInfo struct {
	Id       int    //国家id
	Area_id  int    //地区id
	Name_zh  string //中文
	Name_zht string //粤语
	Name_en  string //英文
	Logo     string //logo,国旗url前缀：http://cdn.sportnanoapi.com/football/country/
}

type NamiLeagueInfo struct {
	Id             int    //赛事id
	Area_id        int    //地区id
	Country_id     int    //国家id
	Type           int    //赛事类型 0-未知 1-联赛 2-杯赛 3-友谊赛
	Level          int    //忽略，兼容使用
	Name_zh        string //中文
	Short_name_zh  string //中文缩写
	Name_zht       string //粤语
	Short_name_zht string //粤语缩写
	Name_en        string //英文
	Short_name_en  string //英文缩写
	Logo           string //logo,赛事url
}

type NamiLeagueDetailInfo struct {
	Id             string //赛事id
	Area_id        string //地区id
	Country_id     string //国家id
	Type           string //赛事类型 0-未知 1-联赛 2-杯赛 3-友谊赛
	Level          string //忽略，兼容使用
	Name_zh        string //中文
	Short_name_zh  string //中文缩写
	Name_zht       string //粤语
	Short_name_zht string //粤语缩写
	Name_en        string //英文
	Short_name_en  string //英文缩写
	Logo           string //logo,赛事url
}

type NamiTeamList []NamiTeamListInfo

type NamiTeamListInfo struct {
	Id             int    //球队id
	Matchevent_id  int    //赛事id
	Name_zh        string //中文
	Name_zht       string //粤语
	Name_en        string //英文
	Short_name_zh  string //中文简称
	Short_name_zht string //粤语简称
	Short_name_en  string //英文简称
	Logo           string //logo
	Found          string //球队成立日期
	Website        string //球队官网
	National       int    //是否国家队
	Country_logo   string //国家队logo，国家队才有
}

type NamiTeamInfo struct {
	Id                    int    //球队id
	Name_zh               string //中文
	Name_zht              string //粤语
	Name_en               string //英文
	Logo                  string //logo
	National              int    //是否国家队
	Market_value          int    //市值
	Market_value_currency string //市值单位
	Manager               struct {
		Id      int
		Name_zh string //中文
		Name_en string //英文
		Logo    string //	头像
	} //教練
	Venue struct {
		Id       int
		Name_zh  string //中文
		Capacity int    //容量
	} //场馆
	Players  []NamiTeamPlayerInfo   //球员列表
	Transfer []NamiTeamTransferInfo //球员转会
	Honor    []NamiTeamHonor        //榮譽
	Injury   []NamiTeamInjury       //伤停情况

}

type NamiTeamPlayerInfo struct {
	Shirt_number int
	Position     string //F-前锋 M-中场 D-后卫 G-守门员
	Player       struct {
		Id                    int
		Name_zh               string
		Name_en               string
		Logo                  string //球员logo,url前缀:http://cdn.sportnanoapi.com/football/player/
		Weight                int
		Height                int
		Birthday              int
		Country_id            int
		Nationality           string
		Preferred_foot        int    //0-未知 1-左脚 2-右脚 3-左右脚
		Contract_until        int    //俱乐部合同到期时间,可能无
		Market_value          int    //可能无
		Market_value_currency string //可能无
	}
}
type NamiTeamTransferInfo struct {
	Player struct {
		Id       int
		Name_zh  string
		Name_en  string
		Position string //F-前锋 M-中场 D-后卫 G-守门员
		Logo     string //球员logo,url前缀:http://cdn.sportnanoapi.com/football/player/
	}
	Transfer struct {
		Type          int //1-租借 2-租借结束 3-转会 4-退役 5-选秀 6-已解约 7-已签约 8-未知
		Transfer_fee  int //单位：欧元
		Transfer_time int
		From_id       int
		From_name     string
		To_id         int
		To_name       string
	}
}

type NamiTeamHonor struct {
	Id      int    //荣誉id
	Logo    string //荣誉logo,url前缀:http://cdn.sportnanoapi.com/football/honor/
	Name_zh string
	Seasons []string
	Detail  []struct {
		Competition_id int //赛事id,可能无,表示未关联
		Season_id      int //赛季id,可能无,表示未关联
		Season         string
	}
}

type NamiTeamInjury struct {
	Player struct {
		Id       int
		Name_zh  string
		Name_en  string
		Position string //F-前锋 M-中场 D-后卫 G-守门员
		Logo     string //球员logo,url前缀:http://cdn.sportnanoapi.com/football/player/
	}
	Reasons        string //伤停原因
	Missed_matches int    //影响场次
	Start_time     int    //开始时间
	End_time       int    //归队时间
	Type           int    //类型，1-受伤 2-停赛
}

type NamiRankListFifa struct {
	Items    []NamiRankInfoFifa
	Pub_time int
}
type NamiRankInfoFifa struct {
	Team struct {
		Id           int
		Name_zh      string //中文
		Name_zht     string //粤语
		Name_en      string //英文
		Logo         string //logo
		Country_logo string //国家队logo，国家队才有
	}
	Region_id        int //区域id  1-欧洲足联 2-南美洲足联 3-中北美洲及加勒比海足协 4-非洲足联 5-亚洲足联 6-大洋洲足联
	Ranking          int //当前排名
	Points           int //当前积分
	Previous_points  int //上次积分
	Position_changed int //排名变化

}

type NamiGameDetail struct {
	Info       NamiGameDetailInfo
	Matchevent struct {
		Season  string //赛季
		Id      int    //赛事id
		Name_zh string //赛事名称
	}
	Home_team NamiGameDetailTeamInfo
	Away_team NamiGameDetailTeamInfo
	Tlive     []struct {
		Main     int //是否重要事件
		Data     string
		Position int    //事件发生方,0-中立 1,主队 2,客队
		Type     int    //事件类型
		Time     string //事件时间
	}
	Stats []struct {
		Home int //主队值
		Away int //客队值
		Type int //技术统计类型，详见状态码->足球技术统计
	}
}

type NamiGameDetailInfo struct {
	Matchtime   int      //比赛时间
	Realtime    int      //开球时间，可能是上半场开球时间或者下半场开球时间
	Round       int      //比赛的轮次
	Statusid    int      //比赛状态，具体详见状态码
	Environment struct { //比赛环境,有数据才有此字段
		Pressure      string //气压
		Temperature   string //温度
		Wind          string //风速
		Humidity      string //湿度
		Weather_id    int    //天气id
		Weather       string //天气
		Weather_image string //天气logo
	}
}

type NamiGameDetailTeamInfo struct {
	Logo       string
	Half_score int
	Score      int
	Id         int
	Name_zh    string
}

type NamiGameLineup struct {
	Confirmed      int                    //是否确认是正式阵容，否则是预测阵容
	Home_formation string                 //主队阵型
	Away_formation string                 //客队阵型
	Home           []NamiGamePlayerLineup //主队阵容
	Away           []NamiGamePlayerLineup //客队阵容

}

type NamiGamePlayerLineup struct {
	Id           int        //球员id
	First        int        //是否首发，1-是 0-否
	Name         string     //球员名称
	Logo         string     //球员logo,地址前缀:http://cdn.sportnanoapi.com/football/player/
	Shirt_number int        //球衣号
	Position     string     //球员位置,F前锋 M中场 D后卫 G守门员,其他为未知
	X            int        //阵容x坐标 总共100
	Y            int        //阵容y坐标 总共100
	Rating       string     //评分
	Incidents    []struct { //球员事件
		Type       int         //事件类型
		Time       string      //事件发生时间（含加时时间）
		Minute     int         //事件发生时比赛分钟数
		Addtime    int         //加时时间（eg:中场时间伤停两分钟 time:'45+2' minute:45 addtime:2）
		Belong     int         //发生方 0-中立,1-主队,2-客队
		Text       string      //文本描述
		Home_score int         //主队比分
		Away_score int         //客队比分
		Player     interface{} //球员信息 player-相关球员 assist1-助攻球员1 assist2-助攻球员2 in_player 换上球员 out_player-换下球员
		In_player  interface{}
		Out_player interface{}
	}
}

type NamiPlayerInfo struct {
	Id                    int    //球员id
	Name_zh               string //中文
	Name_en               string //英文
	Team_id               int    //球队id，当1.球员退役 2.自由球员 3.球队未知 时team_id可能为0
	Logo                  string
	Birthday              int               //生日
	Weight                int               //体重
	Height                int               //身高
	Preferred_foot        int               //0-未知 1-左脚 2-右脚 3-左右脚
	Country_id            int               //国籍id
	Nationality           string            //国籍
	Position              string            //擅长位置
	Market_value          int               //可能无
	Market_value_currency string            //可能无
	Contract_until        int               //俱乐部合同到期时间,可能无
	Ability               [][]interface{}   //能力评分列表
	Characteristics       [][][]interface{} //技术特点列表
	Positions             []interface{}     //详细位置列表
	Honor                 []struct {        //荣誉列表
		Id      int
		Logo    string //logo,url前缀:http://cdn.sportnanoapi.com/football/honor/
		Name_zh string
		Seasons []struct {
			Competition_id int //赛事id,可能无,表示未关联
			Season_id      int //赛季id,可能无,表示未关联
			Season         string
			Team           struct {
				Id       int
				Name_zh  string
				Name_zht string
				Name_en  string
				Logo     string //球队logo，url前缀：http://cdn.sportnanoapi.com/football/team/
			}
		}
	}
}

type NamiSeasonList struct {
	Areas        []NamiAreaInfo
	Countries    []NamiCountryInfo
	Competitions []NamiSeasonListInfo
	Updated_at   int
}

type NamiSeasonListInfo struct {
	NamiLeague
	Seasons []struct {
		Id     int    //赛季id
		Season string //赛季年份
	}
}

type NamiSeasonInfo struct {
	Competition NamiSeasonLeague //赛事数据
	Stages      []NamiStage
	Teams       []struct {
		Id           int    //球队id
		Name_zh      string //中文名称
		Name_zht     string //粤语名称
		Name_en      string //英文名称
		Logo         string //logo,url前缀:http://cdn.sportnanoapi.com/football/team/
		Country_logo string //赛事id
	}
	Matches []struct {
		Id             int      //纳米比赛id
		Season_id      int      //赛季id
		Competition_id int      //赛事id
		Home_team_id   int      //主队id
		Away_team_id   int      //客队id
		Match_time     int      //比赛时间
		Status_id      int      //比赛状态
		Note           string   //比赛说明
		Home_scores    [6]int   //主队详细比分[比分(常规时间),半场比分,红牌,黄牌,角球，-1表示没有角球数据,加时比分(120分钟),加时赛才有,点球大战比分,点球大战才有]
		Away_scores    [6]int   //主队详细比分[比分(常规时间),半场比分,红牌,黄牌,角球，-1表示没有角球数据,加时比分(120分钟),加时赛才有,点球大战比分,点球大战才有]
		Round          struct { //阶段数据
			Stage_id  int //阶段id
			Group_num int //第几组,1-A 2-B以此类推
			Round_num int //第几轮
		}
		Position struct { //当前比赛的排名数据
			Home string //主队排名
			Away string //客队排名
		}
	}
	Table struct {
		Promotions []NamiPromotion
		Tables     []struct {
			Rows []struct {
				Team_id            int    //球队id
				Promotion_id       int    //升降级id
				Points             int    //积分
				Deduct_points      int    //扣除积分
				Note_zh            string //说明
				Position           int    //排名
				Total              int    //比赛场次
				Won                int    //胜的场次
				Draw               int    //平的场次
				Loss               int    //负的场次
				Goals              int    //进球
				Goals_against      int    //失球
				Goal_diff          int    //净胜球
				Home_points        int    //主场积分
				Home_position      int    //主场排名
				Home_total         int    //主场比赛场次
				Home_won           int    //主场胜的场次
				Home_draw          int    //主场平的场次
				Home_loss          int    //主场负的场次
				Home_goals         int    //主场进球
				Home_goals_against int    //主场失球
				Home_goal_diff     int    //主场净胜球
				Away_points        int    //客场积分
				Away_position      int    //客场排名
				Away_total         int    //客场比赛场次
				Away_won           int    //客场胜的场次
				Away_draw          int    //客场平的场次
				Away_loss          int    //客场负的场次
				Away_goals         int    //客场进球:失球
				Away_goals_against int    //客场失球
				Away_goal_diff     int    //客场净胜球
			}
			Id         int
			Conference string //分区信息，极少部分赛事才有，比如美职联
			Group      int    //不为0表示分组赛的第几组
			Stage_id   int    //所属阶段id
		}
	}
}
type NamiSeasonLeague struct {
	Id             int    //赛事id
	Type           int    //赛事类型 0-未知 1-联赛 2-杯赛 3-友谊赛
	Name_zh        string //中文
	Short_name_zh  string //中文缩写
	Name_zht       string //粤语
	Short_name_zht string //粤语缩写
	Name_en        string //英文
	Short_name_en  string //英文缩写
	Logo           string //logo,赛事url
	Cur_season_id  int    //当前赛季id
	Cur_stage_id   int    //当前阶段id
	Cur_round      int    //当前轮次
}

type NamiPromotion struct {
	Id       int
	Name_zh  string //中文名称
	Name_zht string //粤语名称
	Name_en  string //英文名称
	Color    string //颜色值
}

type NamiGameAnaysis struct {
	Teams             interface{} //球队列表
	Matchevents       interface{} //赛事列表
	Info              interface{} //当前比赛信息
	History           interface{} //历史交锋/近期战绩
	Goal_distribution interface{} //进球分布
	Injury            interface{} //伤停情况
	Table             interface{} //联赛积分
}
