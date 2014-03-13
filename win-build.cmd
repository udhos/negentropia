@rem win-build

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

@rem build server
set GOPATH=%DEVEL%\negentropia\webserv
go install negentropia\webserv
go install negentropia\world
call %DEVEL%\negentropia\win-goinstall.cmd

@rem build client
set NEG_DART_SDK=%DART_SDK%
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart
set OLD_CD=%CD%
cd %NEG_DART_SRC%
call %NEG_DART_SDK%\bin\pub get
@echo on
call %NEG_DART_SDK%\bin\pub upgrade
@echo on
call %DEVEL%\negentropia\win-dart2js.cmd
cd %OLD_CD%

@rem eof
