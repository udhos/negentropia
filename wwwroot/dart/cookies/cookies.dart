library cookies;

import 'dart:html';

const ONEDAY_MILLISEC = 24 * 60 * 60 * 1000;

class Cookie {
  static Map<String, String> _readCookie() {
    Map<String, String> cookie = new Map<String, String>();
    String c = document.cookie;
    c.split(';').forEach((t) {
      int k = t.indexOf('=');
      if (k > 0) {
        cookie[Uri.decodeComponent(t.substring(0, k))] =
            Uri.decodeComponent(t.substring(k + 1));
      }
    });
    return cookie;
  }

  static void _writeCookie(Map<String, String> m) {
    String sb;
    Iterator<String> t = m.keys.iterator;
    if (!t.moveNext()) {
      String k = t.current;
      sb = '${Uri.encodeComponent(k)}=${Uri.encodeComponent(m[k])}';
      while (t.moveNext()) {
        k = t.current;
        sb = '${sb};${Uri.encodeComponent(k)}=${Uri.encodeComponent(m[k])}';
      }
    }
    document.cookie = sb.toString();
  }

  static void setCookie(String name, String value, int days) {
    Map<String, String> t = _readCookie();
    t[name] = value;

    DateTime now = new DateTime.now();
    DateTime date = new DateTime.fromMillisecondsSinceEpoch(
        now.millisecondsSinceEpoch + days * ONEDAY_MILLISEC);
    t['expires'] = date.toString();

    _writeCookie(t);
  }

  static String getCookie(String name) {
    Map<String, String> t = _readCookie();
    if (t.containsKey(name)) return t[name];
    return null;
  }
}
