@rem win-goinstall

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem test
go test -test.v negentropia\world\server

@rem eof
