import 'dart:html';

DivElement logbox = new DivElement();
bool canvasFullscreen = false;

void log(String msg) {
  msg = "${new DateTime.now()} $msg";
  
  print(msg);
  
  DivElement entry = new DivElement();
  entry.text = msg;
  logbox.append(entry);
  
  if (logbox.children.length > 20) {
    logbox.children.removeAt(0);
  }
}

void main() {
  CanvasElement canvas = new CanvasElement();
  canvas.width = 300;
  canvas.height = 300;
  canvas.style.border = "2px solid black";
  document.body.append(canvas);
  document.body.append(logbox);
  
  canvas.onClick.listen((e) {
    if (canvasFullscreen) {
      document.exitFullscreen();
    }
    else {
      canvas.requestFullscreen();
    }
    canvasFullscreen = !canvasFullscreen;
  });
}
