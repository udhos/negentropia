library buffer;

import 'dart:html';
import 'dart:async';
import 'dart:json';

void fetchSquare(String jsonUrl) {
  
  void handleResponse(String response) {
    print("fetched square JSON from URL: $jsonUrl: [$response]");
    Map square;
    try {
      square = parse(response);
    }
    catch (e) {
      print("failure parsing square JSON: $e");
      return;
    }
    print("square JSON parsed: [$square]");
    print("FIXME: create square GL buffer");
  }
  
  void handleError(AsyncError err) {
    print("failure fetching square JSON from URL: $jsonUrl: $err");
  }
  
  // dart magic :-)
  HttpRequest.getString(jsonUrl)
    .then(handleResponse)
    .catchError(handleError);
}



