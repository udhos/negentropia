library asset;

class Asset {
  String _mesh;
  String _mtl;
  String _obj;
  String _shader;
  String _texture;

  String get mesh => _mesh;
  String get mtl => _mtl;
  String get obj => _obj;
  String get shader => _shader;
  String get texture => _texture;

  void setRoot(String root) {
    _mesh = "${root}mesh";
    _mtl = "${root}mtl";
    _obj = "${root}obj";
    _shader = "${root}shader";
    _texture = "${root}texture";
  }

  Asset(String root) {
    setRoot(root);
  }
}
