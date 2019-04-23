package gm
import( "fmt"
		"github.com/kataras/iris/websocket"
		"math/rand"
		"time"
		"encoding/json"
		"strings"
		"strconv"
		"github.com/go-xorm/xorm"
		"github.com/go-xorm/core"
		_ "github.com/go-sql-driver/mysql"
)

func ReceiveCustomPacket(c websocket.Connection,msg string) {
	//fmt.Printf("%s : %s\n", "receiveCustomPacket", msg)
	
	switch(strings.Split(msg, "_$_")[0]){
		case "101":
			gameResult(c,msg);
		break;
	}
}

type GMSocketPacket struct {
	Cmdid string
	Msg string
}

type game_result struct {
	Gameindex   string `xorm:"not null"`
	Gameid      string `xorm:"not null"`
	Playerid    string `xorm:"not null"`
	Winlose 	int `xorm:"not null"`
	Gameresult  string `xorm:"not null"`
}

/**
產出結果
*/
func gameResult(c websocket.Connection,msg string){
	
	gameContent := strings.Split(msg, "_$_")

	rand.Seed(time.Now().UnixNano());
	//遊戲輪帶
	reelStrip := map[int]map[int]string{
		0:map[int]string{
			0:"3_",1:"5_",2:"1_",3:"2_",4:"3_",5:"4_",6:"4_",7:"4_",8:"5_",9:"6_",10:"3_",11:"5_",12:"7_",
		},
		1:map[int]string{
			0:"1_",1:"2_",2:"4_",3:"3_",4:"4_",5:"7_",6:"7_",7:"5_",8:"6_",9:"7_",
		},
		2:map[int]string{
			0:"1_",1:"6_",2:"2_",3:"7_",4:"6_",5:"7_",6:"6_",7:"6_",8:"6_",9:"3_",10:"5_",11:"4_",12:"5_",13:"6_",14:"7_",15:"7_",
		},
	}
	//賠率
	odds := map[string]int{
		"1_1_1_":100,"2_2_2_":40,"3_3_3_":20,"4_4_4_":10,"5_5_5_":10,"6_6_6_":5,"7_7_7_":4,"A_A_A_":5,"B_B_B_":3,"lose":0,
	}
	
	anyType := map[string]map[int]string{
		"A_":map[int]string{
			0:"1_",1:"2_",2:"3_",3:"4_",
		},
		"B_":map[int]string{
			0:"5_",1:"6_",2:"7_",
		},
	}
	winContent := [3]string{"","",""};
	winStr := "";
	packet := GMSocketPacket{Cmdid:"101",Msg:"",};

	//隨機轉輪
	i := 0;
	for i < len(reelStrip) {
		winContent[i] = reelStrip[i][rand.Intn(len(reelStrip[i]))];
		winStr = winStr + winContent[i];
		i = i + 1;
	}
	
	symbolWinStr := winStr;

	//中獎判斷
	if _, ok := odds[winStr]; ok{
		packet = GMSocketPacket{Cmdid:"101",Msg:winStr};
	}else{
		i := 0;
		for i < len(winContent) {
			//anySymbol判斷
			for anytypekey, anytypevalue := range anyType {
				j := 0;
				for j < len(anytypevalue) {
					if winContent[i] == anytypevalue[j]{
						winStr = strings.Replace(winStr, winContent[i], anytypekey, -1)
					}
					j = j +1;
				}
			}
			i = i + 1;
		}
		
		if _, ok := odds[winStr]; ok{ 
			packet = GMSocketPacket{Cmdid:"101",Msg:winStr};
		}else{
			packet = GMSocketPacket{Cmdid:"101",Msg:"lose"};
		}
	}
	//封包轉JSON
	packettojson, _ := json.Marshal(packet)
	//fmt.Printf("%s : %s\n", "packettojson", packettojson)
	
	//資料庫連線
	engine, _ := xorm.NewEngine("mysql", "gmlog:$test@/slot_777?charset=utf8")  // dbname是taoge
	//engine.ShowSQL(true)  // 显示SQL的执行, 便于调试分析
	
	engine.SetTableMapper(core.SnakeMapper{})
	//結構同步
	_ = engine.Sync2(new(game_result))
	
	//資料插入
	result := new(game_result)
	result.Gameindex = strconv.FormatInt(time.Now().Unix(),10)  + "_" + gameContent[1]
    result.Gameid = "1001"
    result.Playerid = "test"
	result.Winlose = odds[winStr]
	result.Gameresult = symbolWinStr;

	fmt.Printf("%+v\n", result)

	_ , _ = engine.Insert(result);
	//發送至客端
	c.Emit("gmsocket", string(packettojson));
}