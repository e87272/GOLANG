

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
  var playerLogout = { "cmd" : "4", "idem" : Date.now().toString() ,"payload": { roomname : document.getElementById("roomname").value}}
  
  if (ws) {
    ws.send(JSON.stringify(playerLogout));
  }
}
  
function joinroom(){
	var playerEnterroom = { "cmd" : "10", "idem" : Date.now().toString() ,"payload": { roomname : document.getElementById("roomname").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(playerEnterroom));
	}
}

function exitroom(){
	var playerExitroom = { "cmd" : "8", "idem" : Date.now().toString() ,"payload": { roomname : document.getElementById("roomname").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(playerExitroom));
	}
}

function roomsearch(){
	var roomsearch = { "cmd" : "16", "idem" : Date.now().toString() }
	  
	if (ws) {
	  ws.send(JSON.stringify(roomsearch));
	}
}

function usersearch(){
	var roomsearch = { "cmd" : "6", "idem" : Date.now().toString() ,"payload": { roomname : document.getElementById("roomlist").value}}
	  
	if (ws) {
	  ws.send(JSON.stringify(roomsearch));
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
  var loginInfo = {"roomname" : document.getElementById("roomname").value,"nickname" : document.getElementById("nickname").value};
	var changeToken = { "cmd" : "2", "idem" : Date.now().toString(), "payload": loginInfo };
	ws.send(JSON.stringify(changeToken));

	let inputTxt = document.getElementById("input");
	let sendBtn = document.getElementById("sendBtn");

	sendBtn.disabled = false;
	sendBtn.onclick = function () {
		var idem = Date.now().toString();
		inputTxtQueue[idem] = inputTxt.value;
		inputTxt.value = "";
		var sendMessage = { "cmd" : "80", "idem" : idem, "payload": {"text":inputTxtQueue[idem],"style":"1","roominfo":{"roomname" : document.getElementById("roomlist").value,"nickname" : document.getElementById("nickname").value}} }
		ws.send(JSON.stringify(sendMessage));
	};
}

function receivePacketHandle(msg) {
	var packet = JSON.parse(msg)
	switch(packet.cmd){
		case "3":
			//addMessage("我加入聊天室");
　			document.getElementById('nickname').disabled=true;　// 變更欄位為禁用
　			document.getElementById('join').hidden=false;
　			document.getElementById('exit').hidden=false;
		break;
		case "5":
			addMessage("我段開連線");
      		ws.close();
	  　	document.getElementById('nickname').disabled=false;　// 變更欄位為禁用
　			document.getElementById('join').hidden=true;
　			document.getElementById('exit').hidden=true;
		break;
		case "7":
			addMessage("聊天室成員:" + packet["payload"]);
		break;
		case "17":
			addMessage("聊天室列表:" + packet["payload"]);
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
			addMessage(packet["payload"]["from"]["nickname"] + "加入" + packet["payload"]["roominfo"]["roomname"] +"聊天室");
			if(packet["payload"]["from"]["nickname"] == document.getElementById("nickname").value){
				var option = document.createElement("option");
				option.text = packet["payload"]["roominfo"]["roomname"];
				option.value = packet["payload"]["roominfo"]["roomname"];
				option.selected = true;
				option.id = "room-"+packet["payload"]["roominfo"]["roomname"];
				var select = document.getElementById("roomlist");
				select.appendChild(option);
			}
		break;
		case "15":
			addMessage(packet["payload"]["from"]["nickname"] + "離開" + packet["payload"]["roominfo"]["roomname"] +"聊天室");
			if(packet["payload"]["from"]["nickname"] == document.getElementById("nickname").value){
				var option = document.getElementById("room-"+packet["payload"]["roominfo"]["roomname"]);
				if (option) document.getElementById("roomlist").removeChild(option);
			}
		break;
		case "51":
			addMessage(packet["payload"]["from"]["nickname"] + " 在 " + packet["payload"]["roominfo"]["roomname"] +" 聊天室說" + ": " + packet["payload"]["text"]);
		break;
		default:

	}
}