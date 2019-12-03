
/*
	'0014eb7efac03000': {
		'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
		'platform': 'MM',
	},
	'0014eb6e4dc03000': {
		'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
		'platform': 'MM',
	},
	'00157ec2da803000': {
		'platformUuid': '5db266d0-8251-2b17-f28e-c280-d81800cd',
		'platform': 'MM',
	},
	'': {
		'platformUuid': '',
		'platform': '',
	},
*/

var testSet = Object.create(null);

testSet['COMMAND_BLOCKROOMCHAT_NOT_ADMIN'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '24',
		'idem': '',
		'payload': {
			'userUuid': '0014eb6e4dc03000',
			'roomUuid': '000572bf8d4a5001',
			'blockTime': '1',
		},
	},
];
testSet['COMMAND_BLOCKROOMCHAT_ROOMUUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '24',
		'idem': '',
		'payload': {
			'userUuid': '0014eb6e4dc03000',
			'roomUuid': '?',
			'blockTime': '1',
		},
	},
];
testSet['COMMAND_BLOCKROOMCHAT_TARGET_IS_ADMIN'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '24',
		'idem': '',
		'payload': {
			'userUuid': '0014eb7efac03000',
			'roomUuid': '000572bf8d4a5001',
			'blockTime': '1',
		},
	},
];
testSet['COMMAND_BLOCKROOMCHAT_TIME_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '24',
		'idem': '',
		'payload': {
			'userUuid': '0014eb6e4dc03000',
			'roomUuid': '000572bf8d4a5001',
			'blockTime': '?',
		},
	},
];
testSet['COMMAND_BLOCKROOMCHAT_USER_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '24',
		'idem': '',
		'payload': {
			'userUuid': '?',
			'roomUuid': '000572bf8d4a5001',
			'blockTime': '1',
		},
	},
];
// testSet['COMMAND_CREATEPRIVATEROOM_NOT_ADMIN'] = [];
// testSet['COMMAND_CREATEPRIVATEROOM_ROOM_ICON_ERROR'] = [];
// testSet['COMMAND_DISMISSROOM_NOT_ADMIN'] = [];
// testSet['COMMAND_DISMISSROOM_ROOM_UUID_ERROR'] = [];
testSet['COMMAND_GETCHATHISTORY_ROOM_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '18',
		'idem': '',
		'payload': {
			'roomCore': {
				'roomUuid': '000572bf8d4a5001',
				'roomType': 'liveGroup',
			},
			'historyUuid': '',
		},
	},
];
testSet['COMMAND_GETMEMBERLIST_ROOM_TYPE_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'liveGroup',
		},
	},
	{
		'cmd': '6',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'vipGroup',
		},
	},
];
testSet['COMMAND_GETMEMBERLIST_ROOM_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '6',
		'idem': '',
		'payload': {
			'roomUuid': '0000000000000000',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_GETSIDETEXTHISTORY_CHATTARGET_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '34',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb6e4dc03000',
			'historyUuid': '',
		},
	},
];
// testSet['COMMAND_KICKROOMUSER_NOT_ADMIN'] = [];
// testSet['COMMAND_KICKROOMUSER_ROOM_TYPE_ERROR'] = [];
// testSet['COMMAND_KICKROOMUSER_ROOM_UUID_ERROR'] = [];
// testSet['COMMAND_MESSAGESEEN_GUEST'] = [];
// testSet['COMMAND_MESSAGESEEN_TARGET_ROOM_TYPE_ERROR'] = [];
// testSet['COMMAND_MESSAGESEEN_TARGET_ROOM_UUID_ERROR'] = [];
// testSet['COMMAND_MESSAGESEEN_TARGET_SIDE_TEXT_UUID_ERROR'] = [];
testSet['COMMAND_PLAYERENTERROOM_IN_ROOM'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'liveGroup',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_PLAYERENTERROOM_ROOM_READ_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '0000000000000000',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_PLAYERENTERROOM_ROOM_TYPE_NOT_WORD'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '0000000000000000',
			'roomType': '?',
		},
	},
];
testSet['COMMAND_PLAYERENTERROOM_ROOM_UUID_NULL'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_PLAYEREXITROOM_ROOM_TYPE_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '8',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': '',
		},
	},
];
testSet['COMMAND_PLAYEREXITROOM_ROOM_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '8',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_PLAYEREXITROOM_ROOM_UUID_NULL'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '8',
		'idem': '',
		'payload': {
			'roomUuid': '',
			'roomType': 'liveGroup',
		},
	},
];
testSet['COMMAND_PLAYERSENDMSG_CHAT_BLOCK'] = [
	{
		"cmd": "2",
		"idem": "",
		"payload": {
			"platformUuid": "5db266d0-8251-2b17-f28e-c280-d81800cd",
			"platform": "MM",
		},
	},
	{
		"cmd": "80",
		"idem": "",
		"payload": {
			"chatTarget": "000572bf8d4a5001",
			"message": "?",
			"style": "1",
		},
	},
];
testSet['COMMAND_PLAYERSENDMSG_GUEST'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '',
			'platform': '',
		},
	},
	{
		'cmd': '80',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '?',
			'style': '1',
		},
	},
];
testSet['COMMAND_PLAYERSENDMSG_MSG_TOO_LONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '80',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '........10........20........30........40........50........60........70........80........90.......100.......110.......120.......130.......140.......150.......160.......170.......180.......190.......200~~~',
			'style': '1',
		},
	},
];
testSet['COMMAND_PLAYERSENDMSG_NOT_IN_ROOM'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '80',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '?',
			'style': '1',
		},
	},
];
testSet['COMMAND_PLAYERSENDMSG_SPEAK_CD'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '10',
		'idem': '',
		'payload': {
			'roomUuid': '000572bf8d4a5001',
			'roomType': 'liveGroup',
		},
	},
	{
		'cmd': '80',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '?',
			'style': '1',
		},
	},
	{
		'cmd': '80',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '?',
			'style': '1',
		},
	},
];
testSet['COMMAND_SIDETEXTDELETE_GUEST'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '',
			'platform': '',
		},
	},
	{
		'cmd': '32',
		'idem': '',
		'payload': '0014eb6e4dc03000',
	},
];
testSet['COMMAND_SIDETEXTDELETE_NOT_SIDETEXT'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '32',
		'idem': '',
		'payload': '0000000000000000',
	},
];
testSet['COMMAND_SIDETEXTLIST_GUEST'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '',
			'platform': '',
		},
	},
	{
		'cmd': '30',
		'idem': '',
	},
];
testSet['COMMAND_SIDETEXTSEND_GUEST'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '',
			'platform': '',
		},
	},
	{
		'cmd': '82',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb6e4dc03000',
			'message': '?',
			'style': '1',
		},
	},
];
testSet['COMMAND_SIDETEXTSEND_MSG_TOO_LONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '82',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb6e4dc03000',
			'message': '........10........20........30........40........50........60........70........80........90.......100.......110.......120.......130.......140.......150.......160.......170.......180.......190.......200~~~',
			'style': '1',
		},
	},
];
testSet['COMMAND_SIDETEXTSEND_SIDE_TEXT_YOURSELF'] =[
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '82',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb7efac03000',
			'message': '?',
			'style': '1',
		},
	},
];
testSet['COMMAND_SIDETEXTSEND_SPEAK_CD'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '82',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb6e4dc03000',
			'message': '?',
			'style': '1',
		},
	},
	{
		'cmd': '82',
		'idem': '',
		'payload': {
			'chatTarget': '0014eb6e4dc03000',
			'message': '?',
			'style': '1',
		},
	},
];
// testSet['COMMAND_TARGETADDROOMBATCH_NOT_ADMIN'] = [];
// testSet['COMMAND_TARGETADDROOMBATCH_TARGET_UUID_ERROR'] = [];
testSet['COMMAND_TOKENCHANGE_JSON_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': 0,
	},
];
testSet['SHELL_BLOCKUSER_NOT_ADMIN'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BU 0014eb6e4dc03000 1',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_PARAMETER_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BU',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_ROOM_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '?',
			'message': '/su BU 0014eb6e4dc03000 1',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_TARGET_IS_ADMIN'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BU 0014eb7efac03000 1',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_TIME_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BU 0014eb6e4dc03000 ?',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_USER_UUID_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BU ? 1',
			'style': '1',
		},
	},
];
testSet['SHELL_BLOCKUSER_UUID_NULL'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '',
			'message': '/su BU 0014eb6e4dc03000 1',
			'style': '1',
		},
	},
];
testSet['SHELL_LINKPROCLAMATION_CONTENT_TOOLONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su LP 1 -c ........10........20........30........40........50~~~',
			'style': '1',
		},
	},
];
testSet['SHELL_LINKPROCLAMATION_ORDER_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su LP 0',
			'style': '1',
		},
	},
];
testSet['SHELL_LINKPROCLAMATION_ROLE_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su LP 1',
			'style': '1',
		},
	},
];
testSet['SHELL_LINKPROCLAMATION_SHELL_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su LP',
			'style': '1',
		},
	},
];
testSet['SHELL_LINKPROCLAMATION_TITLE_TOOLONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su LP 1 -t ........10~~~',
			'style': '1',
		},
	},
];
testSet['SHELL_NORMALPROCLAMATION_CONTENT_TOOLONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su NP -c ........10........20........30........40........50~~~',
			'style': '1',
		},
	},
];
testSet['SHELL_NORMALPROCLAMATION_ROLE_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00bd6-f384-766f-e48b-997e-becc2787',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su NP',
			'style': '1',
		},
	},
];
testSet['SHELL_NORMALPROCLAMATION_SHELL_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su NP 1',
			'style': '1',
		},
	},
];
testSet['SHELL_NORMALPROCLAMATION_TITLE_TOOLONG'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su NP -t ........10~~~',
			'style': '1',
		},
	},
];
testSet['SHELL_QUERYBLOCKLIST_PARAMETER_ERROR'] = [
	{
		'cmd': '2',
		'idem': '',
		'payload': {
			'platformUuid': '5db00be4-b7fa-6e69-46e7-866a-42f161c0',
			'platform': 'MM',
		},
	},
	{
		'cmd': '84',
		'idem': '',
		'payload': {
			'chatTarget': '000572bf8d4a5001',
			'message': '/su BL 1',
			'style': '1',
		},
	},
];
testSet['SHELL_SHELL_SHELL_ERROR'] = [
	{
		"cmd": "2",
		"idem": "",
		"payload": {
			"platformUuid": "5db00be4-b7fa-6e69-46e7-866a-42f161c0",
			"platform": "MM",
		},
	},
	{
		"cmd": "84",
		"idem": "",
		"payload": {
			"chatTarget": "000572bf8d4a5001",
			"message": "/su",
			"style": "1",
		},
	},
];

