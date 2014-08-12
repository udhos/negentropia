@rem win-dart-test

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk

set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart

set OLD_CD=%CD%
cd %NEG_DART_SRC%\test

@rem call %NEG_DART_SDK%\bin\pub get
@rem @echo on
@rem call %NEG_DART_SDK%\bin\pub upgrade
@rem @echo on

%DART_SDK%\bin\dart vec_test.dart
@echo on
cd %OLD_CD%

@rem eof
