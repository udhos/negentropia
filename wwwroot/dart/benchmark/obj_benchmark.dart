import 'dart:io';
import 'dart:async';
import 'dart:convert';

import 'package:benchmark_harness/benchmark_harness.dart';
import 'package:obj/obj.dart';

String objURL;
String objString;
String mtlURL;
String mtlString;

class ObjBenchmark extends BenchmarkBase {
  const ObjBenchmark() : super("Benchmark: Obj.fromString");

  static void main() {
    new ObjBenchmark().report();
  }

  void run() {
    Obj obj = new Obj.fromString(objURL, objString);
  }
}

class MtlBenchmark extends BenchmarkBase {
  const MtlBenchmark() : super("Benchmark: mtllib_parse");

  static void main() {
    new MtlBenchmark().report();
  }

  void run() {
    Map<String, Material> lib = mtllib_parse(mtlString, mtlURL);
  }
}

void main() {
  void fetchMtl(String URL) {
    void done(String response) {
      mtlURL = URL;
      mtlString = response;
      MtlBenchmark.main(); // run benchmark
    }

    var file = new File(URL);
    Future<String> finishedReading = file.readAsString(encoding: ASCII);
    finishedReading.then(done);
  }

  void fetchObj(String URL) {
    void done(String response) {
      objURL = URL;
      objString = response;
      ObjBenchmark.main(); // run benchmark

      fetchMtl(
          "C:/tmp/devel/negentropia/wwwroot/mtl/Colony Ship Ogame Fleet.mtl");
    }

    var file = new File(URL);
    Future<String> finishedReading = file.readAsString(encoding: ASCII);
    finishedReading.then(done);
  }

  fetchObj("C:/tmp/devel/negentropia/wwwroot/obj/Colony Ship Ogame Fleet.obj");
}
