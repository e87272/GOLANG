/**
 * @typedef TestCase
 * @property {string} Method
 * @property {Object.<string, string>} Title
 * @property {Object.<string, any>} Input
 * @property {any} Output
 */

/**
 * @typedef InputPacket
 * @property {string} api
 * @property {string} method
 * @property {any} body
 */

/**
 * @typedef OutputPacket
 * @property {string} result
 * @property {Object.<string, string>} error
 * @property {any} payload
 */

/** @type {TestCase[]} */
var testSet;

/** @type {(selector: string) => HTMLInputElement} */
function $(selector) {
	return document.querySelector(selector);
};

/** @type {(selector: string) => NodeListOf<HTMLInputElement>} */
function $$(selector) {
	return document.querySelectorAll(selector);
};

/** @type {(x: any)=>string} */
function typeOf(x) {
	return Object.prototype.toString.call(x).slice(8, -1);
};

/** @type {(json: string)=>any} */
function jsonEval(json) {
	return Function('return (\n' + json + '\n);')();
};

/** @type {(x: any, space: string, pad: string)=>string} */
function jsonPretty(x, space = '  ', pad = '') {
	let content = '';
	let indent = pad + space;
	switch (typeOf(x)) {
		case 'Array': {
			/** @type {any[]} */
			let a = x;
			for (let i = 0; i < a.length; ++i) {
				let y = a[i];
				content += indent + jsonPretty(y, space, indent) + ',\n';
			};
			return content == '' ? '[]' : '[\n' + content + pad + ']';
		};
		case 'Object': {
			let a = Object.keys(x).sort();
			for (let i = 0; i < a.length; ++i) {
				let key = a[i];
				let y = x[key];
				content += indent + JSON.stringify(key) + ': ' + jsonPretty(y, space, indent) + ',\n';
			};
			return content == '' ? '{}' : '{\n' + content + pad + '}';
		};
		default: {
			return JSON.stringify(x);
		};
	};
};

/** @type {(x: any, comment: any, space: string, pad: string)=>string} */
function jsonComment(x, comment, space = '  ', pad = '') {
	let type = typeOf(x);
	if (typeOf(comment) != type) {
		return jsonPretty(x, space, pad);
	};
	let content = '';
	let indent = pad + space;
	switch (typeOf(x)) {
		case 'Array': {
			/** @type {any[]} */
			let a = x;
			for (let i = 0; i < a.length; ++i) {
				let y = a[i];
				let z = comment[i];
				content += indent + jsonComment(y, z, space, indent);
				content += (typeOf(z) == 'String' ? ',  //' + z : ',') + '\n';
			};
			return content == '' ? '[]' : '[\n' + content + pad + ']';
		};
		case 'Object': {
			let a = Object.keys(x).sort();
			for (let i = 0; i < a.length; ++i) {
				let key = a[i];
				let y = x[key];
				let z = comment[key];
				content += indent + JSON.stringify(key) + ': ' + jsonComment(y, z, space, indent);
				content += (typeOf(z) == 'String' ? ',  //' + z : ',') + '\n';
			};
			return content == '' ? '{}' : '{\n' + content + pad + '}';
		};
		default: {
			return JSON.stringify(x);
		};
	};
};

/** @type {(x: any)=>any} */
function cloneStruct(x) {
	let type = typeOf(x);
	if (type != 'Object') {
		return (new self[type]());
	};
	let y = {};
	for (let i in x) {
		y[i] = cloneStruct(x[i]);
	};
	return y;
};

//================ 功能 ================

/** @type {(message: string)} */
function output(message) {
	$('#Output').insertAdjacentText('beforeend', message + '\n\n');
	$('#OutputBox').scrollTop = 0;
};

/** @type {(isDisable: boolean)} */
function disableInput(isDisable) {
	$('#Input').disabled = isDisable;
	$('#SendBtn').disabled = isDisable;
};

/** @type {(this: HTMLElement, ev: MouseEvent)} */
function onItemClick(ev) {
	let testCase = testSet[+this.dataset.index];
	let api = this.dataset.api;
	/** @type {InputPacket} */
	let packet = {
		'api': api,
		'method': testCase.Method,
		'body': cloneStruct(testCase.Input),
	};
	let json = '//' + testCase.Title[api] + '\n' + jsonComment(packet, { 'body': testCase.Input });
	$('#Input').value = json;
};

/** @type {(this: HTMLInputElement, ev: MouseEvent)} */
function onSendBtnClick(ev) {
	$('#Output').innerHTML = '';
	/** @type {InputPacket} */
	let inputPacket;
	try {
		inputPacket = jsonEval($('#Input').value);
	} catch (e) {
		output('JSON 解析失敗');
		return;
	};
	if (inputPacket == null || inputPacket.body == null) {
		output('JSON 解析失敗');
		return;
	}
	disableInput(true);
	console.log(inputPacket);
	let bodyStr = '';
	Object.keys(inputPacket.body).sort().forEach(function (value, index, array) {
		bodyStr += '&' + encodeURIComponent(value) + '=' + encodeURIComponent(inputPacket.body[value]);
	});
	bodyStr = bodyStr.slice(1);
	let xhr = new XMLHttpRequest();
	xhr.open(inputPacket.method, document.location.origin + inputPacket.api, true);
	xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	xhr.onerror = function () {
		output('API 請求失敗');
		xhr = null;
		disableInput(false);
	};
	xhr.onload = function () {
		/** @type {OutputPacket} */
		let outputPacket;
		try {
			outputPacket = JSON.parse(xhr.responseText);
		} catch (e) {
			outputPacket = null;
		};
		console.log(outputPacket);
		if (outputPacket) {
			let item = $('[data-api="' + inputPacket.api + '"]');
			output(item ? jsonComment(outputPacket, { 'payload': testSet[+item.dataset.index].Output }) : jsonPretty(outputPacket));
		} else {
			output(xhr.responseText);
		}
		xhr = null;
		disableInput(false);
	};
	xhr.send(bodyStr);
};

//================ 畫面 ================

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

$$('.frameset.row>hr').forEach(function (value, key, parent) {
	value.onmousedown = onVrMouseDown;
	value.onmousemove = onVrMouseMove;
});

$$('.frameset.col>hr').forEach(function (value, key, parent) {
	value.onmousedown = onHrMouseDown;
	value.onmousemove = onHrMouseMove;
});

testSet.forEach(function (value, index, array) {
	let testIndex = index;
	Object.keys(value.Title).sort().forEach(function (value, index, array) {
		let item = document.createElement('div');
		item.className = 'item';
		item.dataset.index = testIndex;
		item.dataset.api = value;
		item.textContent = value;
		item.onclick = onItemClick;
		$('#Menu').appendChild(item);
	});
});

$('#SendBtn').onclick = onSendBtnClick;
$('#Input').value = '';
$('#Output').innerHTML = '';
