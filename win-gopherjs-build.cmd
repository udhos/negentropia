@rem win-gopherjs

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia

@rem go install negentropia\webserv
@rem go install negentropia\world

\tmp\devel\negentropia\webserv\bin\gopherjs install negentropia\negoc

copy \tmp\devel\negentropia\webserv\bin\negoc.js     \tmp\devel\negentropia\wwwroot\negoc
copy \tmp\devel\negentropia\webserv\bin\negoc.js.map \tmp\devel\negentropia\wwwroot\negoc

@rem eof
