@rem webserver-build

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem install goauth2
go get code.google.com/p/goauth2/oauth

@rem build
go install negentropia\webserv

@rem eof
