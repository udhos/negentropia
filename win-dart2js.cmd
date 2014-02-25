@rem win-dart2js

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

@rem dartanalyzer munges DART_SDK
set NEG_DART_SDK=%DART_SDK%
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart
set NEG_DART_MAIN=%NEG_DART_SRC%\negentropia_home.dart

call %NEG_DART_SDK%\bin\dartfmt -w %NEG_DART_SRC%

@rem build client
call %NEG_DART_SDK%\bin\dartanalyzer %NEG_DART_MAIN%
@echo on
call %NEG_DART_SDK%\bin\dart2js -c -DDEBUG=true -o %NEG_DART_MAIN%.js %NEG_DART_MAIN%

@rem eof
