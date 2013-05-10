library skybox;

class Obj {
  
  List<double> vertCoord = new List<double>();   
  List<int> textCoord = new List<int>();
  List<int> indices = new List<int>();

  Obj.fromString(String url, String str) {
  
    Map<String,int> indexTable = new Map<String,int>();
    List<double> _vertCoord = new List<double>();   
    List<int> _textCoord = new List<int>();
    int indexCounter = 0;

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
        _vertCoord.add(double.parse(v[1])); // x
        _vertCoord.add(double.parse(v[2])); // y
        _vertCoord.add(double.parse(v[3])); // z
        return;
      }

      if (line.startsWith("vt ")) {
        // texture coord
        List<String> t = line.split(' ');
        if (t.length != 3) {
          print("OBJ: wrong number of texture coordinates (${t.length - 1} != 2) at line=$lineNum from url=$url: [$rawLine]");
          return;
        }
        _textCoord.add(double.parse(t[1])); // u
        _textCoord.add(double.parse(t[2])); // v
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
          String ind = f[i];
          
          // known unified index?
          int index = indexTable[ind];
          if (index != null) {
            indices.add(index);
            continue;
          }
          
          List<String> v = ind.split('/');
          String vi = v[0];
          int vIndex = int.parse(vi) - 1;
          int vOffset = 3 * vIndex; 
          vertCoord.add(_vertCoord[vOffset + 0]); // x
          vertCoord.add(_vertCoord[vOffset + 1]); // y
          vertCoord.add(_vertCoord[vOffset + 2]); // z
          
          if (v.length > 1) {
            // texture index?
            String ti = v[1];
            if (ti != null && !ti.isEmpty) {
              int tIndex = int.parse(ti) - 1;
              int tOffset = 2 * tIndex;
              textCoord.add(_textCoord[tOffset + 0]); // u
              textCoord.add(_textCoord[tOffset + 1]); // v
            }
          }

          if (v.length > 2) {
            // normal index?
            String ni = v[2];
            if (ni != null && !ni.isEmpty) {
              int nIndex = int.parse(ni) - 1;
            }
          }
          
          // add unified index
          indices.add(indexCounter);
          indexTable[ind] = indexCounter;      
          ++indexCounter;
        }
        return;
      }

      print("OBJ: unknown pattern at line=$lineNum from url=$url: [$rawLine]");
    }
    
    List<String> lines = str.split('\n');
    lines.forEach((String line) => parseLine(line));
    
    print("Obj.fromString: indices.length = ${indices.length}");
    print("Obj.fromString: vertCoord.length = ${vertCoord.length}");
    print("Obj.fromString: textCoord.length = ${textCoord.length}");
  }
}