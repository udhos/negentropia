@rem win-build

set DEVEL=c:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

set DART_SDK=c:\dart\dart-sdk

@rem dartanalyzer munges DART_SDK
set NEG_DART_SDK=%DART_SDK%
set NEG_DART_SRC=%DEVEL%\negentropia\wwwroot\dart
set NEG_DART_MAIN=%NEG_DART_SRC%\negentropia_home.dart

@rem install goauth2
@rem
@rem facebook broken:
@rem go get code.google.com/p/goauth2/oauth
@rem
@rem google broken:
@rem go get github.com/robfig/goauth2/oauth
@rem
@rem go get broken:
@rem go get code.google.com/r/jasonmcvetta-goauth2/
@rem
@rem load from git bash:
@rem go get github.com/HairyMezican/goauth2/oauth
@rem
@rem untested:
@rem go get bitbucket.org/gosimple/oauth2

gofmt -s -w %DEVEL%\negentropia\webserv\src\negentropia

@rem build server
go install negentropia\webserv
go install negentropia\world

@rem build client
set OLD_CD=%CD%
cd %NEG_DART_SRC%
call %NEG_DART_SDK%\bin\pub install
@echo on
call %NEG_DART_SDK%\bin\pub update
@echo on
cd %OLD_CD%
call %NEG_DART_SDK%\bin\dartanalyzer %NEG_DART_MAIN%
@echo on
call %NEG_DART_SDK%\bin\dart2js -c -o %NEG_DART_MAIN%.js %NEG_DART_MAIN%

@rem eof
