library ws;

import 'dart:html';
import 'dart:async';
import 'dart:json';

const CM_CODE_FATAL = 0;
const CM_CODE_INFO  = 1;
const CM_CODE_AUTH  = 2; // client->server: let me in
const CM_CODE_ECHO  = 3; // client->server: please echo this
const CM_CODE_KILL  = 4; // server->client: do not attempt reconnect on same session

/*
WebSocket w;
StreamSubscription<Event> subOpen;
StreamSubscription<Event> subClose;
StreamSubscription<Event> subError;
StreamSubscription<Event> subMessage;
*/

void doSend(WebSocket w, String msg) {
  print("websocket: sending: [${msg}]");
  w.send(msg);
}

void initWebSocket(String wsUri, String sid, int retrySeconds, Element status) {
  
  status.text = "opening $wsUri";
  
  if (retrySeconds < 1) {
    retrySeconds = 1;    
  } 
  else if (retrySeconds > 120) {
    retrySeconds = 120;
  }
  
  print("websocket: opening: ${wsUri} (retry=${retrySeconds})");
  
  WebSocket w = new WebSocket(wsUri);

  StreamSubscription<Event> subOpen;
  StreamSubscription<Event> subClose;
  StreamSubscription<Event> subError;
  StreamSubscription<Event> subMessage;

  subOpen = w.onOpen.listen((e) {
    status.text = "connected to $wsUri";   
    print("websocket: CONNECTED");

    var msg = new Map();
    msg["Code"] = CM_CODE_AUTH;
    msg["Data"] = sid;
    
    String jsonMsg = stringify(msg);
    
    doSend(w, jsonMsg);
  });
  
  bool reconnectScheduled = false;
  
  void scheduleReconnect() {
    if (reconnectScheduled) {
      return;
    }
    
    print("websocket: retrying in $retrySeconds seconds");
    new Timer(new Duration(seconds: retrySeconds), () => initWebSocket(wsUri, sid, 2 * retrySeconds, status));
      
    reconnectScheduled = true;
  }
  
  subClose = w.onClose.listen((MessageEvent e) {
    status.text = "disconnected from $wsUri";    
    print("websocket: DISCONNECTED");
    scheduleReconnect();
  });
  
  subError = w.onError.listen((MessageEvent e) {
    print("websocket: error: w.onError.listen");
    print("websocket: error: [${e.data}]");
    scheduleReconnect();
  });
  
  subMessage = w.onMessage.listen((MessageEvent e) {
    print('websocket: received: w.onMessage.listen');
    print('websocket: received: [${e.data}]');
    
    Map msg = parse(e.data);
    
    if ((msg["Code"] == CM_CODE_INFO) && (msg["Data"].startsWith("welcome"))) {
      // test echo loop thru server
      var m = new Map();
      m["Code"] = CM_CODE_ECHO;
      m["Data"] = "hi there";
      doSend(w, stringify(m));
      return;
    }
    
    if (msg["Code"] == CM_CODE_KILL) {
      
      String killInfo = msg["Data"];
      String m = "server killed our session: $killInfo";

      print(m);
      status.text = m;

      subOpen.cancel();
      subClose.cancel();
      subMessage.cancel();
      subError.cancel();
      w.close();
      w = null;
      
      return;
    }
    
  });
}

