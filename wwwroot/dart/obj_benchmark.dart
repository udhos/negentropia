//import 'dart:html';
import 'dart:io';
import 'dart:async';

import 'package:benchmark_harness/benchmark_harness.dart';

import 'asset.dart';
import 'obj.dart';

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

  // Not measured setup code executed prior to the benchmark runs.
  void setup() { }

  // Not measures teardown code executed after the benchark runs.
  void teardown() { }
}

class MtlBenchmark extends BenchmarkBase {
  const MtlBenchmark() : super("Benchmark: mtllib_parse");

  static void main() {
    new MtlBenchmark().report();
  }

  void run() {
    Map<String,Material> lib = mtllib_parse(mtlString, mtlURL);
  }

  // Not measured setup code executed prior to the benchmark runs.
  void setup() { }

  // Not measures teardown code executed after the benchark runs.
  void teardown() { }
}

void main() {

  void fetchMtl(String URL) {
    
    void done(String response) {
      mtlURL    = URL;
      mtlString = response;

      // Run TemplateBenchmark
      MtlBenchmark.main();
    }
            
    var file = new File(URL);    
    Future<String> finishedReading = file.readAsString(encoding: Encoding.ASCII);
    finishedReading.then(done); 
  }
  
  void fetchObj(String URL) {
    
    void done(String response) {
      objURL    = URL;
      objString = response;

      // Run TemplateBenchmark
      ObjBenchmark.main();
      
      fetchMtl("C:/tmp/devel/negentropia/wwwroot/mtl/Colony Ship Ogame Fleet.mtl");
    }
            
    var file = new File(URL);    
    Future<String> finishedReading = file.readAsString(encoding: Encoding.ASCII);
    finishedReading.then(done); 
  }

  fetchObj("C:/tmp/devel/negentropia/wwwroot/obj/Colony Ship Ogame Fleet.obj");    
}
