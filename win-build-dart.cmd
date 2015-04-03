@rem win-build

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

@rem run go fix
call %DEVEL%\negentropia\win-gofix.cmd

@rem run go vet
call %DEVEL%\negentropia\win-govet.cmd

@rem run go tests
call %DEVEL%\negentropia\win-gotest.cmd

@rem build go servers
call %DEVEL%\negentropia\win-goinstall.cmd

@rem format dart code
call %DEVEL%\negentropia\win-dartformat.cmd

@rem build dart client
set NEG_DART_SDK=%DART_SDK%
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart
set OLD_CD=%CD%
@echo on
cd %NEG_DART_SRC%
call %NEG_DART_SDK%\bin\pub get
@echo on
call %NEG_DART_SDK%\bin\pub upgrade
@echo on
call %DEVEL%\negentropia\win-dart2js.cmd
cd %OLD_CD%

@rem run dart tests
call %DEVEL%\negentropia\win-dart-test.cmd
call %DEVEL%\negentropia\win-dart-benchmark.cmd

@rem eof
