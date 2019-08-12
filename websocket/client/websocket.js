

var scheme = document.location.protocol == "https:" ? "wss" : "ws";
var port = document.location.port ? ":" + document.location.port : "";
var wsURL = scheme + "://" + document.location.hostname + port + "/echo";
var ws;

var inputTxtQueue = [];

function connect(){
    if (ws) {
    	return false;
    }
    ws = new WebSocket(wsURL);
    ws.onopen = function(evt) {
      handleNamespaceConnectedConn();
    }
    ws.onclose = function(evt) {
      console.log("CLOSE");
      ws = null;
    }
    ws.onmessage = function(evt) {
      receivePacketHandle(evt.data);
    }
    ws.onerror = function(evt) {
	   console.log("ERROR: " + evt.data);
    }
  return true;   
}

function disconnect(){
  var playerLogout = { "cmd" : "4", "idem" : Date.now().toString() ,"payload": { roomName : document.getElementById("roomName").value}}
  
  if (ws) {
    ws.send(JSON.stringify(playerLogout));
  }
}
  
function joinroom(){
	var playerEnterroom = { "cmd" : "10", "idem" : Date.now().toString() ,"payload": { roomName : document.getElementById("roomName").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(playerEnterroom));
	}
}

function exitroom(){
	var playerExitroom = { "cmd" : "8", "idem" : Date.now().toString() ,"payload": { roomName : document.getElementById("roomName").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(playerExitroom));
	}
}

function roomSearch(){
	var roomSearch = { "cmd" : "16", "idem" : Date.now().toString() }
	  
	if (ws) {
	  ws.send(JSON.stringify(roomSearch));
	}
}

function userSearch(){
	var userSearch = { "cmd" : "6", "idem" : Date.now().toString() ,"payload": { roomName : document.getElementById("roomInput").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(userSearch));
	}
}

function addMessage(msg) {
	document.getElementById("output").innerHTML += msg + "\n";
}

function handleError(reason) {
	console.log(reason);
	window.alert("error: see the dev console");
}

function handleNamespaceConnectedConn() {
  var loginInfo = {"roomName" : document.getElementById("roomName").value,"nickName" : document.getElementById("nickName").value};
	var changeToken = { "cmd" : "2", "idem" : Date.now().toString(), "payload": loginInfo };
	ws.send(JSON.stringify(changeToken));

	let inputTxt = document.getElementById("input");
	let sendBtn = document.getElementById("sendBtn");

	sendBtn.disabled = false;
	sendBtn.onclick = function () {
		var idem = Date.now().toString();
		inputTxtQueue[idem] = inputTxt.value;
		inputTxt.value = "";
		var sendMessage = { "cmd" : "80", "idem" : idem, "payload": {"text":inputTxtQueue[idem],"style":"1","roomInfo":{"roomName" : document.getElementById("roomList").value,"nickName" : document.getElementById("nickName").value}} }
		ws.send(JSON.stringify(sendMessage));
	};
}

function receivePacketHandle(msg) {
	var packet = JSON.parse(msg)
	switch(packet.cmd){
		case "3":
			//addMessage("我加入聊天室");
　			document.getElementById('nickName').disabled=true;　// 變更欄位為禁用
　			document.getElementById('join').hidden=false;
　			document.getElementById('exit').hidden=false;
		break;
		case "5":
			addMessage("我段開連線");
      		ws.close();
	  　	document.getElementById('nickName').disabled=false;　// 變更欄位為禁用
　			document.getElementById('join').hidden=true;
　			document.getElementById('exit').hidden=true;
			document.getElementById("roomList").options.length = 0;
		break;
		case "7":
			var nameList = "";
			for (i= 0; i< packet["payload"].length; i++){
				nameList = nameList + packet["payload"][i].nickName + ","
			}
			addMessage("聊天室成員:" + nameList.substr(0 , nameList.length - 1));
		break;
		case "17":
			var roomList = "";
			for (i= 0; i< packet["payload"].length; i++){
				roomList = roomList + packet["payload"][i].roomName + ","
			}
			addMessage("聊天室列表:" + roomList.substr(0 , roomList.length - 1));
		break;
		case "9":
			console.log("CMD_R_PLAYER_EXIT_ROOM");
		break;
		case "11":
			console.log("CMD_R_PLAYER_ENTER_ROOM");
		break;
		case "81":
			if(packet["result"] == "ok"){
        		//addMessage("Me: " + inputTxtQueue[packet["idem"]]);
        		inputTxtQueue[packet["idem"]] = "";
			}else{
				console.log(packet);
			}
		break;
		case "13":
			addMessage(packet["payload"]["from"]["nickName"] + "加入" + packet["payload"]["roomInfo"]["roomName"] +"聊天室");
			if(packet["payload"]["from"]["nickName"] == document.getElementById("nickName").value){
				var option = document.createElement("option");
				option.text = packet["payload"]["roomInfo"]["roomName"];
				option.value = packet["payload"]["roomInfo"]["roomName"];
				option.selected = true;
				option.id = "room-"+packet["payload"]["roomInfo"]["roomName"];
				var select = document.getElementById("roomList");
				select.appendChild(option);
			}
		break;
		case "15":
			addMessage(packet["payload"]["from"]["nickName"] + "離開" + packet["payload"]["roomInfo"]["roomName"] +"聊天室");
			if(packet["payload"]["from"]["nickName"] == document.getElementById("nickName").value){
				var option = document.getElementById("room-"+packet["payload"]["roomInfo"]["roomName"]);
				if (option) document.getElementById("roomList").removeChild(option);
			}
		break;
		case "51":
			addMessage(packet["payload"]["from"]["nickName"] + " 在 " + packet["payload"]["roomInfo"]["roomName"] +" 聊天室說" + ": " + packet["payload"]["text"]);
		break;
		default:

	}
}