/** @type {(selector: string) => HTMLInputElement} */
function $(selector) {
	return document.querySelector(selector);
};

/** @type {(selector: string) => NodeListOf<HTMLInputElement>} */
function $$(selector) {
	return document.querySelectorAll(selector);
};

/** @type {(obj: any) => string} */
function jsonPretty(obj) {
	return JSON.stringify(obj,null,2);
};

//================ 功能 ================
!function() {
	/** @type {(message: string)} */
	function output(message) {
		$('#SocketOutput').insertAdjacentText('beforeend',message+'\n\n');
		$('#OutputBox').scrollTop = $('#OutputBox').scrollHeight;
	};
	/** @type {(isDisable: bool)} */
	function disableInput(isDisable) {
		$('#SocketInput').disabled = isDisable;
		$('#SendBtn').disabled = isDisable;
	};
	/** @type {(this: HTMLInputElement, ev: MouseEvent)} */
	function onSendBtnClick(ev) {
		$('#SocketOutput').innerHTML = '';
		/** @type {any[]} */
		var packetList;
		try {
			packetList = JSON.parse($('#SocketInput').value);
		} catch(e) {
			output('封包列表 JSON 解析失敗');
			return;
		};
		if (!Array.isArray(packetList)) {
			output('封包列表必須是陣列');
			return;
		};
		disableInput(true);
		var ws = new WebSocket('wss://mml.zanxingbctv.com/echo');
		var packetCount = 0;
		var timer = 0;
		function closeWebSocket() {
			if (ws) {
				ws.close();
				ws = null;
				output('======== 關閉 WebSocket ========');
				disableInput(false);
			};
		};
		function onPacketTimeout() {
			output('======== WebSocket 沒有回應 ========');
			closeWebSocket();
		};
		function sendPacket() {
			if (packetCount < packetList.length) {
				ws.send(JSON.stringify(packetList[packetCount]));
				output('======== 發送封包 ' + packetCount + ' ========');
				++packetCount;
				timer = setTimeout(onPacketTimeout,2000);
			} else {
				output('======== 封包發送完畢 ========');
				closeWebSocket();
			};
		};
		ws.onopen = function(ev) {
			output('======== WebSocket 連線成功 ========');
			sendPacket();
		};
		ws.onerror = function(ev) {
			output('======== WebSocket 連線失敗 ========');
			closeWebSocket();
		};
		ws.onmessage = function(ev) {
			var packet = JSON.parse(ev.data);
			output(jsonPretty(packet));
			clearTimeout(timer);
			timer = setTimeout(sendPacket,400);
		};
	};
	$('#SendBtn').onclick = onSendBtnClick;
}();

