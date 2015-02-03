@rem win-dartformat

set DEVEL=c:\tmp\devel
set DART_SDK=c:\dart\dart-sdk
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart

@rem old formatter
@rem call %NEG_DART_SDK%\bin\dartfmt -t -w %NEG_DART_SRC%

@rem new formatter
@rem it puts dartformat executable in: C:\Users\esmarques\AppData\Roaming\Pub\Cache\bin
@rem
@rem APPDATA=C:\Users\esmarques\AppData\Roaming
@rem
call %DART_SDK%\bin\pub global activate dart_style
set PATH=%PATH%;%DART_SDK%\bin
echo PATH=%PATH%
@echo on
%APPDATA%\pub\cache\bin\dartformat -w %NEG_DART_SRC%

@rem eof
