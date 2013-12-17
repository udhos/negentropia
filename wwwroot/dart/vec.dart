library vec;

import 'package:vector_math/vector_math.dart';

import 'logg.dart';

Vector3 parseVector3(String s) {
  if (s == null) {
    err("parseVector3($s): null");
    return null;
  }

  s = s.trim();
  
  if (s.startsWith('[')) {
    s = s.substring(1);
  }
  
  if (s.endsWith(']')) {
    s = s.substring(0, s.length - 1);
  }

  List<String> list = s.split(',');
  if (list.length != 3) {
    err("parseVector3($s): bad length: ${list.length}");
    return null;
  }
  
  double x;
  double y;
  double z;

  try { x = double.parse(list[0]); }
  catch (e) { err("parseVector3($s): failure parsing x=${list[0]}: exception: $e"); return null; }
  try { y = double.parse(list[1]); }
  catch (e) { err("parseVector3($s): failure parsing y=${list[1]}: exception: $e"); return null; }
  try { z = double.parse(list[2]); }
  catch (e) { err("parseVector3($s): failure parsing z=${list[2]}: exception: $e"); return null; }
  
  return new Vector3(x, y, z);
}
