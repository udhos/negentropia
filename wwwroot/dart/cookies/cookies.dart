library cookies;

import 'dart:html';
import 'dart:uri';

class Cookie {
  static _readCookie() {
    var cookie = new Map();
    var c = document.cookie;
    c.split(';').forEach((t) {
      var k = t.indexOf('=');
      if (k > 0)
      {
        cookie[decodeUriComponent(t.substring(0, k))] = decodeUriComponent(t.substring(k + 1));
      }
    });
    return cookie;
  }

  static _writeCookie(Map m) {
    String sb;
    var t = m.keys.iterator;
    if (!t.moveNext()) {
      var k = t.current;
      sb = '${encodeUriComponent(k)}=${encodeUriComponent(m[k])}';
      while (t.moveNext()) {
        k = t.current;
        sb = '${sb};${encodeUriComponent(k)}=${encodeUriComponent(m[k])}';
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



