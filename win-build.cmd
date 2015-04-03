@rem win-build

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem run go fix
call %DEVEL%\negentropia\win-gofix.cmd

@rem run go vet
call %DEVEL%\negentropia\win-govet.cmd

@rem run go tests
call %DEVEL%\negentropia\win-gotest.cmd

@rem build go servers
call %DEVEL%\negentropia\win-goinstall.cmd

@rem build client
call %DEVEL%\negentropia\win-gopherjs-build.cmd

@rem eof
