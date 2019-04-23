var nickname;
var roomName = "";
function join(){
	nickname = document.getElementById("nickname").value;
	var packet_str = "101_$_" +  nickname;
	socket.Emit("imsocket", packet_str);
	
}

function joinRoom(){
	roomName = document.getElementById("roomName").value;
	var packet_str = "201_$_" + nickname + "_$_" +  roomName;
	socket.Emit("imsocket", packet_str);
	im.innerHTML = "";
	document.getElementById("roomName").value = "";
}

function leaveRoom(){
	var packet_str = "202_$_" + nickname + "_$_" +  roomName;
	socket.Emit("imsocket", packet_str);
	im.innerHTML = "";
	roomName = "";
}

function send(){
	var msg = document.getElementById("msg").value;
	var packet_str
	if(roomName != ""){
		packet_str = "203_$_" + nickname + "_$_" +  msg + "_$_" + roomName;
	}else{
		packet_str = "102_$_" + nickname + "_$_" +  msg;
	}
	socket.Emit("imsocket", packet_str);
	document.getElementById("msg").value = "";
}
function receiveIMPacket(IMPacket) {
	var imMsg = document.getElementById("im");
	switch(IMPacket.Cmdid){
		case "101":
			document.getElementById("join").hidden = true;
			document.getElementById("joinroom").hidden = false;
			document.getElementById("sendmsg").hidden = false;
			document.getElementById("leaveRoom").hidden = true;
			imMsg.innerHTML = "";
		break;
		case "102":
			imMsg.innerHTML = IMPacket.Name + ":" + IMPacket.Msg + "\n" + imMsg.innerHTML;
		break;
		case "201":
			document.getElementById("joinroom").hidden = true;
			document.getElementById("leaveRoom").hidden = false;
			imMsg.innerHTML = "";
		break;
	}
}