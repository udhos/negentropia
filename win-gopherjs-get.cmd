@rem win-gopherjs-get.cmd

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv
set GOPHER_DEP=%DEVEL%\negentropia\gopherjs-deps.txt

@rem gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia
@rem go install negentropia\webserv
@rem go install negentropia\world

@rem go get -u github.com/gopherjs/gopherjs
@rem go get -u github.com/gopherjs/webgl
@rem go get -u github.com/gopherjs/websocket
@rem go get -u github.com/udhos/cookie
@rem go get -u honnef.co/go/js/dom

more %GOPHER_DEP% | %DEVEL%\negentropia\win-goget-stdin.cmd

@rem eof
