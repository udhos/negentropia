import 'dart:html';

void main() {
  
  CanvasElement canvas = new CanvasElement();
  
  canvas.id = "main_canvas";
  
  //  document.body.elements.add(canvas);
  query("#canvasbox").append(canvas);
}
