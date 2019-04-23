package im
import( 
		"encoding/json"
		"strings"
		"github.com/kataras/iris/websocket"
)

type IMSocketPacket struct {
	Cmdid string
	Name string
	Msg string
}

func ReceiveImPacket(c websocket.Connection,msg string) {
	
	switch(strings.Split(msg, "_$_")[0]){
		case "101":
			joinPublicRoom(c,msg);
		break;
		case "102":
			speakPublicRoom(c,msg);
		break;
		case "201":
			joinPrivateRoom(c,msg);
		break;
		case "202":
			leavePrivateRoom(c,msg);
		break;
		case "203":
			speakPrivateRoom(c,msg);
		break;
	}
}

func joinPublicRoom(c websocket.Connection,msg string) {
	
	roomStateChange(c,msg,"加入公頻");
}

func joinPrivateRoom(c websocket.Connection,msg string) {

	roomStateChange(c,msg,"離開公頻");
	
	roomStateChange(c,msg,"加入聊天室");

}

func leavePrivateRoom(c websocket.Connection,msg string) {
	
	roomStateChange(c,msg,"離開聊天室");
	
	roomStateChange(c,msg,"加入公頻");
}

func speakPublicRoom(c websocket.Connection,msg string) {
	
	packet := IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:strings.Split(msg, "_$_")[2]};
	packettojson, _ := json.Marshal(packet)
	c.To("Public").Emit("imsocket", packettojson);
}

func speakPrivateRoom(c websocket.Connection,msg string) {
	
	packet := IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:strings.Split(msg, "_$_")[2]};
	packettojson, _ := json.Marshal(packet)
	c.To(strings.Split(msg, "_$_")[3]).Emit("imsocket", packettojson);
}


func roomStateChange(c websocket.Connection,msg string , state string) {
	var packet IMSocketPacket;
	var packettojson []byte;
	switch(state){
		case "加入公頻":

			c.Join("Public");
			packet = IMSocketPacket{Cmdid:"101",Name:strings.Split(msg, "_$_")[1],Msg:"加入公頻"};
			packettojson, _ = json.Marshal(packet)
			c.Emit("imsocket", packettojson);
			
			packet = IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:"加入公頻"};
			packettojson, _ = json.Marshal(packet)
			c.To("Public").Emit("imsocket", packettojson);

			break;
		case "離開公頻":

			packet = IMSocketPacket{Cmdid:"101",Name:strings.Split(msg, "_$_")[1],Msg:"離開公頻"};
			packettojson,_ = json.Marshal(packet)
			c.Emit("imsocket", packettojson);
		
			packet = IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:"離開公頻"};
			packettojson,_ = json.Marshal(packet)
			c.To("Public").Emit("imsocket", packettojson);
			
			c.Leave("Public");
			break;
		case "加入聊天室":
			
			c.Join(strings.Split(msg, "_$_")[2]);
			packet = IMSocketPacket{Cmdid:"201",Name:strings.Split(msg, "_$_")[1],Msg:"加入" + strings.Split(msg, "_$_")[2] +"聊天室"};
			packettojson,_ = json.Marshal(packet)
			c.Emit("imsocket", packettojson);
			
			packet = IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:"加入" + strings.Split(msg, "_$_")[2] +"聊天室"};
			packettojson, _ = json.Marshal(packet)
			c.To(strings.Split(msg, "_$_")[2]).Emit("imsocket", packettojson);
			
			break;
		case "離開聊天室":

			packet = IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:"離開聊天室"};
			packettojson,_ = json.Marshal(packet)
			c.Emit("imsocket", packettojson);
		
			packet = IMSocketPacket{Cmdid:"102",Name:strings.Split(msg, "_$_")[1],Msg:"離開" + strings.Split(msg, "_$_")[2] +"聊天室"};
			packettojson,_ = json.Marshal(packet)
			c.To(strings.Split(msg, "_$_")[2]).Emit("imsocket", packettojson);
			
			c.Leave(strings.Split(msg, "_$_")[2]);
			break;
	}
}

