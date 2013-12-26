@rem win-benchmark

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart

set OLD_CD=%CD%
cd %NEG_DART_SRC%\benchmark
call %NEG_DART_SDK%\bin\pub get
@echo on
call %NEG_DART_SDK%\bin\pub upgrade
@echo on
%DART_SDK%\bin\dart obj_benchmark.dart
@echo on
cd %OLD_CD%

@rem eof
