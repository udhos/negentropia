
var CM_CODE_FATAL = 0;
var	CM_CODE_INFO  = 1;
var	CM_CODE_AUTH  = 2;

var wsUri = "ws://127.0.0.2:8000/";
	
function initWebSocket(status) {
	console.log("websocket: opening " + wsUri);
	status.innerHTML = "opening " + wsUri;
	websocket = new WebSocket(wsUri);
	websocket.onopen = function(evt) { onOpen(evt, status) };
	websocket.onclose = function(evt) { onClose(evt, status) };
	websocket.onmessage = function(evt) { onMessage(evt, status) };
	websocket.onerror = function(evt) { onError(evt, status) };
}

function onOpen(evt, status) {
	status.innerHTML = "connected to " + wsUri;
	console.log("websocket: CONNECTED");
	
	var msg = {
		Code: CM_CODE_AUTH,
		Data: sid
	};
  
	doSend(JSON.stringify(msg));
}

function onClose(evt, status) {
	status.innerHTML = "disconnected from " + wsUri;
	console.log("websocket: DISCONNECTED");
}

function onMessage(evt, status) {
	console.log("websocket: received: [" + evt.data + "]");
}

function onError(evt, status) {
	console.log("websocket: error: [" + evt.data + "]");
}

function doSend(message) {
	console.log("websocket: sending: [" + message + "]");
	websocket.send(message);
}
