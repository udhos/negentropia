library cookies;

import 'dart:html';

class Cookie {
  static _readCookie() {
    var cookie = new Map();
    var c = document.cookie;
    c.split(';').forEach((t) {
      var k = t.indexOf('=');
      if (k > 0) {
        cookie[Uri.decodeComponent(t.substring(0, k))] = Uri.decodeComponent(t.substring(k + 1));
      }
    });
    return cookie;
  }

  static _writeCookie(Map m) {
    String sb;
    var t = m.keys.iterator;
    if (!t.moveNext()) {
      var k = t.current;
      sb = '${Uri.encodeComponent(k)}=${Uri.encodeComponent(m[k])}';
      while (t.moveNext()) {
        k = t.current;
        sb = '${sb};${Uri.encodeComponent(k)}=${Uri.encodeComponent(m[k])}';
      }
    }
    document.cookie = sb.toString();
  }

  static void setCookie(String name, String value, int days) {
    var t = _readCookie();
    t[name] = value;

    DateTime now = new DateTime.now();
    DateTime date = new DateTime.fromMillisecondsSinceEpoch(now.millisecondsSinceEpoch + days*24*60*60*1000);
    t['expires'] = date.toString();

    _writeCookie(t);
  }

  static String getCookie(String name) {
    var t = _readCookie();
    if (t.containsKey(name))
      return t[name];
    return null;
  }
}



