import 'dart:html';

DivElement _log = new DivElement();
int _line = 0;

void log(String msg) {
  ++_line;
  msg = "$_line $msg";
  print(msg);
  DivElement m = new DivElement();
  m.text = msg;
  _log.append(m);

  List<Element> children = _log.children;
  if (children.length > 20) {
    children.removeAt(0);
  }
}

void main() {
  
  document.body.append(_log);
  
  window.onMouseWheel.listen((WheelEvent e) {
    e.preventDefault();
    log("wheel: type=${e.type} dx=${e.deltaX} dy=${e.deltaY}");
  });
}
