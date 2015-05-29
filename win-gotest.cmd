@rem win-goinstall

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem test
@rem go test -test.v negentropia\world\server
go test negentropia\world\server
go test negentropia\world\util
go test negentropia\world\obj

@rem eof
