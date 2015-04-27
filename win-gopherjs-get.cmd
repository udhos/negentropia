@rem win-gopherjs

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia
@rem go install negentropia\webserv
@rem go install negentropia\world

go get -u github.com/gopherjs/gopherjs
go get -u github.com/gopherjs/webgl
go get -u github.com/gopherjs/websocket
go get -u github.com/udhos/cookie
go get honnef.co/go/js/dom

@rem eof
