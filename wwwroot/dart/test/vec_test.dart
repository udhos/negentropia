import 'dart:math' as math;

import 'package:unittest/unittest.dart';
import 'package:vector_math/vector_math.dart';

import '../vec.dart';
import '../string.dart';

void main() {
  vec_test();
  quat_test();
  string_test();
}

typedef void echo_func(String);

void vec_test() {
  echo_func echo = null;

  void testParseVector3(var expected, String target, String str) {
    test("parseVector3: '$str' => '$target'",
        () => expect(parseVector3(str, echoFunc: echo).toString(), expected));
  }

  List<String> patterns1 = [
    "1.0,2.0,3.0",
    " 1.0 , 2.0 , 3.0 ",
    "1,2,3",
    " 1 , 2 , 3 "
  ];

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
    test("parseVector3: '$str' => null", () => expect(
        parseVector3(str, echoFunc: echo).toString(), equals(null.toString())));
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
    "1,2,3,4"
  ];

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
    "1,2,3]["
  ];

  bad_patterns2.forEach((p) => testParseVector3bad(p));
}

void quat_test() {
  double radAngle = math.PI;
  Vector3 axis = new Vector3.all(1.0);
  axis.normalize();
  Quaternion quat = new Quaternion.axisAngle(axis, radAngle);
  test("quat_test: 90deg around (1,1,1)", () {
    expect(quat.toString(), equals(
        "0.5773502588272095, 0.5773502588272095, 0.5773502588272095 @ 6.123234262925839e-17"));
  });

  Vector3 vec = new Vector3(1.0, 0.0, 0.0);
  quat.rotate(vec);
  test("quat_test: quat applied to vec=(1,0,0)", () {
    expect(vec.toString(),
        equals("[-0.3333333134651184,0.6666666269302368,0.6666666269302368]"));
  });
}

void string_test() {
  test('stringIsTrue(null)', () {
    expect(stringIsTrue(null), isFalse);
  });
  test('stringIsTrue("")', () {
    expect(stringIsTrue(""), isFalse);
  });
  test('stringIsTrue(" ")', () {
    expect(stringIsTrue(" "), isFalse);
  });
  test('stringIsTrue(" 1 ")', () {
    expect(stringIsTrue(" 1 "), isTrue);
  });
  test('stringIsTrue(" 0.1 ")', () {
    expect(stringIsTrue(" 0.1 "), isTrue);
  });
  test('stringIsTrue(" 0 ")', () {
    expect(stringIsTrue(" 0 "), isFalse);
  });
  test('stringIsTrue(" 05 ")', () {
    expect(stringIsTrue(" 05 "), isTrue);
  });
  test('stringIsTrue(" f ")', () {
    expect(stringIsTrue(" f "), isFalse);
  });
  test('stringIsTrue(" false ")', () {
    expect(stringIsTrue(" false "), isFalse);
  });
  test('stringIsTrue(" F ")', () {
    expect(stringIsTrue(" F "), isFalse);
  });
  test('stringIsTrue(" o ")', () {
    expect(stringIsTrue(" o "), isTrue);
  });
  test('stringIsTrue(" of ")', () {
    expect(stringIsTrue(" of "), isFalse);
  });
  test('stringIsTrue(" off ")', () {
    expect(stringIsTrue(" off "), isFalse);
  });
  test('stringIsTrue(" OFF ")', () {
    expect(stringIsTrue(" OFF "), isFalse);
  });
  test('stringIsTrue(" on ")', () {
    expect(stringIsTrue(" on "), isTrue);
  });
  test('stringIsTrue(" 0XXX45 ")', () {
    expect(stringIsTrue(" 0XXX45 "), isTrue);
  });
}
