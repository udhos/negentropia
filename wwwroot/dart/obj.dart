library obj;

class Obj {
  
  static final String prefix_mtllib = "mtllib ";
  static final String prefix_usemtl = "usemtl ";
  static final int prefix_mtllib_len = prefix_mtllib.length;
  static final int prefix_usemtl_len = prefix_usemtl.length;
  
  List<int> indices = new List<int>();
  List<double> vertCoord = new List<double>();   
  List<double> textCoord = new List<double>();
  List<double> normCoord = new List<double>();
  String mtllib;
  String usemtl;

  Obj.fromString(String url, String str) {
  
    Map<String,int> indexTable = new Map<String,int>();
    List<double> _vertCoord = new List<double>();
    List<double> _textCoord = new List<double>();
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
        if (v.length == 4) {
          _vertCoord.add(double.parse(v[1])); // x
          _vertCoord.add(double.parse(v[2])); // y
          _vertCoord.add(double.parse(v[3])); // z
          return;
        }
        if (v.length == 5) {
          double w = double.parse(v[4]);
          _vertCoord.add(double.parse(v[1]) / w); // x
          _vertCoord.add(double.parse(v[2]) / w); // y
          _vertCoord.add(double.parse(v[3]) / w); // z
          return;
        }
        
        print("OBJ: wrong number of vertex coordinates: ${v.length - 1} at line=$lineNum from url=$url: [$line]");
        return;
      }

      if (line.startsWith("vt ")) {
        // texture coord
        List<String> t = line.split(' ');
        if (t.length != 3) {
          print("OBJ: wrong number of texture coordinates (${t.length - 1} != 2) at line=$lineNum from url=$url: [$line]");
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
          print("OBJ: wrong number of face indices (${f.length - 1} != 3) at line=$lineNum from url=$url: [$line]");
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
      
      if (line.startsWith(prefix_mtllib)) {
        String new_mtllib = line.substring(prefix_mtllib_len);
        if (mtllib != null) {
          print("OBJ: mtllib redefinition: from mtllib=$mtllib to mtllib=$new_mtllib");
        }
        mtllib = new_mtllib;
        return;
      }

      if (line.startsWith(prefix_usemtl)) {
        String new_usemtl = line.substring(prefix_usemtl_len);
        if (usemtl != null) {
          print("OBJ: usemtl redefinition: from usemtl=$usemtl to usemtl=$new_usemtl");          
        }
        usemtl = new_usemtl;
        return;
      }

      print("OBJ: unknown pattern at line=$lineNum from url=$url: [$line]");
    }
    
    List<String> lines = str.split('\n');
    lines.forEach((String line) => parseLine(line));
    
    print("Obj.fromString: indices.length = ${indices.length}");
    print("Obj.fromString: vertCoord.length = ${vertCoord.length}");
    print("Obj.fromString: textCoord.length = ${textCoord.length}");
    print("Obj.fromString: normCoord.length = ${normCoord.length}");
    print("Obj.fromString: mtllib = $mtllib");
    print("Obj.fromString: usemtl = $usemtl");
  }
}

class Material {

  static final String prefix_newmtl = "newmtl ";
  static final String prefix_map_Kd = "map_Kd ";
  static final int prefix_newmtl_len = prefix_newmtl.length;
  static final int prefix_map_Kd_len = prefix_map_Kd.length;

  String map_Kd;
  Material(this.map_Kd);
}

Map<String,Material> mtllib_parse(String str, String url) {
  
  Map<String,Material> lib = new Map<String,Material>();
  String currMaterialName;
  int lineNum = 0;
  
  void parseLine(String rawLine) {
    ++lineNum;
    
    String line = rawLine.trim();

    if (line.isEmpty) {
      return;
    }
    
    if (line[0] == '#') {
      return;
    }
    
    if (line.startsWith(Material.prefix_newmtl)) {
      currMaterialName = line.substring(Material.prefix_newmtl_len);
      return;
    }

    if (line.startsWith(Material.prefix_map_Kd)) {
      String map_Kd = line.substring(Material.prefix_map_Kd_len);
      
      if (currMaterialName == null) {
        print("mtllib_parse: url=$url: line=$lineNum: map_Kd=$map_Kd found for undefined material: [$line]");
        return;
      }     
      
      lib[currMaterialName] = new Material(map_Kd);
            
      return;
    }
    
    print("mtllib_parse: url=$url: line=$lineNum: unknown pattern: [$line]");    
  }
  
  List<String> lines = str.split('\n');
  lines.forEach((String line) => parseLine(line));

  return lib;
}