//================ 畫面 ================
!function() {
	/** @type {(this: HTMLElement, ev: MouseEvent)} */
	function onVrMouseDown(ev) {
		ev.stopImmediatePropagation();
		ev.preventDefault();
		this.dataset.reisze = ev.x - this.previousElementSibling.offsetWidth;
		return false;
	};
	/** @type {(this: HTMLElement, ev: MouseEvent)} */
	function onHrMouseDown(ev) {
		ev.stopImmediatePropagation();
		ev.preventDefault();
		this.dataset.reisze = ev.y - this.previousElementSibling.offsetHeight;
		return false;
	};
	/** @type {(this: HTMLElement, ev: MouseEvent)} */
	function onVrMouseMove(ev) {
		if (this.matches(':scope:active')) {
			this.previousElementSibling.style.width = (ev.x - this.dataset.reisze) + 'px';
		};
	};
	/** @type {(this: HTMLElement, ev: MouseEvent)} */
	function onHrMouseMove(ev) {
		if (this.matches(':scope:active')) {
			this.previousElementSibling.style.height = (ev.y - this.dataset.reisze) + 'px';
		};
	};
	$$('.frameset.row>hr').forEach(function(value,key,parent) {
		value.onmousedown = onVrMouseDown;
		value.onmousemove = onVrMouseMove;
	});
	$$('.frameset.col>hr').forEach(function(value,key,parent) {
		value.onmousedown = onHrMouseDown;
		value.onmousemove = onHrMouseMove;
	});
	/** @type {(this: HTMLElement, ev: MouseEvent)} */
	function onItemClick(ev) {
		var message = this.dataset.value;
		$('#SocketInput').value = jsonPretty(testSet[message]);
		console.log(message);
	};
	Object.keys(testSet).sort().forEach(function(value,index,array) {
		var item = document.createElement('div');
		item.className = 'item';
		item.dataset.value = value;
		item.textContent = value;
		item.onclick = onItemClick;
		$('#Menu').appendChild(item);
	});
}();
