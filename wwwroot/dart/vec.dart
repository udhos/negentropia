library vec;

import 'package:vector_math/vector_math.dart';

import 'logg.dart';

class ParseError {
  final String _label;
  final String _arg;
  final String _msg;
  ParseError(this._label, this._arg, this._msg);
  String toString() {
    return super.toString() + formatParseError(_label, _arg, _msg);
  }
}

String formatParseError(String label, String arg, String msg) {
  return "ParseError: $label($arg): $msg";
}

void _log_error(String msg) {
  err(msg);
}

void _fail(void echo(String), bool exception, String label, String argument,
    String msg) {
  if (echo != null) {
    echo(formatParseError(label, argument, msg));
  }
  if (exception) {
    throw new ParseError(label, argument, msg);
  }
}

Vector3 parseVector3(String s,
    {void echoFunc(String): _log_error, bool throwException: false}) {
  const LABEL = "parseVector3";
  String save = s;

  if (s == null) {
    _fail(echoFunc, throwException, LABEL, save, "null argument");
    return null;
  }

  s = s.trim();

  bool bracket_open = s.startsWith('[');
  bool bracket_close = s.endsWith(']');

  if (bracket_open != bracket_close) {
    _fail(echoFunc, throwException, LABEL, save, "bracket open/close mismatch");
    return null;
  }

  if (bracket_open) {
    s = s.substring(1);
  }

  if (bracket_close) {
    s = s.substring(0, s.length - 1);
  }

  List<String> list = s.split(',');
  if (list.length != 3) {
    _fail(echoFunc, throwException, LABEL, save,
        "bad length=${list.length}: string='$save'");
    return null;
  }

  double x;
  double y;
  double z;

  try {
    x = double.parse(list[0]);
  } catch (e) {
    _fail(echoFunc, throwException, LABEL, save,
        "failure parsing x=${list[0]}: exception: $e");
    return null;
  }
  try {
    y = double.parse(list[1]);
  } catch (e) {
    _fail(echoFunc, throwException, LABEL, save,
        "failure parsing y=${list[1]}: exception: $e");
    return null;
  }
  try {
    z = double.parse(list[2]);
  } catch (e) {
    _fail(echoFunc, throwException, LABEL, save,
        "failure parsing z=${list[2]}: exception: $e");
    return null;
  }

  return new Vector3(x, y, z);
}

const double MAX_CLOSE_TO_ZERO = 1e-7;

bool closeToZero(double d) => d.abs() < MAX_CLOSE_TO_ZERO;

bool vector3Orthogonal(Vector3 v1, Vector3 v2) => closeToZero(v1.dot(v2));

bool vector3Unit(Vector3 v) => closeToZero(v.length - 1.0);
