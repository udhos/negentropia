library ws;

import 'dart:html';
import 'dart:async';
import 'dart:json';
import 'dart:collection';

const CM_CODE_FATAL = 0;
const CM_CODE_INFO  = 1;
const CM_CODE_AUTH  = 2; // client->server: let me in
const CM_CODE_ECHO  = 3; // client->server: please echo this
const CM_CODE_KILL  = 4; // server->client: do not attempt reconnect on same session
const CM_CODE_REQZ  = 5; // client->server: please send current zone
const CM_CODE_ZONE  = 6; // server->client: current zone
const CM_CODE_SKYBOX   = 7;  // server->client: set full skybox
const CM_CODE_PROGRAM  = 8;  // server->client: set shader program
const CM_CODE_MODEL    = 9;  // server->client: set model
const CM_CODE_INSTANCE = 10; // server->client: set instance

WebSocket _ws;
ListQueue<String> _wsQueue = new ListQueue<String>();
typedef void dispatcherFunc(int code, String data);
dispatcherFunc _dispatcher;

void requestZone() {
  Map msg = new Map();
  
  msg["Code"] = CM_CODE_REQZ;
  msg["Data"] = "";
  
  String json = stringify(msg);
  
  wsSend(json);  
}

void _write(String msg) {
  //print("websocket: writing: [${msg}]");
  _ws.send(msg);
}

void wsSend(String msg) {
  _wsQueue.add(msg);
  wsFlush();
}

void wsFlush() {
  while (_ws != null && _ws.readyState == WebSocket.OPEN && !_wsQueue.isEmpty) {
    try {
      _write(_wsQueue.first);
    }
    catch(e) {
      print("websocket flush: send failure: $e");
      return;
    }
    _wsQueue.removeFirst();
  }
}

void initWebSocket(String wsUri, String sid, int retrySeconds, Element status, dispatcherFunc dispatch) {
  
  _dispatcher = dispatch;
  
  status.text = "opening $wsUri";
  
  if (retrySeconds < 1) {
    retrySeconds = 1;    
  } 
  else if (retrySeconds > 120) {
    retrySeconds = 120;
  }
  
  print("websocket: opening: ${wsUri} (retry=${retrySeconds})");
  
  _ws = new WebSocket(wsUri);

  StreamSubscription<Event> subOpen;
  StreamSubscription<Event> subClose;
  StreamSubscription<Event> subError;
  StreamSubscription<Event> subMessage;
  
  bool reconnectScheduled = false;
  
  void scheduleReconnect() {
    if (reconnectScheduled) {
      return;
    }
    
    print("websocket: retrying in $retrySeconds seconds");
    new Timer(new Duration(seconds: retrySeconds), () => initWebSocket(wsUri, sid, 2 * retrySeconds, status, dispatch));
      
    reconnectScheduled = true;
  }

  subOpen = _ws.onOpen.listen((e) {
    status.text = "connected to $wsUri";   
    print("websocket: ${status.text}");

    var msg = new Map();
    msg["Code"] = CM_CODE_AUTH;
    msg["Data"] = sid;
    
    String jsonMsg = stringify(msg);
    
    try {
      _write(jsonMsg);
    }
    catch (e) {
      print("websocket auth: send failure: $e");
      scheduleReconnect();
    }
  });
  
  subClose = _ws.onClose.listen((Event e) {
    status.text = "disconnected from $wsUri";    
    print("websocket: DISCONNECTED");
    scheduleReconnect();
  });
  
  subError = _ws.onError.listen((Event e) {
    print("websocket: error: [$e]");
    scheduleReconnect();
  });
  
  subMessage = _ws.onMessage.listen((MessageEvent e) {
    //print('websocket: received: [${e.data}]');
    
    Map msg = parse(e.data);
    int code = msg["Code"];
    String data = msg["Data"];
    
    if (code == CM_CODE_KILL) {
      
      String killInfo = data;
      String m = "server killed our session: $killInfo";

      print(m);
      status.text = m;

      subOpen.cancel();
      subClose.cancel();
      subMessage.cancel();
      subError.cancel();
      _ws.close();
      _ws = null;
      
      return;
    }
    
    _dispatcher(code, data);
  });
}

