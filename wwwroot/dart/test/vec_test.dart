import 'dart:math' as math;

import 'package:unittest/unittest.dart';
import 'package:vector_math/vector_math.dart';

import '../vec.dart';

void main() {
  vec_test();
  quat_test();
}

typedef void echo_func(String);

void vec_test() {

  echo_func echo = null;

  void testParseVector3(var expected, String target, String str) {
    test(
        "parseVector3: '$str' => '$target'",
        () => expect(parseVector3(str, echoFunc: echo).toString(), expected));
  }

  List<String> patterns1 = [
      "1.0,2.0,3.0",
      " 1.0 , 2.0 , 3.0 ",
      "1,2,3",
      " 1 , 2 , 3 "];

  String target1 = new Vector3(1.0, 2.0, 3.0).toString();
  var expected1 = equals(target1);

  patterns1.forEach((p) => testParseVector3(expected1, target1, p));
  patterns1.forEach((p) => testParseVector3(expected1, target1, "[$p]"));

  List<String> patterns2 = ["1.1,2.2,3.3", " 1.1 , 2.2 , 3.3 ",];

  String target2 = new Vector3(1.1, 2.2, 3.3).toString();
  var expected2 = equals(target2);

  patterns2.forEach((p) => testParseVector3(expected2, target2, p));
  patterns2.forEach((p) => testParseVector3(expected2, target2, "[$p]"));

  void testParseVector3bad(String str) {
    test(
        "parseVector3: '$str' => null",
        () =>
            expect(parseVector3(str, echoFunc: echo).toString(), equals(null.toString())));
  }

  List<String> bad_patterns = [
      null,
      "",
      " ",
      ",",
      ",,",
      ",,,",
      " , ",
      " , , ",
      " , , , " "a",
      "a,a",
      "a,a,a",
      "a,a,a,a",
      " a ",
      " a , a ",
      " a , a , a ",
      " a , a , a , a ",
      "1,a,3",
      "1,2,a",
      "1,2,3,a",
      "1",
      "1,2",
      "1,2,3,4"];

  bad_patterns.forEach((p) => testParseVector3bad(p));
  bad_patterns.forEach((p) => testParseVector3bad("[$p]"));

  List<String> bad_patterns2 = [
      "[[1,2,3]]",
      "[1,2,3",
      "[[1,2,3",
      "1,2,3]",
      "1,2,3]]",
      "][",
      "]][[",
      "[[]]",
      "]1,2,3[",
      "]]1,2,3[[",
      "[]1,2,3",
      "1,2,3[]",
      "][1,2,3",
      "1,2,3]["];

  bad_patterns2.forEach((p) => testParseVector3bad(p));
}

void quat_test() {
  double radAngle = math.PI;
  Vector3 axis = new Vector3.all(1.0);
  axis.normalize();
  Quaternion quat = new Quaternion.axisAngle(axis, radAngle);
  test("quat_test: 90deg around (1,1,1)", () {
    expect(quat.toString(), equals(""));
  });
}
