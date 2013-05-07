library skybox;

class Obj {

  List<double> vertCoord = new List<double>();   
  List<int> vertInd = new List<int>(); 

  Obj.fromString(String url, String str) {
    // TODO FIXME parse OBJ string
    
    int lineNum = 0;
    
    void parseLine(String rawLine) {
      ++lineNum;
      
      //print("line: $lineNum [$rawLine]");
      
      String line = rawLine.trim();
      
      if (line.isEmpty) {
        return;
      }
      
      if (line[0] == '#') {
        return;
      }
      
      if (line.startsWith("v ")) {
        // vertex coord
        List<String> v = line.split(' ');
        if (v.length != 4) {
          print("OBJ: wrong number of vertex coordinates (${v.length - 1} != 3) at line=$lineNum from url=$url: [$rawLine]");
          return;
        }
        vertCoord.add(double.parse(v[1])); // x
        vertCoord.add(double.parse(v[2])); // y
        vertCoord.add(double.parse(v[3])); // z
        return;
      }

      if (line.startsWith("vt ")) {
        // texture coord
        return;
      }

      if (line.startsWith("vn ")) {
        // normal
        return;
      }
      
      if (line.startsWith("f ")) {
        // face
        List<String> f = line.split(' ');
        if (f.length != 4) {
          print("OBJ: wrong number of face indices (${f.length - 1} != 3) at line=$lineNum from url=$url: [$rawLine]");
          return;
        }
        for (int i = 1; i < f.length; ++i) {
          List<String> v = f[i].split('/');
          String vi = v[0];
          int vIndex = int.parse(vi) - 1;
          vertInd.add(vIndex);
          if (v.length > 1) {
            String ti = v[1];
            if (ti != null && !ti.isEmpty) {
              int tIndex = int.parse(ti) - 1;
            }
          }
          if (v.length > 2) {
            String ni = v[2];
            if (ni != null && !ni.isEmpty) {
              int nIndex = int.parse(ni) - 1;
            }
          }
        }
        return;
      }

      print("OBJ: unknown pattern at line=$lineNum from url=$url: [$rawLine]");
    }
    
    List<String> lines = str.split('\n');
    lines.forEach((String line) => parseLine(line));
  }
}