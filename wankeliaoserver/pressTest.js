
var cmdList = {};
cmdList['3']  = onWebSocketTokenChange;
cmdList['11'] = onWebSocketAddRoom;
function onWebSocketOpen(ev) {
	changeToken(this);
};
function onWebSocketMessage(ev) {
	var packet = JSON.parse(ev.data);
	if (packet.cmd in cmdList && packet.result != 'err') {
		cmdList[packet.cmd](this,packet);
	};
};
function changeToken(ws) {
	var packet = {
		'cmd': '2',
		'idem': '',
		'payload': {
			'platform': 'MM',
			'platformUuid': '5d831e3e-5d2e-3bed-5c09-68b1-90e15ddf',
		},
	};
	ws.send(JSON.stringify(packet));
};
function addRoom(ws) {
	var packet = {
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomType': 'liveGroup',
			'roomUuid': '0005786c733a5000',
			'roomName': '',
			'adminSet': '',
		},
	};
	ws.send(JSON.stringify(packet));
};
var chatCount = 0;
function sendChatMessage(ws,roomInfo) {
	var packet = {
		'cmd': '80',
		'idem': '',
		'payload': {
			'roomInfo': roomInfo,
			'message': Date.now().toString()+'-'+chatCount,
			'style': '1',
		},
	};
	ws.send(JSON.stringify(packet));
	++chatCount;
};
function onWebSocketTokenChange(ws,packet) {
	addRoom(ws)
};
function onWebSocketAddRoom(ws,packet) {
	setInterval(sendChatMessage,1000,ws,packet.payload);
};
function createWebSocket() {
	var ws = new WebSocket('ws://' + document.location.host + '/echo');
	ws.onopen    = onWebSocketOpen;
	ws.onmessage = onWebSocketMessage;
};
!function() {
	for (var i=0; i<10; ++i) {
		setTimeout(createWebSocket,i*100);
	};
}();
DanMu.style.display = 'none';