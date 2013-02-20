library ws;

import 'dart:html';
import 'dart:async';
import 'dart:json';

const CM_CODE_FATAL = 0;
const CM_CODE_INFO  = 1;
const CM_CODE_AUTH  = 2;

WebSocket w;

void doSend(String msg) {
  print("websocket: sending: [${msg}]");
  w.send(msg);
}

void initWebSocket(String wsUri, String sid, int retrySeconds, Element status) {
  
  status.text = "opening $wsUri";
  
  var fail = false;
  
  if (retrySeconds < 1) {
    retrySeconds = 1;    
  } 
  else if (retrySeconds > 120) {
    retrySeconds = 120;
  }
  
  print("websocket: opening: ${wsUri} (retry=${retrySeconds})");
  
  w = new WebSocket(wsUri);
  
  w.onOpen.listen((e) {
    status.text = "connected to $wsUri";   
    print("websocket: CONNECTED");

    var msg = new Map();
    msg["Code"] = CM_CODE_AUTH;
    msg["Data"] = sid;
    
    String jsonMsg = stringify(msg);
    
    doSend(jsonMsg);
  });
  
  w.onClose.listen((e) {
    status.text = "disconnected from $wsUri";    
    print("websocket: DISCONNECTED");
    if (!fail) {
      print("websocket: retrying in $retrySeconds seconds");
      new Timer(1000 * retrySeconds, (Timer t) => initWebSocket(wsUri, sid, 2 * retrySeconds, status));
    }
    fail = true;
  });
  
  w.onError.listen((e) {
    print("websocket: error: [${e.data}]");
    if (!fail) {
      print("websocket: retrying in $retrySeconds seconds");
      new Timer(1000 * retrySeconds, (Timer t) => initWebSocket(wsUri, sid, 2 * retrySeconds, status));
    }
    fail = true;
  });
  
  w.onMessage.listen((MessageEvent e) {
    print('websocket: received: [${e.data}]');
  });
}

