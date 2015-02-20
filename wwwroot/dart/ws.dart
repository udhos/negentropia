library ws;

import 'dart:html';
import 'dart:async';
import 'dart:convert';
import 'dart:collection';

import 'logg.dart';

const CM_CODE_FATAL = 0;
const CM_CODE_INFO = 1;
const CM_CODE_AUTH = 2; // client->server: let me in
const CM_CODE_ECHO = 3; // client->server: please echo this
const CM_CODE_KILL = 4;
// server->client: do not attempt reconnect on same session
const CM_CODE_REQZ = 5; // client->server: please send current zone
const CM_CODE_ZONE = 6; // server->client: reset client zone info
const CM_CODE_SKYBOX = 7; // server->client: set full skybox
const CM_CODE_PROGRAM = 8; // server->client: set shader program
const CM_CODE_INSTANCE = 9; // server->client: set instance
const CM_CODE_INSTANCE_UPDATE = 10; // server->client: update instance
const CM_CODE_MESSAGE = 11; // server->client: message for user
const CM_CODE_MISSION_NEXT = 12; // client->server: switch mission
const CM_CODE_SWITCH_ZONE = 13; // client->server: switch zone

WebSocket _ws;
ListQueue<String> _wsQueue = new ListQueue<String>();
typedef void dispatcherFunc(int code, String data, Map<String, String> tab);
dispatcherFunc _dispatcher;

void missionNext(Map m) {
  wsSendMap({'Code': CM_CODE_MISSION_NEXT, 'Data': "", 'Tab': m});
}

void switchZone() {
  wsSendMap({'Code': CM_CODE_SWITCH_ZONE, 'Data': ""});
}

void requestZone() {
  /*
  Map msg = new Map();  
  msg["Code"] = CM_CODE_REQZ;
  msg["Data"] = "";
  */
  Map msg = {'Code': CM_CODE_REQZ, 'Data': ""};

  wsSendMap(msg);
}

void _write(String msg) {
  _ws.send(msg);
}

void wsSendString(String msg) {
  _wsQueue.add(msg);
  wsFlush();
}

void wsSendMap(Map msg) {
  wsSendString(JSON.encode(msg));
}

void wsFlush() {
  while (_ws != null && _ws.readyState == WebSocket.OPEN && !_wsQueue.isEmpty) {
    try {
      _write(_wsQueue.first);
    } catch (e) {
      err("websocket flush: send failure: $e");
      return;
    }
    _wsQueue.removeFirst();
  }
}

void initWebSocket(String wsUri, String sid, int retrySeconds, Element status,
    dispatcherFunc dispatch) {
  _dispatcher = dispatch;

  status.text = "opening $wsUri";

  if (retrySeconds < 1) {
    retrySeconds = 1;
  } else if (retrySeconds > 20) {
    retrySeconds = 20;
  }

  debug("websocket: opening: ${wsUri} (retry=${retrySeconds})");

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

    debug("websocket: retrying in $retrySeconds seconds");
    new Timer(new Duration(seconds: retrySeconds),
        () => initWebSocket(wsUri, sid, 2 * retrySeconds, status, dispatch));

    reconnectScheduled = true;
  }

  subOpen = _ws.onOpen.listen((e) {
    status.text = "connected to $wsUri";
    debug("websocket: ${status.text}");

    /*
    var msg = new Map();
    msg["Code"] = CM_CODE_AUTH;
    msg["Data"] = sid;
    */
    Map msg = {'Code': CM_CODE_AUTH, 'Data': sid};

    String jsonMsg = JSON.encode(msg);

    try {
      _write(jsonMsg);
    } catch (e) {
      err("websocket auth: send failure: $e");
      scheduleReconnect();
    }
  });

  subClose = _ws.onClose.listen((Event e) {
    status.text = "disconnected from $wsUri";
    warn("websocket: ${status.text}: [$e]");
    scheduleReconnect();
  });

  subError = _ws.onError.listen((Event e) {
    err("websocket: error: [$e]");
    scheduleReconnect();
  });

  subMessage = _ws.onMessage.listen((MessageEvent e) {
    Map msg = JSON.decode(e.data);
    int code = msg["Code"];
    String data = msg["Data"];
    Map<String, String> tab = msg["Tab"];

    if (code == CM_CODE_KILL) {
      String killInfo = data;
      String m = "server killed our session: $killInfo";

      warn(m);
      status.text = m;

      subOpen.cancel();
      subClose.cancel();
      subMessage.cancel();
      subError.cancel();
      _ws.close();
      _ws = null;

      return;
    }

    _dispatcher(code, data, tab);
  });
}
