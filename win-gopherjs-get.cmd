@rem win-gopherjs

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia
@rem go install negentropia\webserv
@rem go install negentropia\world

go get -u github.com/gopherjs/gopherjs

@rem eof
