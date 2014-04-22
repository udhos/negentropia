@rem win-build

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

@rem run tests
call %DEVEL%\negentropia\win-gotest.cmd
call %DEVEL%\negentropia\win-benchmark-dart.cmd

@rem build server
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
