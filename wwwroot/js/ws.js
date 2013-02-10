
var CM_CODE_FATAL = 0;
var	CM_CODE_INFO  = 1;
var	CM_CODE_AUTH  = 2;
	
function initWebSocket() {
	var wsUri = "ws://127.0.0.2:8000/";
	console.log("websocket: opening " + wsUri);
	websocket = new WebSocket(wsUri);
	websocket.onopen = function(evt) { onOpen(evt) };
	websocket.onclose = function(evt) { onClose(evt) };
	websocket.onmessage = function(evt) { onMessage(evt) };
	websocket.onerror = function(evt) { onError(evt) };
}

function onOpen(evt) {
	console.log("websocket: CONNECTED");
	
	var msg = {
		Code: CM_CODE_AUTH,
		Data: sid
	};
  
	doSend(JSON.stringify(msg));
}

function onClose(evt) {
	console.log("websocket: DISCONNECTED");
}

function onMessage(evt) {
	console.log("websocket: received: [" + evt.data + "]");
}

function onError(evt) {
	console.log("websocket: error: [" + evt.data + "]");
}

function doSend(message) {
	console.log("websocket: sending: [" + message + "]");
	websocket.send(message);
}
