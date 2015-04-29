@rem win-goget-stdin.cmd

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia
@rem go install negentropia\webserv
@rem go install negentropia\world

@rem go get -u github.com/gopherjs/gopherjs
@rem go get -u github.com/gopherjs/webgl
@rem go get -u github.com/gopherjs/websocket
@rem go get -u github.com/udhos/cookie
@rem go get -u honnef.co/go/js/dom

@echo off
setlocal DisableDelayedExpansion

for /F "tokens=*" %%a in ('findstr /n $') do (
  set "line=%%a"
  setlocal EnableDelayedExpansion
  set "line=!line:*:=!"
  @rem echo(!line!
  echo go get !line!
  go get !line!
  endlocal
)

@rem eof
