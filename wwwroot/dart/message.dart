library message;

import 'dart:html';
import 'dart:collection';

import 'logg.dart';

DivElement _messagebox;

int _maxList = 15;
ListQueue<String> _msgList = new ListQueue<String>(_maxList);

void newMessagebox(Element e, String id, CanvasElement c) {
  assert(_messagebox == null);
  _messagebox = _createMessagebox(id, c);
  e.append(_messagebox);
}

DivElement _createMessagebox(String id, CanvasElement c) {

  DivElement mbox = new DivElement();
  mbox.id = id;

  int left = 10 + c.offsetLeft;
  int top = 28 + c.offsetTop;

  mbox.style.border = '2px solid #FFF';
  mbox.style.zIndex = "1";
  mbox.style.position = "absolute";
  mbox.style.width = "300px";
  mbox.style.color = "lightgreen";
  mbox.style.background = "rgba(50,50,50,0.7)";
  mbox.style.textAlign = "left";
  mbox.style.padding = "2px";
  mbox.style.fontSize = 'x-small';

  void repositionBox(Event e) {
    int left = 10 + c.offsetLeft;
    int top = 28 + c.offsetTop;

    mbox.style.left = "${left}px";
    mbox.style.top = "${top}px";

    log("repositionBox: event=$e: left=${mbox.style.left} top=${mbox.style.top}"
        );
  }

  repositionBox(null);

  c.onChange.listen(repositionBox);

  return mbox;
}

void messageUser(String m) {

  _msgList.add(m);

  while (_msgList.length > _maxList) {
    _msgList.removeFirst();
  }

  _messagebox.children.clear();

  _msgList.forEach((m) {
    DivElement d = new DivElement();
    d.text = m;
    _messagebox.children.add(d);
  });
}
